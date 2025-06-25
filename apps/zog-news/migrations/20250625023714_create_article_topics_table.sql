-- +goose Up
-- +goose StatementBegin
-- https://kb.objectrocket.com/postgresql/postgresql-composite-primary-keys-629
CREATE TABLE article_topics (
    article_id UUID REFERENCES articles(id), -- ON UPDATE ON DELETE?
    topic_id UUID REFERENCES topics(id) -- ON UPDATE ON DELETE?
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE article_topics;
-- +goose StatementEnd
