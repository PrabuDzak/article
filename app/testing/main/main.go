package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/subosito/gotenv"

	"github.com/prabudzak/article/model"
)

type testContext struct {
	url string
	wg  *sync.WaitGroup
}

func main() {
	gotenv.Load()
	t := &testContext{
		url: fmt.Sprintf("%s:%s", os.Getenv("URL"), os.Getenv("PORT")),
		wg:  &sync.WaitGroup{},
	}
	t.wg.Add(4)
	go testHealthz(t)
	go testCreateArticle(t)
	go testGetArticleByAuthor(t)
	go testGetArticleByKeyword(t)
	t.wg.Wait()
	fmt.Println("done")
}

func testHealthz(t *testContext) {
	defer t.wg.Done()

	resp, err := http.DefaultClient.Get(t.url + "/healthz")
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("article service not running")
	}
}

func testCreateArticle(t *testContext) {
	defer t.wg.Done()

	tests := []struct {
		name               string
		body               string
		expectedStatusCode int
	}{
		{
			name: "success created",
			body: `
				{
					"author": "john doe",
					"title": "A Valid Title",
					"body": "A very interesting content"
				}
			`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "bad request body",
			body: `
				not a valid json body
			`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "author is blank",
			body: `
				{
					"title": "A Valid Title",
					"body": "A very interesting content"
				}
			`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "title is blank",
			body: `
				{
					"author": "john doe",
					"body": "A very interesting content"
				}
			`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "body is blank",
			body: `
				{
					"author": "john doe",
					"title": "A Valid Title"
				}
			`,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		resp, err := http.DefaultClient.Post(t.url+"/articles", "application/json", strings.NewReader(tc.body))
		assertNoResponseError(tc.name, err)
		assertStatusCode(tc.name, tc.expectedStatusCode, resp.StatusCode)
	}
}

func testGetArticleByAuthor(t *testContext) {
	defer t.wg.Done()

	name := "get article by author"
	authorNum := time.Now().Unix()
	body := fmt.Sprintf(`
		{
			"author": "author%d",
			"title": "A Valid Title",
			"body": "A very interesting content"
		}
	`, authorNum)

	n := rand.Int()%5 + 3
	for i := 0; i < n; i++ {
		resp, err := http.DefaultClient.Post(t.url+"/articles", "application/json", strings.NewReader(body))
		assertNoResponseError(name, err)
		assertStatusCode(name, http.StatusCreated, resp.StatusCode)
	}

	time.Sleep(1000 * time.Millisecond)

	resp, err := http.DefaultClient.Get(fmt.Sprintf("%s/articles?author=author%d", t.url, authorNum))
	assertNoResponseError(name, err)
	assertStatusCode(name, http.StatusOK, resp.StatusCode)
	assertArticleResponseAuthor(name, resp.Body, fmt.Sprintf("author%d", authorNum))
}

func testGetArticleByKeyword(t *testContext) {
	defer t.wg.Done()

	name := "get article by keyword"
	titleNum := time.Now().Unix()
	bodyNum := time.Now().Add(2 * time.Hour).Unix()
	body := fmt.Sprintf(`
		{
			"author": "bob",
			"title": "title%d",
			"body": "content%d"
		}
	`, titleNum, bodyNum)

	n := rand.Int()%5 + 3
	for i := 0; i < n; i++ {
		resp, err := http.DefaultClient.Post(t.url+"/articles", "application/json", strings.NewReader(body))
		assertNoResponseError(name, err)
		assertStatusCode(name, http.StatusCreated, resp.StatusCode)
	}

	time.Sleep(1000 * time.Millisecond)

	resp, err := http.DefaultClient.Get(fmt.Sprintf("%s/articles?keyword=title%d", t.url, titleNum))
	assertNoResponseError(name, err)
	assertStatusCode(name, http.StatusOK, resp.StatusCode)
	assertArticleResponseKeyword(name, resp.Body, fmt.Sprintf("title%d", titleNum))

	resp, err = http.DefaultClient.Get(fmt.Sprintf("%s/articles?keyword=content%d", t.url, bodyNum))
	assertNoResponseError(name, err)
	assertStatusCode(name, http.StatusOK, resp.StatusCode)
	assertArticleResponseKeyword(name, resp.Body, fmt.Sprintf("content%d", bodyNum))
}

func assertStatusCode(name string, expectation int, actual int) bool {
	if actual != expectation {
		fmt.Printf("FAIL %s:\n"+
			"Want: %d\n"+
			"Actual: %d\n",
			name, expectation, actual)
	}
	return actual == expectation
}

func assertNoResponseError(name string, err error) bool {
	if err != nil {
		fmt.Printf("FAIL %s:\n"+
			"Got an error: %s\n"+
			name, err.Error())
	}
	return err == nil
}

func assertArticleResponseAuthor(name string, reader io.ReadCloser, author string) {
	type response struct {
		Message string          `json:"message"`
		Data    []model.Article `json:"data"`
	}

	var resp response
	err := json.NewDecoder(reader).Decode(&resp)
	if err != nil {
		log.Fatalln(err)
	}

	if len(resp.Data) == 0 {
		fmt.Printf("FAIL %s:\n"+
			"article length is 0\n",
			name)
	}

	for _, article := range resp.Data {
		if article.Author != author {
			fmt.Printf("FAIL %s:\n"+
				"Want: %s\n"+
				"Actual: %s\n",
				name, author, article.Author)
		}
	}
}

func assertArticleResponseKeyword(name string, reader io.ReadCloser, keyword string) {
	type response struct {
		Message string          `json:"message"`
		Data    []model.Article `json:"data"`
	}
	var resp response
	err := json.NewDecoder(reader).Decode(&resp)
	if err != nil {
		log.Fatalln(err)
	}

	if len(resp.Data) == 0 {
		fmt.Printf("FAIL %s:\n"+
			"article length is 0\n",
			name)
	}

	for _, article := range resp.Data {
		if strings.Contains(article.Title, keyword) || strings.Contains(article.Body, keyword) {
			return
		}
	}

	fmt.Printf("FAIL %s:\n"+
		"article with keyword %s not found\n",
		name, keyword)
}
