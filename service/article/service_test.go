package article_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prabudzak/article/model"
	"github.com/prabudzak/article/service"
	"github.com/prabudzak/article/service/article"
	"github.com/prabudzak/article/service/article/mock"
	"github.com/stretchr/testify/assert"
)

type dependency struct {
	database *mock.MockDatabase
	cache    *mock.MockCache
	indexer  *mock.MockIndexer
}

func initialize(ctrl *gomock.Controller) dependency {
	return dependency{
		database: mock.NewMockDatabase(ctrl),
		cache:    mock.NewMockCache(ctrl),
		indexer:  mock.NewMockIndexer(ctrl),
	}
}

func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name            string
		article         model.Article
		dbGenerateIDErr error
		indexErr        error
		dbCreateErr     error

		expectError bool
	}{
		{
			name: "success",
			article: model.Article{
				Author: "John Doe",
				Title:  "A Valid Title",
				Body:   "A very interesting content",
			},
			expectError: false,
		},
		{
			name: "missing author",
			article: model.Article{
				Title: "A Valid Title",
				Body:  "A very interesting content",
			},
			expectError: true,
		},
		{
			name: "missing title",
			article: model.Article{
				Author: "John Doe",
				Body:   "A very interesting content",
			},
			expectError: true,
		},
		{
			name: "missing body",
			article: model.Article{
				Author: "John Doe",
				Title:  "A Valid Title",
			},
			expectError: true,
		},
		{
			name: "unable to assign article id",
			article: model.Article{
				Author: "John Doe",
				Title:  "A Valid Title",
				Body:   "A very interesting content",
			},
			dbGenerateIDErr: assert.AnError,
			expectError:     true,
		},
		{
			name: "unable to index article",
			article: model.Article{
				Author: "John Doe",
				Title:  "A Valid Title",
				Body:   "A very interesting content",
			},
			indexErr:    assert.AnError,
			expectError: true,
		},
		{
			name: "unable to create article to database",
			article: model.Article{
				Author: "John Doe",
				Title:  "A Valid Title",
				Body:   "A very interesting content",
			},
			dbCreateErr: assert.AnError,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dep := initialize(ctrl)
			dep.database.EXPECT().GenerateID(gomock.Any()).AnyTimes().Return(123, tc.dbGenerateIDErr)
			dep.indexer.EXPECT().Index(gomock.Any(), gomock.Any()).AnyTimes().Return(tc.indexErr)
			dep.database.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(tc.dbCreateErr)

			articleService := article.NewArticleService(dep.database, dep.cache, dep.indexer)

			err := articleService.CreateArticle(context.Background(), tc.article)
			assert.Equal(t, tc.expectError, err != nil)
		})
	}

}

func TestSearchArticle(t *testing.T) {
	tests := []struct {
		name                   string
		indexSearchArticleIDs  []int
		indexSearchErr         error
		cacheGetErr            error
		expectedArticlesLength int
		expectErr              bool
	}{
		{
			name:                   "all indexed article returned",
			indexSearchArticleIDs:  []int{1, 2, 3, 4, 5},
			expectedArticlesLength: 5,
			expectErr:              false,
		},
		{
			name:                   "not all indexed article returned, article not found in cache",
			indexSearchArticleIDs:  []int{1, 2, 3, 4, 400},
			expectedArticlesLength: 4,
			expectErr:              false,
		},
		{
			name:                  "unable to search articles",
			indexSearchArticleIDs: []int{1, 2, 3},
			indexSearchErr:        assert.AnError,
			expectErr:             true,
		},
		{
			name:                   "not all indexed article returned, cache error",
			indexSearchArticleIDs:  []int{1, 2, 3},
			cacheGetErr:            assert.AnError,
			expectedArticlesLength: 0,
			expectErr:              false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dep := initialize(ctrl)
			dep.indexer.EXPECT().Search(gomock.Any(), gomock.Any()).AnyTimes().Return(tc.indexSearchArticleIDs, tc.indexSearchErr)
			dep.cache.EXPECT().Get(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, id int) (model.Article, error) {
				// simulate  condition all article with id > 100 not found
				if id > 100 {
					return model.Article{}, service.ErrArticleNotFound
				}

				return model.Article{
					ID:     id,
					Title:  fmt.Sprintf("title %d", id),
					Body:   fmt.Sprintf("body %d", id),
					Author: fmt.Sprintf("author%d", id),
				}, tc.cacheGetErr
			})

			articleService := article.NewArticleService(dep.database, dep.cache, dep.indexer)

			articles, err := articleService.SearchArticle(context.Background(), model.ArticleSearchQuery{})
			assert.Equal(t, tc.expectErr, err != nil)
			assert.Len(t, articles, tc.expectedArticlesLength)
		})
	}
}
