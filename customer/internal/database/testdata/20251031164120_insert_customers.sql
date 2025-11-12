-- +goose Up
-- +goose StatementBegin
INSERT INTO CUSTOMERS VALUES ('9484aea1-3ff0-4d6e-8925-dee68b9db7ff', 'Clark', 'solscard_228_322', 'Zalupkino-14');
INSERT INTO CUSTOMERS VALUES ('601be6b0-542b-439f-ac04-4a8a5364639b', 'Dave', 'solscard_322_288', 'Chumazovck-11');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM CUSTOMERS;
-- +goose StatementEnd
