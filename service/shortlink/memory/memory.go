package memory

import (
	"context"
	"errors"
	"log"

	"github.com/prabudzak/article/model"
)

type UrlShortenerInMemoryDatabase struct {
	db map[string]model.ShortenURL
}

func NewUrlShortenerInMemoryDatabase() *UrlShortenerInMemoryDatabase {
	return &UrlShortenerInMemoryDatabase{
		db: make(map[string]model.ShortenURL),
	}
}

func (u *UrlShortenerInMemoryDatabase) Create(ctx context.Context, shortenURL model.ShortenURL) error {
	if _, found := u.db[shortenURL.ID]; found {
		return errors.New("already exist")
	}

	u.db[shortenURL.ID] = shortenURL
	log.Println(u.db)

	return nil
}

func (u *UrlShortenerInMemoryDatabase) GetByID(ctx context.Context, ID string) (model.ShortenURL, error) {
	shortenURL, found := u.db[ID]
	if !found {
		return model.ShortenURL{}, errors.New("id not found")
	}

	return shortenURL, nil
}

func (u *UrlShortenerInMemoryDatabase) IncrementViewCounter(ctx context.Context, ID string) error {
	shortenURL, found := u.db[ID]
	if !found {
		return errors.New("not found")
	}

	shortenURL.Counter = shortenURL.Counter + 1
	u.db[shortenURL.ID] = shortenURL

	return nil
}
