-- +goose Up
-- +goose StatementBegin
ALTER TABLE topics
    ALTER COLUMN deleted_at SET DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE topics
    ALTER COLUMN deleted_at SET DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd
