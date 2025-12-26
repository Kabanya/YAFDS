-- +goose Up
-- +goose StatementBegin
CREATE TABLE RESTAURANTS(
  emp_id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  wallet_address TEXT NOT NULL,
  address TEXT NOT NULL,
  status BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE RESTAURANTS;
-- +goose StatementEnd
