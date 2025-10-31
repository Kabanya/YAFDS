-- +goose up
CREATE TABLE COURIERS (
  empId UUID PRIMARY KEY,
  name TEXT NOT NULL,
  transport_type TEXT NOT NULL,
  is_active BOOLEAN NOT NULL,
  geolocation TEXT NOT NULL
);

-- +goose down
DROP TABLE COURIERS;