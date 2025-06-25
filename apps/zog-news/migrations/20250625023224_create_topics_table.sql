-- +goose Up
-- +goose StatementBegin
-- Using id over name primary key allows renaming topics easier
CREATE TABLE topics (
    id UUID PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE topics;
-- +goose StatementEnd
