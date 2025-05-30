    CREATE TABLE IF NOT EXISTS urinary_clinic (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL,
    has_a_urinary_disease BOOL DEFAULT false,
    nervous_disease VARCHAR(255) DEFAULT '',
    relationship_with_diabetes BOOL DEFAULT false,
    comments VARCHAR(255) DEFAULT '',
  
    CONSTRAINT fk_patient FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE

);
