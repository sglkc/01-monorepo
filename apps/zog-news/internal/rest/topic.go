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

type TopicService interface {
	CreateTopic(ctx context.Context, topic *domain.CreateTopicRequest) (*domain.Topic, error)
	GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, topic *domain.Topic) (*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error

	GetTopicArticles(ctx context.Context, id uuid.UUID) ([]domain.Article, error)
}

type TopicHandler struct {
	Service TopicService
}

func NewTopicHandler(e *echo.Group, svc TopicService) {
	handler := &TopicHandler{
		Service: svc,
	}
	topicGroup := e.Group("/topics") // topics group

	topicGroup.GET("", handler.GetTopicList)
	topicGroup.GET("/:id", handler.GetTopic)
	topicGroup.POST("", handler.CreateTopic)
	topicGroup.PUT("/:id", handler.UpdateTopic)
	topicGroup.DELETE("/:id", handler.DeleteTopic)

	topicArticlesGroup := topicGroup.Group("/:id/articles")
	topicArticlesGroup.GET("", handler.GetTopicArticles)
}

// GetTopicList retrieves a list of topics with optional filtering
//
//	@Summary		Get topics list
//	@Description	Get a list of all topics with optional search filtering
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string										false	"Search in topic name"
//	@Success		200		{object}	domain.ResponseMultipleData[domain.Topic]	"Successfully retrieved topics list"
//	@Failure		500		{object}	domain.ResponseMultipleData[domain.Empty]	"Internal server error"
//	@Router			/topics [get]
func (h *TopicHandler) GetTopicList(c echo.Context) error {
	filter := new(domain.TopicFilter)
	if err := c.Bind(filter); err != nil {
		fmt.Println(err)
	}

	ctx := c.Request().Context()
	topics, err := h.Service.GetTopicList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to list topics: " + err.Error(),
		})
	}
	if topics == nil {
		topics = []domain.Topic{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Topic]{
		Data:    topics,
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Successfully retrieve topic list",
	})
}

// GetTopic retrieves a single topic by ID
//
//	@Summary		Get topic by ID
//	@Description	Get a single topic by its unique identifier
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Topic ID"	format(uuid)
//	@Success		200	{object}	domain.ResponseSingleData[domain.Topic]	"Successfully retrieved topic"
//	@Failure		400	{object}	domain.ResponseSingleData[domain.Empty]	"Invalid topic ID format"
//	@Failure		404	{object}	domain.ResponseSingleData[domain.Empty]	"Topic not found"
//	@Failure		500	{object}	domain.ResponseSingleData[domain.Empty]	"Internal server error"
//	@Router			/topics/{id} [get]
func (h *TopicHandler) GetTopic(c echo.Context) error {
	tracer := otel.Tracer("http.handler.topic")
	ctx, span := tracer.Start(c.Request().Context(), "GetTopicHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid UUID")
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid topic ID format",
		})
	}

	span.SetAttributes(attribute.String("topic.id", id.String()))
	topic, err := h.Service.GetTopic(ctx, id)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "not found")
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "Topic not found",
			})
		}

		span.SetStatus(codes.Error, "service error")
		fmt.Println("GetTopic error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get topic: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Topic]{
		Data:    *topic,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved topic",
	})
}

// CreateTopic creates a new topic
//
//	@Summary		Create new topic
//	@Description	Create a new topic with the provided information
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			topic	body		domain.CreateTopicRequest				true	"Topic creation data"
//	@Success		201		{object}	domain.ResponseSingleData[domain.Topic]	"Topic successfully created"
//	@Failure		400		{object}	domain.ResponseSingleData[domain.Empty]	"Invalid request payload"
//	@Failure		500		{object}	domain.ResponseSingleData[domain.Empty]	"Internal server error"
//	@Router			/topics [post]
func (h *TopicHandler) CreateTopic(c echo.Context) error {
	var topic domain.CreateTopicRequest
	if err := c.Bind(&topic); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	createdTopic, err := h.Service.CreateTopic(ctx, &topic)
	if err != nil {
		fmt.Println("CreateTopic error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create topic: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.Topic]{
		Data:    *createdTopic,
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Topic successfully created",
	})
}

// UpdateTopic updates an existing topic
//
//	@Summary		Update topic
//	@Description	Update an existing topic by ID with new information
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Topic ID"	format(uuid)
//	@Param			topic	body		domain.UpdateTopicRequest				true	"Topic update data"
//	@Success		200		{object}	domain.ResponseSingleData[domain.Topic]	"Topic successfully updated"
//	@Failure		400		{object}	domain.ResponseSingleData[domain.Empty]	"Invalid request payload or topic ID"
//	@Failure		500		{object}	domain.ResponseSingleData[domain.Empty]	"Internal server error"
//	@Router			/topics/{id} [put]
func (h *TopicHandler) UpdateTopic(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid topic ID format",
		})
	}

	var topic domain.Topic
	if err := c.Bind(&topic); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	updatedTopic, err := h.Service.UpdateTopic(ctx, id, &topic)
	if err != nil {
		fmt.Println("UpdateTopic error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update topic: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Topic]{
		Data:    *updatedTopic,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Topic successfully updated",
	})
}

// DeleteTopic deletes a topic by ID
//
//	@Summary		Delete topic
//	@Description	Delete a topic by its unique identifier (soft delete)
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Topic ID"	format(uuid)
//	@Success		204	{object}	domain.ResponseSingleData[domain.Empty]	"Topic successfully deleted"
//	@Failure		400	{object}	domain.ResponseSingleData[domain.Empty]	"Invalid topic ID format"
//	@Failure		404	{object}	domain.ResponseSingleData[domain.Empty]	"Topic not found"
//	@Failure		500	{object}	domain.ResponseSingleData[domain.Empty]	"Internal server error"
//	@Router			/topics/{id} [delete]
func (h *TopicHandler) DeleteTopic(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid topic ID format",
		})
	}

	ctx := c.Request().Context()
	if err := h.Service.DeleteTopic(ctx, id); err != nil {
		fmt.Println("DeleteTopic error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete topic: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Status:  "success",
		Message: "Topic successfully deleted",
	})
}

// GetTopicArticles retrieves articles associated with a topic
//
//	@Summary		Get topic articles
//	@Description	Get all articles associated with a specific topic
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string										true	"Topic ID"	format(uuid)
//	@Success		200	{object}	domain.ResponseMultipleData[domain.Article]	"Successfully retrieved articles for topic"
//	@Failure		400	{object}	domain.ResponseSingleData[domain.Empty]		"Invalid topic ID format"
//	@Failure		500	{object}	domain.ResponseMultipleData[domain.Empty]	"Internal server error"
//	@Router			/topics/{id}/articles [get]
func (h *TopicHandler) GetTopicArticles(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid topic ID format",
		})
	}

	ctx := c.Request().Context()
	articles, err := h.Service.GetTopicArticles(ctx, id)
	if err != nil {
		fmt.Println("GetTopicArticles error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get topic articles: " + err.Error(),
		})
	}
	if articles == nil {
		articles = []domain.Article{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Article]{
		Data:    articles,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved topic articles",
	})
}
