package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/corpus/internal/blog"
)

type ArticleRepo interface {
	GetByAuthorId(author_id int) ([]*blog.Article, error)
	GetById(id int) (*blog.Article, error)
	Add(article *blog.Article) error
	Delete(id int) error
}

type SQLArticleRepo struct {
	db *sql.DB
}

func NewSQLArticleRepo(db *sql.DB) *SQLArticleRepo {
	return &SQLArticleRepo{db: db}
}

func (repo *SQLArticleRepo) GetByAuthorId(author_id int) ([]*blog.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
        SELECT id, author_id, title, content, created_at, updated_at
        FROM articles
		WHERE author_id = ?
        ORDER BY created_at DESC
    `

	rows, err := repo.db.QueryContext(ctx, q, author_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out, err := scanArticlesUtil(rows)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func scanArticleUtil(rows *sql.Row) (blog.Article, error){
	
}
func scanArticlesUtil(rows *sql.Rows) ([]*blog.Article, error){
	var out []*blog.Article
	for rows.Next() {
		var (
			id        int
			authorID  string
			title     string
			content   string
			createdAt time.Time
			updatedAt time.Time
		)

		if err := rows.Scan(&id, &authorID, &title, &content, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		a := &blog.Article{
			Id:        id,
			AuthorID:  authorID,
			Title:     title,
			Content:   content,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		out = append(out, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}


func (repo *SQLArticleRepo) GetById(id int) (*blog.Article, error){

}
	Add(article *blog.Article) error
	Delete(id int) error