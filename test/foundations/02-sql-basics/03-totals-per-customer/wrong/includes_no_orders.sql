SELECT c.name, SUM(o.amount) AS total FROM customers c LEFT JOIN orders o ON o.customer_id = c.id GROUP BY c.id;
