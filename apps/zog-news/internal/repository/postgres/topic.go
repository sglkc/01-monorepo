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

type TopicRepository struct {
	Conn *pgxpool.Pool
}

func NewTopicRepository(conn *pgxpool.Pool) *TopicRepository {
	return &TopicRepository{
		Conn: conn,
	}
}

func (a *TopicRepository) CreateTopic(ctx context.Context, topic *domain.CreateTopicRequest) (*domain.Topic, error) {

	query := `
		INSERT INTO topics (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id`

	var id uuid.UUID
	var err = a.Conn.QueryRow(ctx, query, topic.Name).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &domain.Topic{
		ID:    id.String(),
		Name:  topic.Name,
	}, nil
}

func (a *TopicRepository) GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error) {
	query := `
		SELECT
            a.id,
            a.name,
            a.created_at,
            a.updated_at
		FROM topics a
        WHERE a.deleted_at is NULL`

	var args []interface{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, `(a.name ILIKE $1)`)
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

	var topics []domain.Topic

	for rows.Next() {
		var topic domain.Topic
		err := rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

func (a *TopicRepository) GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error) {
	tracer := otel.Tracer("repo.topic")
	ctx, span := tracer.Start(ctx, "TopicRepository.GetTopic")
	defer span.End()

	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM topics
		WHERE id = $1 AND deleted_at IS NULL`

	span.SetAttributes(attribute.String("query.statement", query))
	span.SetAttributes(attribute.String("query.parameter", id.String()))
	row := a.Conn.QueryRow(ctx, query, id)

	var topic domain.Topic
	err := row.Scan(
		&topic.ID,
		&topic.Name,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &topic, nil
}

func (a *TopicRepository) UpdateTopic(ctx context.Context, id uuid.UUID, topic *domain.Topic) (*domain.Topic, error) {
	query := `
		UPDATE topics
		SET name = $1,
			updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL`

	_, err := a.Conn.Exec(ctx, query, topic.Name, id)
	if err != nil {
		return nil, err
	}

	updatedTopic, err := a.GetTopic(ctx, id)
	if err != nil {
		return nil, err
	}
	return updatedTopic, nil
}

func (a *TopicRepository) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE topics
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := a.Conn.Exec(ctx, query, id)
	return err
}

func (a *TopicRepository) GetTopicArticles(
    ctx context.Context,
    id uuid.UUID,
) ([]domain.Article, error) {
    query := `
        SELECT a.id, a.title, a.content, a.author, a.status, a.created_at, a.updated_at
        FROM articles a
        JOIN article_topics at ON a.id = at.article_id
        WHERE at.topic_id = $1 AND a.deleted_at IS NULL`
    rows, err := a.Conn.Query(ctx, query, id)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var articles []domain.Article

    for rows.Next() {
        var article domain.Article
        if err := rows.Scan(
            &article.ID,
            &article.Title,
            &article.Content,
            &article.Author,
            &article.Status,
            &article.CreatedAt,
            &article.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        articles = append(articles, article)
    }

    return articles, nil
}
