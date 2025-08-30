CREATE TABLE IF NOT EXISTS  treatment_drugs (
  id SERIAL PRIMARY KEY,
    treatment_id INT NOT NULL,
    drug_id INT NOT NULL,
    dosage_per_day VARCHAR(255),
    quantity INT DEFAULT 0,
    FOREIGN KEY (treatment_id) REFERENCES treatments(id) ON DELETE CASCADE,
    
);
