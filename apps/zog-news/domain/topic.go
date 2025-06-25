package domain

import (
	"time"
)

type Topic struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Status defaults to "draft"
type CreateTopicRequest struct {
	Name     string    `json:"name" validate:"required"`
}

type UpdateTopicRequest struct {
	Name     string    `json:"name" validate:"required"`
}

type TopicFilter struct {
	Search string `json:"search" query:"search"`
}
