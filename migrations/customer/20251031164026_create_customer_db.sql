-- +goose Up
-- +goose StatementBegin
CREATE TABLE CUSTOMERS (
  emp_id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  wallet_address TEXT NOT NULL,
  address TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE CUSTOMERS;
-- +goose StatementEnd
