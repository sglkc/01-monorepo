-- +goose Up
-- +goose StatementBegin
ALTER TABLE article_topics ADD PRIMARY KEY (article_id, topic_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE article_topics DROP CONSTRAINT article_topics_pkey;
-- +goose StatementEnd
