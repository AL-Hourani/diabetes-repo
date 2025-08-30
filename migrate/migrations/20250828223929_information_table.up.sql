CREATE TABLE IF NOT EXISTS information (
    id SERIAL PRIMARY KEY,
    name_arabic TEXT NOT NULL,
    name_english TEXT NOT NULL,
    requested_quantity INTEGER NOT NULL,
    center_id INT NOT NULL,
    information_status  TEXT NOT NULL,
    createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    request_id INT NOT NULL
);
