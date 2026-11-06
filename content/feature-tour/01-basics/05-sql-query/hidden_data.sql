-- Same table, different rows: a level-5 edge case and more heroes.
CREATE TABLE heroes (id INTEGER PRIMARY KEY, name TEXT NOT NULL, level INTEGER NOT NULL);
INSERT INTO heroes VALUES (1, 'Ana', 6), (2, 'Bo', 5), (3, 'Cy', 1), (4, 'Di', 12);
