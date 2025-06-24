-- +goose Up
-- +goose StatementBegin
-- https://stackoverflow.com/questions/20342717/postgresql-change-column-type-from-int-to-uuid
-- https://www.postgresql.org/docs/current/functions-uuid.html
-- https://stackoverflow.com/questions/41149554/default-for-column-xxxx-cannot-be-cast-automatically-to-type-boolean-in-postgr
ALTER TABLE articles
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id TYPE uuid USING (gen_random_uuid()),
    ALTER COLUMN id SET DEFAULT (gen_random_uuid());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Wallahua'lam
-- https://www.postgresql.org/docs/current/datatype-numeric.html
ALTER TABLE articles
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN id TYPE integer USING (nextval('articles_id_seq')),
    ALTER COLUMN id SET DEFAULT (nextval('articles_id_seq'));
-- +goose StatementEnd
