-- +goose Up
-- +goose StatementBegin
INSERT INTO CUSTOMERS VALUES ('a1e2b3c4-d5f6-4a7b-8c9d-0e1f2a3b4c5d', 'Alice Johnson', '0x1234567890abcdef1234567890abcdef12345678', '123 Maple St, Springfield', '6282McNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcxY', '\x6ec4186b9010827e43035715e4ff9136');
INSERT INTO CUSTOMERS VALUES ('b2f3c4d5-e6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Bob Smith', '0x2345678901abcdef2345678901abcdef23456789', '456 Oak Ave, Metropolis', '7a9b3cNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcxZ', '\x7ec4186b9010827e43035715e4ff9137');
INSERT INTO CUSTOMERS VALUES ('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'Charlie Brown', '0x3456789012abcdef3456789012abcdef34567890', '789 Pine Rd, Gotham', '8b2c4dNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyA', '\x8ec4186b9010827e43035715e4ff9138');
INSERT INTO CUSTOMERS VALUES ('f2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c', 'Laura Palmer', '0x4567890123abcdef4567890123abcdef45678901', '101 Twin Peaks, Washington', '9c3d5eNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyB', '\x9ec4186b9010827e43035715e4ff9139');
INSERT INTO CUSTOMERS VALUES ('a3b4c5d6-e7f8-4a9b-0c1d-2e3f4a5b6c7d', 'Mallory Archer', '0x5678901234abcdef5678901234abcdef56789012', '212 ISIS Hq, New York', 'ad4e6fNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyC', '\xaec4186b9010827e43035715e4ff913a');
INSERT INTO CUSTOMERS VALUES ('d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'David Miller', '0x6789012345abcdef6789012345abcdef67890123', '321 Elm St, Central City', 'bd5f7gNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyD', '\xbec4186b9010827e43035715e4ff913b');
INSERT INTO CUSTOMERS VALUES ('e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'Eve Polastri', '0x7890123456abcdef7890123456abcdef78901234', '432 Birch Ln, London', 'ce6g8hNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyE', '\xcec4186b9010827e43035715e4ff913c');
INSERT INTO CUSTOMERS VALUES ('f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c', 'Frank Castle', '0x8901234567abcdef8901234567abcdef89012345', '543 Cedar Blvd, Hell''s Kitchen', 'df7h9iNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyF', '\xdec4186b9010827e43035715e4ff913d');
INSERT INTO CUSTOMERS VALUES ('a7b8c9d0-e1f2-4a3b-4c5d-6e7f8a9b0c1d', 'Grace Hopper', '0x9012345678abcdef9012345678abcdef90123456', '654 Walnut Dr, Arlington', 'eg8i0jNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyG', '\xeec4186b9010827e43035715e4ff913e');
INSERT INTO CUSTOMERS VALUES ('b8c9d0e1-f2a3-4b4c-5d6e-7f8a9b0c1d2e', 'Heidi Klum', '0xa012345678abcdefa012345678abcdefa0123456', '765 Cherry Ct, Beverly Hills', 'fh9j1kNtgqsNcaN979RTt8cCoN76pSLIPm04jE+wcyH', '\xfec4186b9010827e43035715e4ff913f');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM CUSTOMERS;
-- +goose StatementEnd
