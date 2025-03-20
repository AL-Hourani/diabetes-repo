CREATE TABLE IF NOT EXISTS health_info (
    id SERIAL PRIMARY KEY,
    age INT NOT NULL,
    gender VARCHAR(255) NOT NULL,
    diabetes_type VARCHAR(255) NOT NULL,
    patient_id INT NOT NULL,
    diagnosis_data VARCHAR(255) NOT NULL,
    current_blood_suger_level VARCHAR(255) NOT NULL,
    blood_suger_trends VARCHAR(255) NOT NULL,
    createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_patient FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
