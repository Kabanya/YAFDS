-- +goose Up
-- +goose StatementBegin
ALTER TABLE CUSTOMERS ADD COLUMN password_hash TEXT;
ALTER TABLE CUSTOMERS ADD COLUMN password_salt BYTEA;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE CUSTOMERS DROP COLUMN password_hash;
ALTER TABLE CUSTOMERS DROP COLUMN password_salt;
-- +goose StatementEnd
