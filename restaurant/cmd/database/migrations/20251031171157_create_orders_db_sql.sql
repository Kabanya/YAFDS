-- +goose Up
-- +goose StatementBegin
CREATE TABLE ORDERS (
  empId UUID PRIMARY KEY,
  customer_id UUID NOT NULL,
  courrier_id UUID NOT NULL,
  started_at TIMESTAMP NOT NULL,
  finished_at TIMESTAMP NOT NULL,
  status TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ORDERS;
-- +goose StatementEnd
