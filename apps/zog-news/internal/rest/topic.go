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
}

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
