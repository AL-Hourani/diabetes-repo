CREATE TABLE IF NOT EXISTS patient_health_info (
    id SERIAL PRIMARY KEY,
    patient_id INT NOT NULL,
    bloodSugar VARCHAR(255) NOT NULL,
    hemoglobin VARCHAR(255) NOT NULL,
    bloodPressure VARCHAR(255) NOT NULL,
    sugarType VARCHAR(255) NOT NULL,
    diseaseDetection  VARCHAR(255) NOT NULL,
    otherDisease VARCHAR(255) NOT NULL,
    typeOfMedicine VARCHAR(255) NOT NULL,
    urineAcid VARCHAR(255) NOT NULL,
    cholesterol VARCHAR(255) NOT NULL,
    grease VARCHAR(255) NOT NULL,
    historyOfFamilyDisease VARCHAR(255) NOT NULL,
    createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_patient FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
