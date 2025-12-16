-- +goose Up
-- +goose StatementBegin
ALTER TABLE COURIERS ADD COLUMN password_hash TEXT;
ALTER TABLE COURIERS ADD COLUMN password_salt BYTEA;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE COURIERS DROP COLUMN password_hash;
ALTER TABLE COURIERS DROP COLUMN password_salt;
-- +goose StatementEnd
