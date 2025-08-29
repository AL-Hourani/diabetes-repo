ALTER TABLE information
ADD COLUMN IF NOT EXISTS request_id INT UNIQUE REFERENCES medication_requests(id);
