package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/corpus/internal/blog"
)

type ArticleRepo interface {
	GetAll() ([]*blog.Article, error)
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

func (repo *SQLArticleRepo) GetAll() ([]*blog.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
        SELECT id, author_id, title, content, created_at, updated_at
        FROM articles
        ORDER BY created_at DESC
    `

	rows, err := repo.db.QueryContext(ctx, q)
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

// rowScanner is satisfied by *sql.Row and *sql.Rows
type rowScanner interface {
	Scan(dest ...interface{}) error
}

func (repo *SQLArticleRepo) GetById(id int) (*blog.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
		SELECT id, author_id, title, content, created_at, updated_at
		FROM articles
		WHERE id = ?
		LIMIT 1
	`

	row := repo.db.QueryRowContext(ctx, q, id)
	article, err := scanArticle(row)

	if err != nil {
		return nil, err
	}

	return article, nil
}

// Once added, the id foe the created id is inserted into the article type
func (repo *SQLArticleRepo) Add(article *blog.Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
		INSERT INTO articles (author_id, title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	res, err := repo.db.ExecContext(ctx, q, article.AuthorID, article.Title, article.Content, article.CreatedAt, article.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	article.Id = int(id)
	return nil
}

func (repo *SQLArticleRepo) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `DELETE FROM articles WHERE id = ?`

	res, err := repo.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	// treat delete as idempotent; don't return error if no rows were affected
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

// scanArticle scans a single article from any object that implements Scan
func scanArticle(rs rowScanner) (*blog.Article, error) {
	var (
		id        int
		authorID  string
		title     string
		content   string
		createdAt time.Time
		updatedAt time.Time
	)

	if err := rs.Scan(&id, &authorID, &title, &content, &createdAt, &updatedAt); err != nil {
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

	return a, nil
}

func scanArticlesUtil(rows *sql.Rows) ([]*blog.Article, error) {
	var out []*blog.Article
	for rows.Next() {
		article, err := scanArticle(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, article)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
