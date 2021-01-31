package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/prabudzak/article/model"
	"github.com/prabudzak/article/service"
)

const articleKey = "20210130/articles/%d"

// ArticleCache represent article cache redis implementation
type ArticleCache struct {
	client *redis.Client
}

// NewArticleCache create a new instance of redis implmentation article cache
func NewArticleCache(redisClient *redis.Client) *ArticleCache {
	return &ArticleCache{
		client: redisClient,
	}
}

// Cache write article to cache storage
func (a *ArticleCache) Cache(ctx context.Context, article model.Article) error {
	if article.ID == 0 {
		return errors.New("article id is invalid")
	}

	jsoned, _ := json.Marshal(article)

	key := fmt.Sprintf(articleKey, article.ID)
	err := a.client.Set(key, jsoned, 0).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Get retrieve an article by id from cache storage
func (a *ArticleCache) Get(ctx context.Context, id int) (model.Article, error) {
	var article model.Article

	if id == 0 {
		return article, errors.New("id parameter is invalid")
	}

	key := fmt.Sprintf(articleKey, id)
	result, err := a.client.Get(key).Result()
	if err == redis.Nil {
		return article, service.ErrArticleNotFound
	} else if err != nil {
		log.Println(err)
		return article, err
	}

	err = json.Unmarshal([]byte(result), &article)
	if err != nil {
		log.Println(err)
		return article, err
	}

	if article.ID == 0 {
		return article, service.ErrArticleNotFound
	}

	return article, nil
}
