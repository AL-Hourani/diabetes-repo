ALTER TABLE medication_requests
ADD COLUMN IF NOT EXISTS medication_id INT REFERENCES medications(id);
