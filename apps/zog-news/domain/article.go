package domain

import (
	"time"
)

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Status defaults to "draft"
type CreateArticleRequest struct {
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	Status    string    `json:"status" validate:"oneof=draft published archived"`
}

type UpdateArticleRequest struct {
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	Status    string    `json:"status" validate:"oneof=draft published archived"`
}

type ArticleFilter struct {
	Search string `json:"search" query:"search"`
}
