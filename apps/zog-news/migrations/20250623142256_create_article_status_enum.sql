-- @docs https://www.postgresql.org/docs/current/datatype-enum.html
-- +goose Up
-- +goose StatementBegin
CREATE TYPE article_status AS ENUM ('draft', 'deleted', 'published');
-- +goose StatementEnd

-- @docs https://www.postgresql.org/docs/current/sql-droptype.html
-- +goose Down
-- +goose StatementBegin
DROP TYPE article_status;
-- +goose StatementEnd
