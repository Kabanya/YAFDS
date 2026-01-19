SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
AND table_name NOT IN (
    'goose_restaurant_version',
    'goose_courier_version',
    'goose_orders_version',
    'goose_customers_version',
    'goose_customer_version'
)
ORDER BY table_name;
