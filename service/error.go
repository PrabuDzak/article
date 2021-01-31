package service

import "errors"

// GeneralError repesent general service error
type GeneralError error

var (
	// ErrArticleNotFound represent article not found service error
	ErrArticleNotFound GeneralError = errors.New("article not found")
)
