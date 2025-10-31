-- +goose up
CREATE TABLE CUSTOMERS (
  empId UUID PRIMARY KEY,
  name TEXT NOT NULL,
  wallet_address TEXT NOT NULL,
  address TEXT NOT NULL
);

-- +goose down
DROP TABLE CUSTOMERS;