package service

import (
	"context"

	"github.com/prabudzak/article/model"
)

//go:generate mockgen -package=mock -source=service.go -destination=mock/service.go

// ArticleService represent article service interface
type ArticleService interface {
	CreateArticle(ctx context.Context, article model.Article) error
	SearchArticle(ctx context.Context, query model.ArticleSearchQuery) ([]model.Article, error)
}

type ShortenUrlService interface {
	CreateShortURL(ctx context.Context, originalURL string) (model.ShortenURL, error)
	GetShortURL(ctx context.Context, shortUrlID string) (model.ShortenURL, error)
	GetShortUrlForRedirect(ctx context.Context, shortUrlID string) (model.ShortenURL, error)
}
