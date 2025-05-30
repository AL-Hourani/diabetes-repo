    CREATE TABLE IF NOT EXISTS eyes_clinic (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL,
    has_a_eye_disease BOOL DEFAULT false,
    in_kind_disease VARCHAR(255) DEFAULT '',
    relationship_with_diabetes BOOL DEFAULT false,
    comments VARCHAR(255) DEFAULT '',
  
    CONSTRAINT fk_patient FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE

);
