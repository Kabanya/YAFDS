-- +goose Up
-- +goose StatementBegin
INSERT INTO RESTAURANTS VALUES ('11e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Burger King', '0x1111111111111111111111111111111111111111', '100 Fast Food Way, Springfield', true, 'hash1', '\x01');
INSERT INTO RESTAURANTS VALUES ('12f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Pizza Hut', '0x2222222222222222222222222222222222222222', '200 Pizza Plaza, Metropolis', true, 'hash2', '\x02');
INSERT INTO RESTAURANTS VALUES ('13d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Sushi Zen', '0x3333333333333333333333333333333333333333', '300 Zen Garden, Gotham', true, 'hash3', '\x03');
INSERT INTO RESTAURANTS VALUES ('14e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Taco Bell', '0x4444444444444444444444444444444444444444', '400 Taco Terrace, Washington', true, 'hash4', '\x04');
INSERT INTO RESTAURANTS VALUES ('15f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Pasta House', '0x5555555555555555555555555555555555555555', '500 Pasta Path, New York', true, 'hash5', '\x05');
INSERT INTO RESTAURANTS VALUES ('16a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Wok Express', '0x6666666666666666666666666666666666666666', '600 Wok Way, Central City', true, 'hash6', '\x06');
INSERT INTO RESTAURANTS VALUES ('17b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 'Steakhouse', '0x7777777777777777777777777777777777777777', '700 Steak St, London', true, 'hash7', '\x07');
INSERT INTO RESTAURANTS VALUES ('18c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'Vegan Delight', '0x8888888888888888888888888888888888888888', '800 Green Grove, Hell''s Kitchen', true, 'hash8', '\x08');
INSERT INTO RESTAURANTS VALUES ('19d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', 'Bakery & Co', '0x9999999999999999999999999999999999999999', '900 Baker Blvd, Arlington', true, 'hash9', '\x09');
INSERT INTO RESTAURANTS VALUES ('10e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', 'Coffee Shop', '0x0000000000000000000000000000000000000000', '10 Coffee Court, Beverly Hills', true, 'hash0', '\x00');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
