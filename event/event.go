package event

import "github.com/prabudzak/article/model"

// Event represent a dispatched event message
type Event interface {
	String() string
}

// ArticleCreated represent article created event
type ArticleCreated struct {
	Article model.Article
}

func (a ArticleCreated) String() string {
	return "event_article_created"
}

// ArticleCreateFailed represent article create failed event
type ArticleCreateFailed struct {
	Article model.Article
}

func (a ArticleCreateFailed) String() string {
	return "event_article_create_failed"
}

// ArticleNotFound represent article not found event
type ArticleNotFound struct {
	ArticleID int
}

func (a ArticleNotFound) String() string {
	return "event_article_not_found"
}

// ArticleCachingFailed represent article caching failed event
type ArticleCachingFailed struct {
	Article model.Article
}

func (a ArticleCachingFailed) String() string {
	return "event_article_caching_failed"
}
