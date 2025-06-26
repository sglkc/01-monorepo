package domain

import (
	"errors"
	"time"
)

// https://maddevs.io/writeups/personal-experience-using-clean-architecture-in-golang/
type ArticleStatus string

const (
    StatusDraft ArticleStatus = "draft"
    StatusPublished ArticleStatus = "published"
    StatusDeleted ArticleStatus = "deleted"
)

type Article struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Author    string        `json:"author"`
	Status    ArticleStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`

    // List of topic IDs associated with the article
    TopicIDs    []string  `json:"topic_ids,omitempty" db:"-"`

    // Full Topic objects associated with the article for responses
    Topics      []Topic   `json:"topics,omitempty" db:"-"`
}

// Status defaults to "draft"
type CreateArticleRequest struct {
	Title     string        `json:"title" validate:"required"`
	Content   string        `json:"content" validate:"required"`
	Author    string        `json:"author" validate:"required"`
	Status    ArticleStatus `json:"status" validate:"oneof=draft published archived"`
}

type UpdateArticleRequest struct {
	Title     string        `json:"title" validate:"required"`
	Content   string        `json:"content" validate:"required"`
	Author    string        `json:"author" validate:"required"`
	Status    ArticleStatus `json:"status" validate:"oneof=draft published archived"`
}

type ArticleFilter struct {
    Search string           `json:"search" query:"search"`
    Status ArticleStatus    `json:"status" query:"status"`
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
