    CREATE TABLE IF NOT EXISTS heart_clinic (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL,
    has_a_heart_disease BOOL DEFAULT false,
    heart_disease VARCHAR(255) DEFAULT '',
    relationship_with_diabetes VARCHAR(255) DEFAULT '',
    comments VARCHAR(255) DEFAULT '',
  
    CONSTRAINT fk_patient FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE

);
