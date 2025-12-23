-- +goose Up
-- +goose StatementBegin
ALTER TABLE RESTAURANTS ADD COLUMN password_hash TEXT;
ALTER TABLE RESTAURANTS ADD COLUMN password_salt BYTEA;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE RESTAURANTS DROP COLUMN password_hash;
ALTER TABLE RESTAURANTS DROP COLUMN password_salt;
-- +goose StatementEnd
