package mocks

import (
	"context"
	"zog-news/domain"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

type ArticleService struct {
    mock.Mock
}

func (_m *ArticleService) CreateArticle(ctx context.Context, article *domain.CreateArticleRequest) (*domain.Article, error) {
	ret := _m.Called(ctx, article)

	var r0 *domain.Article
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Article)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ArticleService) GetArticle(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Article
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Article)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ArticleService) UpdateArticle(ctx context.Context, id uuid.UUID, u *domain.Article) (*domain.Article, error) {
	ret := _m.Called(ctx, id, u)

	var r0 *domain.Article
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.Article)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ArticleService) GetArticleList(ctx context.Context, filter *domain.ArticleFilter) ([]domain.Article, error) {
	ret := _m.Called(ctx, filter)

	var r0 []domain.Article
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]domain.Article)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ArticleService) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *ArticleService) AddTopicToArticle(ctx context.Context, articleID uuid.UUID, topicID string) error {
    ret := _m.Called(ctx, articleID, topicID)

    var r0 error
    if ret.Get(0) != nil {
        r0 = ret.Error(0)
    }

    return r0
}

func (_m *ArticleService) RemoveTopicFromArticle(ctx context.Context, articleID uuid.UUID, topicID string) error {
    ret := _m.Called(ctx, articleID, topicID)

    var r0 error
    if ret.Get(0) != nil {
        r0 = ret.Error(0)
    }

    return r0
}

func (_m *ArticleService) GetTopicsByArticleID(ctx context.Context, articleID uuid.UUID) ([]domain.Topic, error) {
    ret := _m.Called(ctx, articleID)

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
