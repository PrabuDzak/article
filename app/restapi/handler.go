package restapi

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/prabudzak/article/model"
)

func (a *API) createArticle(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var body createArticleRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.responseMessage(w, http.StatusBadRequest, "bad request")
		return
	}
	defer r.Body.Close()

	err = body.Validate()
	if err != nil {
		a.responseError(w, http.StatusUnprocessableEntity, err)
		return
	}

	article := model.Article{
		Author: body.Author,
		Title:  body.Title,
		Body:   body.Body,
	}

	err = a.articleService.CreateArticle(r.Context(), article)
	if err != nil {
		a.responseError(w, http.StatusInternalServerError, err)
		return
	}

	a.responseMessage(w, http.StatusCreated, "article created")
}

func (a *API) listArticle(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	queryParam, _ := url.ParseQuery(r.URL.RawQuery)

	offset, _ := strconv.ParseInt(queryParam.Get("offset"), 10, 32)
	limit, _ := strconv.ParseInt(queryParam.Get("limit"), 10, 32)

	query := model.ArticleSearchQuery{
		Author:     queryParam.Get("author"),
		Keyword:    queryParam.Get("query"),
		Pagination: model.Pagination{Limit: int(limit), Offset: int(offset)},
	}

	articles, err := a.articleService.SearchArticle(r.Context(), query)
	if err != nil {
		a.responseError(w, http.StatusInternalServerError, err)
		return
	}

	response := response{
		Message: "articles retrieved",
		Data:    articles,
	}

	a.response(w, http.StatusOK, response)
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
