-- +goose Up
-- +goose StatementBegin
CREATE TABLE ORDERS (
  id UUID PRIMARY KEY,
  customer_id UUID NOT NULL,
  courier_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  status TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ORDERS;
-- +goose StatementEnd
