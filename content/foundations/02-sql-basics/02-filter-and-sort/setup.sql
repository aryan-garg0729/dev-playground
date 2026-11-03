CREATE TABLE customers (
  id      INTEGER PRIMARY KEY,
  name    TEXT NOT NULL,
  country TEXT NOT NULL,
  email   TEXT NOT NULL
);
CREATE TABLE orders (
  id          INTEGER PRIMARY KEY,
  customer_id INTEGER NOT NULL REFERENCES customers (id),
  amount      INTEGER NOT NULL
);
INSERT INTO customers VALUES
  (1, 'Ada',   'UK',    'ada@example.com'),
  (2, 'Ravi',  'India', 'ravi@example.com'),
  (3, 'Mei',   'China', 'mei@example.com'),
  (4, 'Anaya', 'India', 'anaya@example.com');
INSERT INTO orders VALUES
  (1, 1, 30), (2, 2, 15), (3, 2, 25), (4, 4, 10), (5, 1, 5);
