package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/corpus/internal/blog"
	_ "github.com/glebarez/go-sqlite"
)

func CreateArticleTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS articles (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id  TEXT    NOT NULL,
			title      TEXT    NOT NULL,
			content    TEXT    NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP),
			updated_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP)
			);
			`
	return db.Exec(sql)
}

func FillArticleTable(db *sql.DB) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(`INSERT INTO articles (author_id, title, content) VALUES (?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	var res sql.Result
	for i := 1; i <= 10; i++ {
		author := fmt.Sprintf("author-%d", i)
		title := fmt.Sprintf("Sample Article %d", i)

		// create a richer markdown content for the sample article
		content := fmt.Sprintf("# %s\n\n_By %s_\n\nThis is an example article written to populate the database for testing purposes.\n\nSample Article %d explores the idea of creating useful placeholder content, and demonstrates code blocks, lists, and links in Markdown.\n\n## Overview\n\n- Purpose: demonstrate article rendering\n- Format: Markdown with headings, code, and lists\n\n## Example Code\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello from Sample Article %d\")\n}\n```\n\n## Conclusion\n\nThis sample article is intentionally verbose so you can see Markdown rendering in the TUI. Enjoy!\n\n", title, author, i, i)

		r, err := stmt.Exec(author, title, content)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		res = r
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func GetAllArticles(db *sql.DB) ([]*blog.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
        SELECT id, author_id, title, content, created_at, updated_at
        FROM articles
        ORDER BY created_at DESC
    `

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

		// The sqlite driver you're using usually maps TIMESTAMP -> time.Time.
		// If your driver returns string, change these vars to string and parse with time.Parse.
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

func DropArticleTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(`DROP TABLE IF EXISTS articles`)
}
