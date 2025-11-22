-- +goose Up
-- +goose StatementBegin
CREATE TABLE CUSTOMERS (
  empId UUID PRIMARY KEY,
  name TEXT NOT NULL,
  walletAddress TEXT NOT NULL,
  address TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE CUSTOMERS;
-- +goose StatementEnd
