package service

import (
	"context"
	"fmt"
	"zog-news/domain"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *domain.CreateArticleRequest) (*domain.Article, error)
	GetArticleList(ctx context.Context, filter *domain.ArticleFilter) ([]domain.Article, error)
	GetArticle(ctx context.Context, id uuid.UUID) (*domain.Article, error)
	UpdateArticle(ctx context.Context, id uuid.UUID, article *domain.Article) (*domain.Article, error)
	DeleteArticle(ctx context.Context, id uuid.UUID) error

    GetTopicsByArticleID(ctx context.Context, articleID uuid.UUID) ([]domain.Topic, error)
    AddTopicToArticle(ctx context.Context, articleID uuid.UUID, topicID string) error
    // AddTopicsToArticle(ctx context.Context, articleID uuid.UUID, topicIDs []string) error
    RemoveTopicFromArticle(ctx context.Context, articleID uuid.UUID, topicID string) error
    // GetArticlesByTopicID(ctx context.Context, topicID string) ([]domain.Article, error)
}

type ArticleService struct {
	articleRepo ArticleRepository
}

func NewArticleService(a ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: a,
	}
}

// CreateArticle adds a new article.
func (a *ArticleService) CreateArticle(
	ctx context.Context,
	u *domain.CreateArticleRequest,
) (*domain.Article, error) {
	createdArticle, err := a.articleRepo.CreateArticle(ctx, u)
	if err != nil {
		return nil, err
	}
	return createdArticle, nil
}

// GetArticle fetches a article by ID.
func (a *ArticleService) GetArticle(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Article, error) {
	tracer := otel.Tracer("service.article")
	ctxTrace, span := tracer.Start(ctx, "ArticleService.GetArticle")
	defer span.End()

	article, err := a.articleRepo.GetArticle(ctxTrace, id)
	if err != nil {
		return nil, err
	}
	return article, nil
}

// UpdateArticle updates title/content/author/status of an existing article.
func (a *ArticleService) UpdateArticle(
	ctx context.Context,
	id uuid.UUID,
	u *domain.Article,
) (*domain.Article, error) {

	existing, err := a.articleRepo.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrArticleNotFound
	}

	existing.Title = u.Title
	existing.Content = u.Content
	existing.Author = u.Author
	existing.Status = u.Status

	_, err = a.articleRepo.UpdateArticle(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteArticle removes a article by ID.
func (a *ArticleService) DeleteArticle(
	ctx context.Context,
	id uuid.UUID,
) error {

	article, err := a.articleRepo.GetArticle(ctx, id)
	if err != nil {
		return err
	}
	if article == nil {
		return domain.ErrArticleNotFound
	}

	err = a.articleRepo.DeleteArticle(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticleService) GetArticleList(ctx context.Context, filter *domain.ArticleFilter) ([]domain.Article, error) {
	articles, err := a.articleRepo.GetArticleList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return articles, nil
}

func (a *ArticleService) GetTopicsByArticleID(
    ctx context.Context,
    articleID uuid.UUID,
) ([]domain.Topic, error) {
    topics, err := a.articleRepo.GetTopicsByArticleID(ctx, articleID)
    if err != nil {
        return nil, err
    }
    return topics, nil
}

func (a *ArticleService) AddTopicToArticle(
    ctx context.Context,
    articleID uuid.UUID,
    topicID string,
) error {
    article, err := a.articleRepo.GetArticle(ctx, articleID)
    if err != nil {
        return err
    }
    if article == nil {
        return domain.ErrArticleNotFound
    }

    if err := article.AddTopicID(topicID); err != nil {
        return err
    }

    return a.articleRepo.AddTopicToArticle(ctx, articleID, topicID)
}

func (a *ArticleService) RemoveTopicFromArticle(
    ctx context.Context,
    articleID uuid.UUID,
    topicID string,
) error {
    article, err := a.articleRepo.GetArticle(ctx, articleID)
    if err != nil {
        return err
    }
    if article == nil {
        return domain.ErrArticleNotFound
    }
    if err := article.RemoveTopicID(topicID); err != nil {
        return err
    }
    return a.articleRepo.RemoveTopicFromArticle(ctx, articleID, topicID)
}
