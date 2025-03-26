CREATE TABLE IF NOT EXISTS basic_patient_info (
    id SERIAL PRIMARY KEY,
    patient_id INT NOT NULL,
    gender VARCHAR(255) NOT NULL,
    wight VARCHAR(255) NOT NULL,
    length_patient  VARCHAR(255) NOT NULL,
    address_patient  VARCHAR(255) NOT NULL,
    createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_patient FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
