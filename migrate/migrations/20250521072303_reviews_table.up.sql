CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    patient_id INT NOT NULL,
    address_patient VARCHAR(255) DEFAULT '',
    wight VARCHAR(255) DEFAULT '',
    length_patient  VARCHAR(255)  DEFAULT '',
    sugarType VARCHAR(255) DEFAULT '',
    otherDisease  VARCHAR(255) DEFAULT '',
    historyOfFamilyDisease JSONB DEFAULT '[]',
    diseaseDetection  VARCHAR(255) DEFAULT '',
    gender VARCHAR(255) DEFAULT '',
    hemoglobin VARCHAR(255)  DEFAULT '',
    grease VARCHAR(255) DEFAULT '',
    urineAcid VARCHAR(255) DEFAULT '',
    bloodPressure VARCHAR(255) DEFAULT '',
    cholesterol VARCHAR(255) DEFAULT '',
    LDL  VARCHAR(255) DEFAULT '',
    HDL  VARCHAR(255) DEFAULT '',
    creatine  VARCHAR(255) DEFAULT '',
    normal_clucose VARCHAR(255) DEFAULT '',
    clucose_after_meal  VARCHAR(255) DEFAULT '',
    triple_grease VARCHAR(255) DEFAULT '',
    hba1c VARCHAR(255) DEFAULT '',
    comments VARCHAR(255) DEFAULT '',
    date_review DATE DEFAULT CURRENT_DATE,
    CONSTRAINT fk_patient FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE

);



