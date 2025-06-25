-- +goose Up
-- +goose StatementBegin
ALTER TABLE topics
    ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE article_topics
    ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE topics
    DROP COLUMN created_at,
    DROP COLUMN updated_at,
    DROP COLUMN deleted_at;

ALTER TABLE article_topics
    DROP COLUMN created_at,
    DROP COLUMN updated_at,
    DROP COLUMN deleted_at;
-- +goose StatementEnd
