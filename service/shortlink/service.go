package shortlink

import (
	"context"
	"log"
	"time"

	"github.com/prabudzak/article/model"
)

type Database interface {
	Create(ctx context.Context, shortenURL model.ShortenURL) error
	GetByID(ctx context.Context, ID string) (model.ShortenURL, error)
	IncrementViewCounter(ctx context.Context, ID string) error
}

type ShortUrlIdGenerator interface {
	Generate(ctx context.Context) (string, error)
}

type UrlShortenerService struct {
	idGenerator ShortUrlIdGenerator
	database    Database
}

func NewUrlShorten(idGenerator ShortUrlIdGenerator, db Database) *UrlShortenerService {
	return &UrlShortenerService{
		idGenerator: idGenerator,
		database:    db,
	}
}

func (u *UrlShortenerService) CreateShortURL(ctx context.Context, originalURL string) (model.ShortenURL, error) {
	id, err := u.idGenerator.Generate(ctx)
	if err != nil {
		log.Println(err)
		return model.ShortenURL{}, err
	}

	shortenURL := model.ShortenURL{
		ID:          id,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		URL:         "generate",
	}

	err = u.database.Create(ctx, shortenURL)
	if err != nil {
		log.Println(err)
		return model.ShortenURL{}, err
	}

	return shortenURL, nil
}

func (u *UrlShortenerService) GetShortURL(ctx context.Context, shortUrlID string) (model.ShortenURL, error) {
	return u.database.GetByID(ctx, shortUrlID)
}

func (u *UrlShortenerService) GetShortUrlForRedirect(ctx context.Context, shortUrlID string) (model.ShortenURL, error) {
	shortUrl, err := u.database.GetByID(ctx, shortUrlID)
	if err != nil {
		log.Println(err)
		return model.ShortenURL{}, err
	}

	err = u.database.IncrementViewCounter(ctx, shortUrlID)
	if err != nil {
		log.Println(err)
	}

	return shortUrl, err
}
