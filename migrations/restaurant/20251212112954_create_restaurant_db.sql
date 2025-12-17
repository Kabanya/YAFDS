-- +goose Up
-- +goose StatementBegin
CREATE TABLE RESTAURANTS(
  empId UUID PRIMARY KEY,
  name TEXT NOT NULL,
  address_wallet TEXT NOT NULL,
  address TEXT NOT NULL,
  status BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE RESTAURANTS;
-- +goose StatementEnd
