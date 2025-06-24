-- +goose Up
-- +goose StatementBegin
ALTER TABLE articles ADD COLUMN deleted_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE articles DROP COLUMN deleted_at;
-- +goose StatementEnd
