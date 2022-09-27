package restapi

import (
	"errors"
)

type createArticleRequest struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (c createArticleRequest) Validate() error {
	if c.Author == "" {
		return errors.New("author is blank")
	}
	if c.Title == "" {
		return errors.New("title is blank")
	}
	if c.Body == "" {
		return errors.New("body is blank")
	}
	return nil
}

type createShortUrlRequest struct {
	URL string `json:"url"`
}
