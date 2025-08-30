CREATE TABLE IF NOT EXISTS medication_requests (
    id SERIAL PRIMARY KEY,
    name_arabic TEXT NOT NULL,
    dosage TEXT NOT NULL,
    medication_type TEXT NOT NULL,
    requested_quantity INTEGER NOT NULL,
    center_id INT NOT NULL,
    requested_at DATE NOT NULL DEFAULT CURRENT_DATE,
    medication_id INT NOT NULL
);
