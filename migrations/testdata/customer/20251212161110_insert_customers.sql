-- +goose Up
-- +goose StatementBegin
INSERT INTO CUSTOMERS VALUES ('a1e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Alice',  'solscard_101_201', 'Mainstreet-1', '6282McNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcxY', '\x6ec4186b9010827e43035715e4ff9136');
INSERT INTO CUSTOMERS VALUES ('b2f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Bob',    'solscard_102_202', 'Mainstreet-2', '7a9b3cNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcxZ', '\x7ec4186b9010827e43035715e4ff9137');
INSERT INTO CUSTOMERS VALUES ('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Charlie','solscard_103_203', 'Mainstreet-3', '8b2c4dNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyA', '\x8ec4186b9010827e43035715e4ff9138');
INSERT INTO CUSTOMERS VALUES ('f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c', 'Laura',  'solscard_112_212', 'Mainstreet-12', '9c3d5eNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyB', '\x9ec4186b9010827e43035715e4ff9139');
INSERT INTO CUSTOMERS VALUES ('a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d', 'Mallory','solscard_113_213', 'Mainstreet-13', 'ad4e6fNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyC', '\xaec4186b9010827e43035715e4ff913a');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM CUSTOMERS;
-- +goose StatementEnd
