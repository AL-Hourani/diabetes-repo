CREATE TABLE IF NOT EXISTS  treatment_drugs (
  id SERIAL PRIMARY KEY,
    treatment_id INT NOT NULL,
    drug_id INT NOT NULL,
    dosage_per_day VARCHAR(255),
    units VARCHAR(255),
    FOREIGN KEY (treatment_id) REFERENCES treatments(id) ON DELETE CASCADE,
    FOREIGN KEY (drug_id) REFERENCES drugs(id) ON DELETE CASCADE
);
