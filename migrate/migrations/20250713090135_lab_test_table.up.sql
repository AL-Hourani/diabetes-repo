CREATE TABLE IF NOT EXISTS lab_tests (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    aliases TEXT NOT NULL,
    min_value REAL NOT NULL,
    max_value REAL NOT NULL,
    unit TEXT NOT NULL
);
