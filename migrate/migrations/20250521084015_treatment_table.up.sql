CREATE TABLE treatments (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL,
    treatment_type VARCHAR(255),
    treatment_speed VARCHAR(255)  DEFAULT '' ,
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE
);