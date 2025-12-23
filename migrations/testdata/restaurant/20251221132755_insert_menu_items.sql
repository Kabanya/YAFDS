-- +goose Up
-- +goose StatementBegin
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('21e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', '11e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Whopper', 5.99, 100, NULL, 'Classic flame-grilled beef burger');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('22f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', '12f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Pepperoni Pizza', 12.99, 50, NULL, 'Classic pepperoni with mozzarella');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('23d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', '13d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Salmon Roll', 8.50, 30, NULL, 'Fresh salmon with avocado');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('24e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', '14e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'Crunchy Taco', 2.49, 200, NULL, 'Seasoned beef in a crunchy shell');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('25f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', '15f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Spaghetti Carbonara', 14.00, 40, NULL, 'Creamy pasta with pancetta');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('26a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', '16a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Kung Pao Chicken', 11.50, 60, NULL, 'Spicy stir-fry with peanuts');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('27b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', '17b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 'Ribeye Steak', 25.00, 20, NULL, 'Grilled 12oz ribeye');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('28c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', '18c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'Quinoa Salad', 9.99, 45, NULL, 'Healthy quinoa with vegetables');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('29d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', '19d0e1f2-a3b4-4c5d-6e7f-8a9b0c1d2e3f', 'Chocolate Croissant', 3.50, 80, NULL, 'Flaky pastry with chocolate');
INSERT INTO RESTAURANT_MENU_ITEMS VALUES ('20e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', '10e1f2a3-b4c5-4d6e-7f8a-9b0c1d2e3f4a', 'Cappuccino', 4.25, 150, NULL, 'Rich espresso with steamed milk');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM RESTAURANT_MENU_ITEMS;
-- +goose StatementEnd
