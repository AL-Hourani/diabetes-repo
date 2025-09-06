CREATE TABLE IF NOT EXISTS patient_m (
    id SERIAL PRIMARY KEY,
    patient_id INT NOT NULL,
    gender VARCHAR(10),
    sugarType VARCHAR(255) ,
    historyOfFamilyDisease JSONB DEFAULT '[]',
    diseaseDetection  VARCHAR(255) DEFAULT '',
    CONSTRAINT fk_patient FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
