package service

import (
	"context"
	"fmt"
	"zog-news/domain"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *domain.CreateTopicRequest) (*domain.Topic, error)
	GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, topic *domain.Topic) (*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error

    GetTopicArticles(ctx context.Context, id uuid.UUID) ([]domain.Article, error)
}

type TopicService struct {
	topicRepo TopicRepository
}

func NewTopicService(a TopicRepository) *TopicService {
	return &TopicService{
		topicRepo: a,
	}
}

// CreateTopic adds a new topic.
func (a *TopicService) CreateTopic(
	ctx context.Context,
	u *domain.CreateTopicRequest,
) (*domain.Topic, error) {
	createdTopic, err := a.topicRepo.CreateTopic(ctx, u)
	if err != nil {
		return nil, err
	}
	return createdTopic, nil
}

// GetTopic fetches a topic by ID.
func (a *TopicService) GetTopic(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Topic, error) {
	tracer := otel.Tracer("service.topic")
	ctxTrace, span := tracer.Start(ctx, "TopicService.GetTopic")
	defer span.End()

	topic, err := a.topicRepo.GetTopic(ctxTrace, id)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

// UpdateTopic updates name/email of an existing topic.
func (a *TopicService) UpdateTopic(
	ctx context.Context,
	id uuid.UUID,
	u *domain.Topic,
) (*domain.Topic, error) {

	existing, err := a.topicRepo.GetTopic(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrTopicNotFound
	}

	existing.Name = u.Name

	_, err = a.topicRepo.UpdateTopic(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteTopic removes a topic by ID.
func (a *TopicService) DeleteTopic(
	ctx context.Context,
	id uuid.UUID,
) error {

	topic, err := a.topicRepo.GetTopic(ctx, id)
	if err != nil {
		return err
	}
	if topic == nil {
		return domain.ErrTopicNotFound
	}

	err = a.topicRepo.DeleteTopic(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *TopicService) GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error) {
	topics, err := a.topicRepo.GetTopicList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return topics, nil
}

func (a *TopicService) GetTopicArticles(ctx context.Context, id uuid.UUID) ([]domain.Article, error) {
    articles, err := a.topicRepo.GetTopicArticles(ctx, id)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    return articles, nil
}
