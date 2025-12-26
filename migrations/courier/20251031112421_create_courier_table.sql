-- +goose Up
-- +goose StatementBegin
CREATE TABLE COURIERS (
  emp_id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  wallet_address TEXT NOT NULL,
  transport_type TEXT NOT NULL,
  is_active BOOLEAN NOT NULL,
  geolocation TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE COURIERS;
-- +goose StatementEnd
