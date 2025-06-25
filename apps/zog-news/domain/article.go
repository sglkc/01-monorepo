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

    // List of topic IDs associated with the article
    Topics    []string  `json:"topics,omitempty" db:"-"`
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

func (a *Article) AddTopic(topic string) {
    for _, t := range a.Topics {
        if t == topic {
            return // Topic already exists
        }
    }
    a.Topics = append(a.Topics, topic)
}

func (a *Article) RemoveTopic(topic string) {
    for i, t := range a.Topics {
        if t == topic {
            a.Topics = append(a.Topics[:i], a.Topics[i+1:]...)
            break
        }
    }
}

func (a *Article) HasTopic(topic string) bool {
    for _, t := range a.Topics {
        if t == topic {
            return true
        }
    }
    return false
}
