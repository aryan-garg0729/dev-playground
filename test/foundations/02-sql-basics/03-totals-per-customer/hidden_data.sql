-- Same tables, different rows: two customers share a name.
CREATE TABLE customers (id INTEGER PRIMARY KEY, name TEXT NOT NULL, country TEXT NOT NULL, email TEXT NOT NULL);
CREATE TABLE orders (id INTEGER PRIMARY KEY, customer_id INTEGER NOT NULL, amount INTEGER NOT NULL);
INSERT INTO customers VALUES (1, 'Sam', 'US', 'sam1@example.com'), (2, 'Sam', 'UK', 'sam2@example.com'), (3, 'Lee', 'KR', 'lee@example.com');
INSERT INTO orders VALUES (1, 1, 7), (2, 2, 11), (3, 2, 1), (4, 3, 100);
