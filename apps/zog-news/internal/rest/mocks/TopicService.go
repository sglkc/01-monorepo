package mocks

import (
	"context"
	"zog-news/domain"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

type TopicService struct {
    mock.Mock
}

func (_m *TopicService) CreateTopic(ctx context.Context, topic *domain.CreateTopicRequest) (*domain.Topic, error) {
	ret := _m.Called(ctx, topic)

	var r0 *domain.Topic
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Topic)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *TopicService) GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Topic
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Topic)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *TopicService) UpdateTopic(ctx context.Context, id uuid.UUID, u *domain.Topic) (*domain.Topic, error) {
	ret := _m.Called(ctx, id, u)

	var r0 *domain.Topic
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Topic)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *TopicService) GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error) {
	ret := _m.Called(ctx, filter)

	var r0 []domain.Topic
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]domain.Topic)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *TopicService) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}
