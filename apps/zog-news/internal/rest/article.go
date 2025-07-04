package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"zog-news/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article *domain.CreateArticleRequest) (*domain.Article, error)
	GetArticleList(ctx context.Context, filter *domain.ArticleFilter) ([]domain.Article, error)
	GetArticle(ctx context.Context, id uuid.UUID) (*domain.Article, error)
	UpdateArticle(ctx context.Context, id uuid.UUID, article *domain.Article) (*domain.Article, error)
	DeleteArticle(ctx context.Context, id uuid.UUID) error

	GetTopicsByArticleID(ctx context.Context, articleID uuid.UUID) ([]domain.Topic, error)
	AddTopicToArticle(ctx context.Context, articleID uuid.UUID, topicID string) error
	RemoveTopicFromArticle(ctx context.Context, articleID uuid.UUID, topicID string) error
}

type ArticleHandler struct {
	Service ArticleService
}

func NewArticleHandler(e *echo.Group, svc ArticleService) {
	handler := &ArticleHandler{
		Service: svc,
	}
	articleGroup := e.Group("/articles") // articles group

	articleGroup.GET("", handler.GetArticleList)
	articleGroup.GET("/:id", handler.GetArticle)
	articleGroup.POST("", handler.CreateArticle)
	articleGroup.PUT("/:id", handler.UpdateArticle)
	articleGroup.DELETE("/:id", handler.DeleteArticle)

	topicGroup := articleGroup.Group("/:id/topics") // topics group under articles
	topicGroup.GET("", handler.GetTopicsByArticleID)
	topicGroup.POST("/:topic_id", handler.AddTopicToArticle)
	topicGroup.DELETE("/:topic_id", handler.RemoveTopicFromArticle)
}

// GetArticleList retrieves a list of articles with optional filtering
//
//	@Summary		Get articles list
//	@Description	Get a paginated list of articles with optional filtering by search, status, and topic
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string										false	"Search in title and content"
//	@Param			status	query		string										false	"Filter by status"	Enums(draft,published,archived)
//	@Param			topic	query		string										false	"Filter by topic name"
//	@Success		200		{object}	domain.ResponseMultipleData[domain.Article]	"Successfully retrieved articles list"
//	@Failure		500		{object}	domain.ResponseMultipleData[domain.Empty]	"Internal server error"
//	@Router			/articles [get]
func (h *ArticleHandler) GetArticleList(c echo.Context) error {
	filter := new(domain.ArticleFilter)
	if err := c.Bind(filter); err != nil {
		fmt.Println(err)
	}

	ctx := c.Request().Context()
	articles, err := h.Service.GetArticleList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to list articles: " + err.Error(),
		})
	}
	if articles == nil {
		articles = []domain.Article{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Article]{
		Data:    articles,
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Successfully retrieve article list",
	})
}

// GetArticle retrieves a single article by ID
//
//	@Summary		Get article by ID
//	@Description	Get a single article by its unique identifier
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string										true	"Article ID"	format(uuid)
//	@Success		200	{object}	domain.ResponseSingleData[domain.Article]	"Successfully retrieved article"
//	@Failure		400	{object}	domain.ResponseSingleData[domain.Empty]		"Invalid article ID format"
//	@Failure		404	{object}	domain.ResponseSingleData[domain.Empty]		"Article not found"
//	@Failure		500	{object}	domain.ResponseSingleData[domain.Empty]		"Internal server error"
//	@Router			/articles/{id} [get]
func (h *ArticleHandler) GetArticle(c echo.Context) error {
	tracer := otel.Tracer("http.handler.article")
	ctx, span := tracer.Start(c.Request().Context(), "GetArticleHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid UUID")
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid article ID format",
		})
	}

	span.SetAttributes(attribute.String("article.id", id.String()))
	article, err := h.Service.GetArticle(ctx, id)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "not found")
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Article not found",
			})
		}

		span.SetStatus(codes.Error, "service error")
		fmt.Println("GetArticle error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get article: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Article]{
		Data:    *article,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved article",
	})
}

// CreateArticle creates a new article
//
//	@Summary		Create new article
//	@Description	Create a new article with the provided information
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			article	body		domain.CreateArticleRequest					true	"Article creation data"
//	@Success		201		{object}	domain.ResponseSingleData[domain.Article]	"Article successfully created"
//	@Failure		400		{object}	domain.ResponseSingleData[domain.Empty]		"Invalid request payload"
//	@Failure		500		{object}	domain.ResponseSingleData[domain.Empty]		"Internal server error"
//	@Router			/articles [post]
func (h *ArticleHandler) CreateArticle(c echo.Context) error {
	var article domain.CreateArticleRequest
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	createdArticle, err := h.Service.CreateArticle(ctx, &article)
	if err != nil {
		fmt.Println("CreateArticle error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create article: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.Article]{
		Data:    *createdArticle,
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Article successfully created",
	})
}

// UpdateArticle updates an existing article
//
//	@Summary		Update article
//	@Description	Update an existing article by ID with new information
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string										true	"Article ID"	format(uuid)
//	@Param			article	body		domain.UpdateArticleRequest					true	"Article update data"
//	@Success		200		{object}	domain.ResponseSingleData[domain.Article]	"Article successfully updated"
//	@Failure		400		{object}	domain.ResponseSingleData[domain.Empty]		"Invalid request payload or article ID"
//	@Failure		500		{object}	domain.ResponseSingleData[domain.Empty]		"Internal server error"
//	@Router			/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid article ID format",
		})
	}

	var article domain.Article
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	updatedArticle, err := h.Service.UpdateArticle(ctx, id, &article)
	if err != nil {
		fmt.Println("UpdateArticle error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update article: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Article]{
		Data:    *updatedArticle,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Article successfully updated",
	})
}

// DeleteArticle deletes an article by ID
//
//	@Summary		Delete article
//	@Description	Delete an article by its unique identifier (soft delete)
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Article ID"	format(uuid)
//	@Success		204	{object}	domain.ResponseSingleData[domain.Empty]	"Article successfully deleted"
//	@Failure		400	{object}	domain.ResponseSingleData[domain.Empty]	"Invalid article ID format"
//	@Failure		404	{object}	domain.ResponseSingleData[domain.Empty]	"Article not found"
//	@Failure		500	{object}	domain.ResponseSingleData[domain.Empty]	"Internal server error"
//	@Router			/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid article ID format",
		})
	}

	ctx := c.Request().Context()
	if err := h.Service.DeleteArticle(ctx, id); err != nil {
		fmt.Println("DeleteArticle error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete article: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Status:  "success",
		Message: "Article successfully deleted",
	})
}

// GetTopicsByArticleID retrieves topics associated with an article
//
//	@Summary		Get article topics
//	@Description	Get all topics associated with a specific article
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string										true	"Article ID"	format(uuid)
//	@Success		200	{object}	domain.ResponseMultipleData[domain.Topic]	"Successfully retrieved topics for article"
//	@Failure		400	{object}	domain.ResponseSingleData[domain.Empty]		"Invalid article ID format"
//	@Failure		500	{object}	domain.ResponseSingleData[domain.Empty]		"Internal server error"
//	@Router			/articles/{id}/topics [get]
func (h *ArticleHandler) GetTopicsByArticleID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid article ID format",
		})
	}
	ctx := c.Request().Context()
	topics, err := h.Service.GetTopicsByArticleID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get topics for article: " + err.Error(),
		})
	}
	if topics == nil {
		topics = []domain.Topic{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Topic]{
		Data:    topics,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved topics for article",
	})
}

// AddTopicToArticle adds a topic to an article
//
//	@Summary		Add topic to article
//	@Description	Associate a topic with an article
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string										true	"Article ID"	format(uuid)
//	@Param			topic_id	path		string										true	"Topic ID"		format(uuid)
//	@Success		200			{object}	domain.ResponseSingleData[domain.Article]	"Topic successfully added to article"
//	@Failure		400			{object}	domain.ResponseSingleData[domain.Empty]		"Invalid article ID or topic ID"
//	@Failure		404			{object}	domain.ResponseSingleData[domain.Empty]		"Article not found"
//	@Failure		500			{object}	domain.ResponseSingleData[domain.Empty]		"Internal server error"
//	@Router			/articles/{id}/topics/{topic_id} [post]
func (h *ArticleHandler) AddTopicToArticle(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid article ID format",
		})
	}
	topicID := c.Param("topic_id")
	if topicID == "" {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Topic ID is required",
		})
	}
	ctx := c.Request().Context()
	article, err := h.Service.GetArticle(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get article: " + err.Error(),
		})
	}
	if article == nil {
		return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Article not found",
		})
	}
	if err := h.Service.AddTopicToArticle(ctx, id, topicID); err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update article with new topic: " + err.Error(),
		})
	}
	// TODO: return article with topics??
	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Article]{
		Data:    *article,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Topic successfully added to article",
	})
}

// RemoveTopicFromArticle removes a topic from an article
//
//	@Summary		Remove topic from article
//	@Description	Disassociate a topic from an article
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string									true	"Article ID"	format(uuid)
//	@Param			topic_id	path		string									true	"Topic ID"		format(uuid)
//	@Success		204			{object}	domain.ResponseSingleData[domain.Empty]	"Topic successfully removed from article"
//	@Failure		400			{object}	domain.ResponseSingleData[domain.Empty]	"Invalid article ID or topic ID"
//	@Failure		404			{object}	domain.ResponseSingleData[domain.Empty]	"Article not found"
//	@Failure		500			{object}	domain.ResponseSingleData[domain.Empty]	"Internal server error"
//	@Router			/articles/{id}/topics/{topic_id} [delete]
func (h *ArticleHandler) RemoveTopicFromArticle(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid article ID format",
		})
	}
	topicID := c.Param("topic_id")
	if topicID == "" {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Topic ID is required",
		})
	}
	ctx := c.Request().Context()
	article, err := h.Service.GetArticle(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get article: " + err.Error(),
		})
	}
	if article == nil {
		return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusNotFound,
			Status:  "error",
			Message: "Article not found",
		})
	}
	if err := h.Service.RemoveTopicFromArticle(ctx, id, topicID); err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update article with removed topic: " + err.Error(),
		})
	}
	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Status:  "success",
		Message: "Article successfully deleted",
	})
}
