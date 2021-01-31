package article

import (
	"context"
	"errors"
	"time"

	"github.com/prabudzak/article/event"
	"github.com/prabudzak/article/model"
	"github.com/prabudzak/article/service"
)

//go:generate mockgen -package=mock -source=service.go -destination=mock/service.go

// Database represent article persistent storage
type Database interface {
	GenerateID(ctx context.Context) (int, error)
	Create(ctx context.Context, article model.Article) error
	Get(ctx context.Context, id int) (model.Article, error)
}

// Cache represent article cache storage
type Cache interface {
	Cache(ctx context.Context, article model.Article) error
	Get(ctx context.Context, id int) (model.Article, error)
}

// Indexer represent article indexer
type Indexer interface {
	Index(ctx context.Context, article model.Article) error
	Remove(ctx context.Context, id int) error
	Search(ctx context.Context, query model.ArticleSearchQuery) ([]int, error)
}

// Service represent article service implementation
type Service struct {
	database Database
	cache    Cache
	indexer  Indexer
}

// NewArticleService create a new article service instance
func NewArticleService(database Database, cache Cache, indexer Indexer) *Service {
	return &Service{
		database: database,
		cache:    cache,
		indexer:  indexer,
	}
}

// CreateArticle write a new article and dispatch article created event if
// successfully written
func (s *Service) CreateArticle(ctx context.Context, article model.Article) error {
	if article.Author == "" {
		return errors.New("article author is blank")
	}
	if article.Title == "" {
		return errors.New("article title is blank")
	}
	if article.Body == "" {
		return errors.New("article body is blank")
	}

	id, err := s.database.GenerateID(ctx)
	if err != nil {
		return err
	}

	article.ID = id
	article.CreatedAt = time.Now().UTC()

	err = s.indexer.Index(ctx, article)
	if err != nil {
		return err
	}

	err = s.database.Create(ctx, article)
	if err != nil {
		event.Dispatch(ctx, event.ArticleCreateFailed{Article: article})
		return err
	}

	event.Dispatch(ctx, event.ArticleCreated{Article: article})
	return nil
}

// SearchArticle search list of article from given parameter
func (s *Service) SearchArticle(ctx context.Context, query model.ArticleSearchQuery) ([]model.Article, error) {
	ids, err := s.indexer.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	articles := []model.Article{}
	for _, id := range ids {
		article, err := s.cache.Get(ctx, id)
		if err == service.ErrArticleNotFound {
			event.Dispatch(ctx, event.ArticleNotFound{ArticleID: id})
			continue
		} else if err != nil {
			continue
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (s *Service) SubscriberRedispatchArticleCreate(ctx context.Context, e event.Event) error {
	var id int

	switch message := e.(type) {
	case event.ArticleNotFound:
		id = message.ArticleID
	case event.ArticleCachingFailed:
		id = message.Article.ID
	default:
		return errors.New("subscribed to unprocessable event")
	}

	article, err := s.database.Get(ctx, id)
	if err != nil {
		return err
	}

	event.Dispatch(ctx, event.ArticleCreated{Article: article})
	return nil
}

func (s *Service) SubscriberCacheArticle(ctx context.Context, e event.Event) error {
	var article model.Article

	switch message := e.(type) {
	case event.ArticleCreated:
		article = message.Article
	default:
		return errors.New("subscribed to unprocessable event")
	}

	err := s.cache.Cache(ctx, article)
	if err != nil {
		event.Dispatch(ctx, event.ArticleCachingFailed{Article: article})
		return err
	}

	return nil
}

func (s *Service) SubscriberRemoveArticleIndex(ctx context.Context, e event.Event) error {
	var id int

	switch message := e.(type) {
	case event.ArticleCreateFailed:
		id = message.Article.ID
	default:
		return errors.New("subscribed to unprocessable event")
	}

	err := s.indexer.Remove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
