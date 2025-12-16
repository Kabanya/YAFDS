-- +goose Up
-- +goose StatementBegin
CREATE TABLE RESTAURANT_MENU_ITEMS (
  order_item_id UUID PRIMARY KEY,
  restaurant_id UUID NOT NULL,
  name TEXT NOT NULL,
  price NUMERIC(10,2) NOT NULL,
  quantity INT NOT NULL,
  image BYTEA NULL,
  description TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE RESTAURANT_MENU_ITEMS;
-- +goose StatementEnd
