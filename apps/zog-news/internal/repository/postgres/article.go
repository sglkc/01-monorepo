package postgres

import (
	"context"
	"strings"
	"zog-news/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ArticleRepository struct {
	Conn *pgxpool.Pool
}

func NewArticleRepository(conn *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{
		Conn: conn,
	}
}

func (a *ArticleRepository) CreateArticle(ctx context.Context, article *domain.CreateArticleRequest) (*domain.Article, error) {

	query := `
		INSERT INTO articles (title, content, author, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id`

	var id uuid.UUID
	var err = a.Conn.QueryRow(ctx, query, article.Title, article.Content, article.Author, article.Status).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &domain.Article{
		ID:    id.String(),
		Title:  article.Title,
		Author: article.Author,
		Content:  article.Content,
		Status: article.Status,
	}, nil
}

func (a *ArticleRepository) GetArticleList(ctx context.Context, filter *domain.ArticleFilter) ([]domain.Article, error) {
	query := `
		SELECT
			a.id,
			a.title,
			a.content,
            a.created_at,
            a.updated_at
		FROM articles a
        WHERE a.deleted_at is NULL`

	var args []interface{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, `(a.title ILIKE $1 OR u.content ILIKE $1)`)
		args = append(args, "%"+filter.Search+"%")
	}

	if len(conditions) > 0 {
		query += strings.Join(conditions, " AND ")
	}
	rows, err := a.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []domain.Article

	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Author,
			&article.Status,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (a *ArticleRepository) GetArticle(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	tracer := otel.Tracer("repo.article")
	ctx, span := tracer.Start(ctx, "ArticleRepository.GetArticle")
	defer span.End()

	query := `
		SELECT
			id,
			title,
			content,
			author,
			status,
			created_at,
			updated_at
		FROM articles
		WHERE id = $1 AND deleted_at IS NULL`

	span.SetAttributes(attribute.String("query.statement", query))
	span.SetAttributes(attribute.String("query.parameter", id.String()))
	row := a.Conn.QueryRow(ctx, query, id)

	var article domain.Article
	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Author,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &article, nil
}

func (a *ArticleRepository) UpdateArticle(ctx context.Context, id uuid.UUID, article *domain.Article) (*domain.Article, error) {
	query := `
		UPDATE articles
		SET title = $1,
			content = $2,
			author = $3,
			status = $4,
			updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL`

	_, err := a.Conn.Exec(ctx, query, article.Title, article.Content, article.Author, article.Status, id)
	if err != nil {
		return nil, err
	}

	updatedArticle, err := a.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	return updatedArticle, nil
}

func (a *ArticleRepository) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE articles
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := a.Conn.Exec(ctx, query, id)
	return err
}
