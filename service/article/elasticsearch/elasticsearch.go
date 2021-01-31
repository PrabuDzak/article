package elasticsearch

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/olivere/elastic"

	"github.com/prabudzak/article/model"
)

// ArticleIndexer represent article indexer elasticsearch implementation
type ArticleIndexer struct {
	client *elastic.Client

	indexName string
}

// NewArticleIndexer create a new instance of elasticsearch implementation article indexer
func NewArticleIndexer(client *elastic.Client, indexName string) *ArticleIndexer {
	return &ArticleIndexer{
		client:    client,
		indexName: indexName,
	}
}

// Index put an index for a given article
func (a *ArticleIndexer) Index(ctx context.Context, article model.Article) error {
	if article.ID == 0 {
		return errors.New("article id is invalid")
	}

	_, err := a.client.Index().
		Index(a.indexName).
		Type("article").
		Id(strconv.FormatInt(int64(article.ID), 10)).
		BodyJson(article).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Remove delete an article index by id
func (a *ArticleIndexer) Remove(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("article id is invalid")
	}

	_, err := a.client.Delete().
		Index(a.indexName).
		Type("article").
		Id(strconv.FormatInt(int64(id), 10)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Search search articles by given search query parameter
func (a *ArticleIndexer) Search(ctx context.Context, query model.ArticleSearchQuery) ([]int, error) {
	q := elastic.NewBoolQuery()

	if query.Pagination.Limit <= 0 || query.Pagination.Limit > 100 {
		query.Pagination.Limit = 20
	}

	if query.Author != "" {
		q.Filter(elastic.NewTermQuery("author", query.Author))
	}

	if query.Keyword != "" {
		q.Must(elastic.NewMultiMatchQuery(query.Keyword, "title", "body"))
	}

	sort := elastic.NewFieldSort("created_at").Desc()

	result, err := a.client.Search().
		Index(a.indexName).
		Type("article").
		Query(q).
		SortBy(sort).
		From(query.Pagination.Offset).
		Size(query.Pagination.Limit).
		FetchSource(false).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ids := []int{}
	for _, hit := range result.Hits.Hits {
		id, err := strconv.ParseInt(hit.Id, 10, 32)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		ids = append(ids, int(id))
	}

	return ids, nil
}
