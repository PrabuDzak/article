package restapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prabudzak/article/app/restapi"
	"github.com/prabudzak/article/model"
	"github.com/prabudzak/article/service/mock"
	"github.com/stretchr/testify/assert"
)

type dependency struct {
	articleService *mock.MockArticleService
}

func initialize(ctrl *gomock.Controller) dependency {
	return dependency{
		articleService: mock.NewMockArticleService(ctrl),
	}
}

func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name               string
		body               string
		createArticleErr   error
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
		{
			name: "unable to create article",
			body: `
				{
					"author": "john doe",
					"title": "A Valid Title",
					"body": "A very interesting content"
				}
			`,
			createArticleErr:   assert.AnError,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dep := initialize(ctrl)
			dep.articleService.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).MaxTimes(1).Return(tc.createArticleErr)

			api := restapi.New(dep.articleService, nil)
			router := api.Router()
			server := httptest.NewServer(router)
			defer server.Close()

			url := server.URL + "/articles"
			resp, err := http.DefaultClient.Post(url, "appplication/json", strings.NewReader(tc.body))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestListArticle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		searchArticle      []model.Article
		searchArticleErr   error
		expectedQuery      model.ArticleSearchQuery
		expectedStatusCode int
	}{
		{
			name:          "articles retrieved",
			path:          "/articles",
			searchArticle: []model.Article{},
			expectedQuery: model.ArticleSearchQuery{
				Author:     "",
				Keyword:    "",
				Pagination: model.Pagination{Limit: 0, Offset: 0},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "articles retrieved, with author param",
			path: "/articles?author=john%20man",
			expectedQuery: model.ArticleSearchQuery{
				Author:     "john man",
				Keyword:    "",
				Pagination: model.Pagination{Limit: 0, Offset: 0},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "articles retrieved, with keyword param",
			path: "/articles?query=some%20keyword",
			expectedQuery: model.ArticleSearchQuery{
				Author:     "",
				Keyword:    "some keyword",
				Pagination: model.Pagination{Limit: 0, Offset: 0},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "articles retrieved, with pagination param",
			path: "/articles?limit=20&offset=5",
			expectedQuery: model.ArticleSearchQuery{
				Author:     "",
				Keyword:    "",
				Pagination: model.Pagination{Limit: 20, Offset: 5},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:             "unable to retreive articles",
			path:             "/articles",
			searchArticle:    []model.Article{},
			searchArticleErr: assert.AnError,
			expectedQuery: model.ArticleSearchQuery{
				Author:     "",
				Keyword:    "",
				Pagination: model.Pagination{Limit: 0, Offset: 0},
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dep := initialize(ctrl)
			dep.articleService.EXPECT().SearchArticle(gomock.Any(), tc.expectedQuery).MaxTimes(1).Return(tc.searchArticle, tc.searchArticleErr)

			api := restapi.New(dep.articleService, nil)
			router := api.Router()
			server := httptest.NewServer(router)
			defer server.Close()

			url := server.URL + tc.path
			resp, err := http.DefaultClient.Get(url)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
		})
	}
}
