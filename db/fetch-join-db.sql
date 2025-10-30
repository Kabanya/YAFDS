-- fetch
SELECT * FROM CUSTOMERS WHERE name = 'Dave';
SELECT * FROM ORDERS WHERE status = 'failed';

-- join
SELECT DISTINCT c.name FROM CUSTOMERS c JOIN ORDERS o  ON c.empId = o.customer_id WHERE o.status = 'failed'