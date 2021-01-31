package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/prabudzak/article/model"
	"github.com/prabudzak/article/service"
)

// ArticleDatabase represent article database mysql implementation
type ArticleDatabase struct {
	db *sql.DB
}

// NewArticleDatabase create a new instance of mysql implmentation article database
func NewArticleDatabase(db *sql.DB) *ArticleDatabase {
	return &ArticleDatabase{
		db: db,
	}
}

// GenerateID generate a new id to be assigned to an article
func (a *ArticleDatabase) GenerateID(ctx context.Context) (int, error) {
	var id int

	trx, err := a.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Println(err)
		return 0, err
	}

	_, err = trx.ExecContext(ctx, "UPDATE article_seq SET num = num + 1")
	if err != nil {
		log.Println(err)
		trx.Rollback()
		return 0, err
	}

	row := trx.QueryRowContext(ctx, "SELECT num FROM article_seq LIMIT 1")
	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
		trx.Rollback()
		return 0, err
	}

	err = trx.Commit()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

// Create write a new article to database
func (a *ArticleDatabase) Create(ctx context.Context, article model.Article) error {
	if article.ID == 0 {
		return errors.New("article id is invalid")
	}

	if article.CreatedAt.IsZero() {
		article.CreatedAt = time.Now().UTC()
	}

	if article.UpdatedAt.IsZero() {
		article.UpdatedAt = article.CreatedAt
	}

	_, err := a.db.QueryContext(ctx, "INSERT INTO article (id, author, title, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		article.ID,
		article.Author,
		article.Title,
		article.Body,
		article.CreatedAt,
		article.UpdatedAt,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Get retrieve an article by in from database
func (a *ArticleDatabase) Get(ctx context.Context, id int) (model.Article, error) {
	var article model.Article

	if id == 0 {
		return article, errors.New("id parameter is invalid")
	}

	row := a.db.QueryRowContext(ctx, "SELECT id, author, title, body, created_at, updated_at FROM article WHERE id = ?", id)
	err := row.Scan(&article.ID, &article.Author, &article.Title, &article.Body, &article.CreatedAt, &article.UpdatedAt)
	if err == sql.ErrNoRows {
		return article, service.ErrArticleNotFound
	} else if err != nil {
		log.Println(err)
		return article, err
	}

	return article, nil
}
