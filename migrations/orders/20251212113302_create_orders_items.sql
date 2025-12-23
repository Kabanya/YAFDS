-- +goose Up
-- +goose StatementBegin
CREATE TABLE ORDERS_ITEMS (
  empId UUID PRIMARY KEY,
  order_id UUID NOT NULL,
  restaurant_item_id UUID NOT NULL,
  price NUMERIC NOT NULL,
  quantity INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ORDERS_ITEMS;
-- +goose StatementEnd
