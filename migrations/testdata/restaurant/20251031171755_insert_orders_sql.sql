-- +goose Up
-- +goose StatementBegin
INSERT INTO ORDERS VALUES ('55ec9cee-a9c5-46f1-b84b-2f84800e412e', '9484aea1-3ff0-4d6e-8925-dee68b9db7ff', '788fbb30-3223-48ae-b85e-22b1ca457cf7',
                            '2025-10-27 10:15:00', '2025-10-27 11:45:00', 'finished');
INSERT INTO ORDERS VALUES ('e6cfe900-2fc8-4b9e-ae76-244c7e04db4e', '9484aea1-3ff0-4d6e-8925-dee68b9db7ff', '788fbb30-3223-48ae-b85e-22b1ca457cf7',
                            '2025-10-27 11:30:00', '2025-10-27 11:50:00', 'failed');
INSERT INTO ORDERS VALUES ('e97950aa-ea84-4487-97be-5a055d602e51', '601be6b0-542b-439f-ac04-4a8a5364639b', '788fbb30-3223-48ae-b85e-22b1ca457cf7',
                            '2025-10-27 12:05:00', '2025-10-27 13:00:00', 'failed');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM ORDERS
-- +goose StatementEnd
