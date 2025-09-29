package blog

import (
	"time"
)

type Article struct {
	Id        int       `json:"id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(authorID string, title string, content string) *Article {
	now := time.Now()
	return &Article{
		Id:        -1,
		AuthorID:  authorID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
