SELECT
    t.table_name,
    string_agg(c.column_name, ', ' ORDER BY c.ordinal_position) AS columns
FROM
    information_schema.tables t
JOIN
    information_schema.columns c ON t.table_name = c.table_name AND t.table_schema = c.table_schema
WHERE
    t.table_schema = 'public'
AND t.table_name NOT IN (
    'goose_restaurant_version',
    'goose_courier_version',
    'goose_orders_version',
    'goose_customers_version',
    'goose_customer_version'
)
GROUP BY
    t.table_name
ORDER BY
    t.table_name;
