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
}

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

func (h *ArticleHandler) CreateArticle(c echo.Context) error {
	var article domain.CreateArticleRequest
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if err := c.Validate(&article); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Validation failed: " + err.Error(),
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

	if err := c.Validate(&article); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Validation failed: " + err.Error(),
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
