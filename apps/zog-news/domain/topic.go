package domain

import (
	"time"
)

// Topic represents a topic entity
// @Description Topic entity for categorizing articles
type Topic struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"Technology"`
	CreatedAt time.Time `json:"created_at" example:"2023-06-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-06-01T12:30:00Z"`
}

// CreateTopicRequest represents the request body for creating a topic
// @Description Request body for creating a new topic
type CreateTopicRequest struct {
	Name string `json:"name" validate:"required" example:"Technology"`
}

// UpdateTopicRequest represents the request body for updating a topic
// @Description Request body for updating an existing topic
type UpdateTopicRequest struct {
	Name string `json:"name" validate:"required" example:"Updated Technology"`
}

// TopicFilter represents query parameters for filtering topics
// @Description Query parameters for filtering topics
type TopicFilter struct {
	Search string `json:"search" query:"search" example:"tech"`
}
