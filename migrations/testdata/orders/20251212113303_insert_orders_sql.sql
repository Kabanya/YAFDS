-- +goose Up
-- +goose StatementBegin
INSERT INTO ORDERS VALUES ('31e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', 'a1e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', 'c1e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', '2025-12-21 10:00:00', '2025-12-21 10:30:00', 'delivered');
INSERT INTO ORDERS VALUES ('32f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', 'b2f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', 'c2f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', '2025-12-21 11:00:00', '2025-12-21 11:45:00', 'delivered');
INSERT INTO ORDERS VALUES ('33d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', '2025-12-21 12:00:00', '2025-12-21 12:15:00', 'cancelled');
INSERT INTO ORDERS VALUES ('34e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c', 'c4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', '2025-12-21 13:00:00', '2025-12-21 13:30:00', 'delivered');
INSERT INTO ORDERS VALUES ('35f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d', 'c5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', '2025-12-21 14:00:00', '2025-12-21 14:45:00', 'delivered');
INSERT INTO ORDERS VALUES ('36a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'd4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'c6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', '2025-12-21 15:00:00', '2025-12-21 15:30:00', 'delivered');
INSERT INTO ORDERS VALUES ('37b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 'e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'c7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', '2025-12-21 16:00:00', '2025-12-21 16:45:00', 'delivered');
INSERT INTO ORDERS VALUES ('38c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'c8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', '2025-12-21 17:00:00', '2025-12-21 17:30:00', 'delivered');
INSERT INTO ORDERS VALUES ('39d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', 'a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 'c9d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', '2025-12-21 18:00:00', '2025-12-21 18:45:00', 'delivered');
INSERT INTO ORDERS VALUES ('30e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', 'b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'c0e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', '2025-12-21 19:00:00', '2025-12-21 19:30:00', 'delivered');

INSERT INTO ORDERS_ITEMS VALUES ('41e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', '31e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', '21e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', 5.99, 2);
INSERT INTO ORDERS_ITEMS VALUES ('42f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', '32f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', '22f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', 12.99, 1);
INSERT INTO ORDERS_ITEMS VALUES ('43d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', '33d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', '23d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 8.50, 3);
INSERT INTO ORDERS_ITEMS VALUES ('44e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', '34e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', '24e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 2.49, 5);
INSERT INTO ORDERS_ITEMS VALUES ('45f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', '35f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', '25f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 14.00, 1);
INSERT INTO ORDERS_ITEMS VALUES ('46a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', '36a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', '26a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 11.50, 2);
INSERT INTO ORDERS_ITEMS VALUES ('47b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', '37b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', '27b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 25.00, 1);
INSERT INTO ORDERS_ITEMS VALUES ('48c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', '38c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', '28c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 9.99, 2);
INSERT INTO ORDERS_ITEMS VALUES ('49d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', '39d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', '29d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', 3.50, 4);
INSERT INTO ORDERS_ITEMS VALUES ('40e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', '30e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', '20e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', 4.25, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM ORDERS
-- +goose StatementEnd
