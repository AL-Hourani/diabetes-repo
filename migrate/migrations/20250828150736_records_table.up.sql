CREATE TABLE IF NOT EXISTS records (
    id SERIAL PRIMARY KEY,
    name_arabic TEXT NOT NULL,
    dosage TEXT NOT NULL,
    medication_type TEXT NOT NULL,
    requested_quantity INTEGER NOT NULL,
    center_id INT NOT NULL,
    created_at  TEXT NOT NULL,
    approval_date TEXT NOT Null,
    record_status  TEXT NOT NULL,
    request_id INT NOT NULL
);
