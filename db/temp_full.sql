-- create
CREATE TABLE CUSTOMERS (
  empId UUID PRIMARY KEY,
  name TEXT NOT NULL,
  wallet_address TEXT NOT NULL,
  address TEXT NOT NULL
);

CREATE TABLE COURIERS (
  empId UUID PRIMARY KEY,
  name TEXT NOT NULL,
  transport_type TEXT NOT NULL,
  is_active BOOLEAN NOT NULL,
  geolocation TEXT NOT NULL
);

CREATE TABLE ORDERS (
  empId UUID PRIMARY KEY,
  customer_id UUID NOT NULL,
  courrier_id UUID NOT NULL,
  started_at TIMESTAMP NOT NULL,
  finished_at TIMESTAMP NOT NULL,
  status TEXT NOT NULL
);

-- insert
INSERT INTO CUSTOMERS VALUES ('9484aea1-3ff0-4d6e-8925-dee68b9db7ff', 'Clark', 'solscard_228_322', 'Zalupkino-14');
INSERT INTO CUSTOMERS VALUES ('601be6b0-542b-439f-ac04-4a8a5364639b', 'Dave', 'solscard_322_288', 'Chumazovck-11');

INSERT INTO COURIERS VALUES ('788fbb30-3223-48ae-b85e-22b1ca457cf7', 'Ava', 'bike', true, 'hirosima-5');

INSERT INTO ORDERS VALUES ('55ec9cee-a9c5-46f1-b84b-2f84800e412e', '9484aea1-3ff0-4d6e-8925-dee68b9db7ff', '788fbb30-3223-48ae-b85e-22b1ca457cf7',
                            '2025-10-27 10:15:00', '2025-10-27 11:45:00', 'finished');
INSERT INTO ORDERS VALUES ('e6cfe900-2fc8-4b9e-ae76-244c7e04db4e', '9484aea1-3ff0-4d6e-8925-dee68b9db7ff', '788fbb30-3223-48ae-b85e-22b1ca457cf7',
                            '2025-10-27 11:30:00', '2025-10-27 11:50:00', 'failed');
INSERT INTO ORDERS VALUES ('e97950aa-ea84-4487-97be-5a055d602e51', '601be6b0-542b-439f-ac04-4a8a5364639b', '788fbb30-3223-48ae-b85e-22b1ca457cf7',
                            '2025-10-27 12:05:00', '2025-10-27 13:00:00', 'failed');
-- fetch
SELECT * FROM CUSTOMERS WHERE name = 'Dave';
SELECT * FROM ORDERS WHERE status = 'failed';

-- join
SELECT DISTINCT c.name FROM CUSTOMERS c JOIN ORDERS o  ON c.empId = o.customer_id WHERE o.status = 'failed'