-- +goose Up
-- +goose StatementBegin
INSERT INTO COURIERS VALUES ('788fbb30-3223-48ae-b85e-22b1ca457cf7', 'Ava', '0x123', 'bike', true, 'hirosima-5');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM COURIERS
-- +goose StatementEnd
