package domain

import (
	"errors"
	"time"
)

// https://maddevs.io/writeups/personal-experience-using-clean-architecture-in-golang/

// ArticleStatus represents the status of an article
//	@Description	Article status enum
//	@Enum			draft,published,archived
type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusPublished ArticleStatus = "published"
	StatusDeleted   ArticleStatus = "deleted"
)

// Article represents an article entity
//	@Description	Article entity with associated topics
type Article struct {
	ID      string        `json:"id" example:"d4b8583d-5038-4838-bcd7-3d8dddfedd6a"`
	Title   string        `json:"title" example:"Breaking News: Important Update"`
	Content string        `json:"content" example:"This is the content of the article..."`
	Author  string        `json:"author" example:"John Doe"`
	Status  ArticleStatus `json:"status" example:"published"`

	// List of topic IDs associated with the article
	TopicIDs []string `json:"-,omitempty" db:"-"`

	// Full Topic objects associated with the article for responses
	Topics []Topic `json:"topics,omitempty" db:"-"`

	CreatedAt time.Time `json:"created_at" example:"2023-06-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-06-01T12:30:00Z"`
}

// CreateArticleRequest represents the request body for creating an article
//	@Description	Request body for creating a new article
type CreateArticleRequest struct {
	Title   string        `json:"title" validate:"required" example:"Breaking News: Important Update"`
	Content string        `json:"content" validate:"required" example:"This is the content of the article..."`
	Author  string        `json:"author" validate:"required" example:"John Doe"`
	Status  ArticleStatus `json:"status" validate:"oneof=draft published archived" example:"draft"`
}

// UpdateArticleRequest represents the request body for updating an article
//	@Description	Request body for updating an existing article
type UpdateArticleRequest struct {
	Title   string        `json:"title" validate:"required" example:"Updated Breaking News"`
	Content string        `json:"content" validate:"required" example:"This is the updated content..."`
	Author  string        `json:"author" validate:"required" example:"Jane Doe"`
	Status  ArticleStatus `json:"status" validate:"oneof=draft published archived" example:"published"`
}

// ArticleFilter represents query parameters for filtering articles
//	@Description	Query parameters for filtering articles
type ArticleFilter struct {
	Search string        `json:"search" query:"search" example:"breaking news"`
	Status ArticleStatus `json:"status" query:"status" example:"published"`
	Topic  string        `json:"topic" query:"topic" example:"technology"`
}

func (a *Article) HasTopicID(topic string) error {
	for _, t := range a.TopicIDs {
		if t == topic {
			return errors.New("topic already exists")
		}
	}
	return nil
}

// TODO: check if topic id is deleted??
func (a *Article) AddTopicID(topic string) error {
	if err := a.HasTopicID(topic); err != nil {
		return err
	}
	a.TopicIDs = append(a.TopicIDs, topic)
	return nil
}

func (a *Article) RemoveTopicID(topic string) error {
	if len(a.TopicIDs) == 0 {
		return errors.New("no topics to remove")
	}
	if err := a.HasTopicID(topic); err != nil {
		return err
	}

	for i, t := range a.TopicIDs {
		if t == topic {
			a.TopicIDs = append(a.TopicIDs[:i], a.TopicIDs[i+1:]...)
			break
		}
	}
	return nil
}
