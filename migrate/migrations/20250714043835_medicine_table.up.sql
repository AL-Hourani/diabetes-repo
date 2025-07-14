CREATE TABLE IF NOT EXISTS medications (

    id SERIAL PRIMARY KEY,
    name_arabic TEXT NOT NULL,                
    name_english TEXT NOT NULL,                
    medication_type TEXT NOT NULL,            
    dosage TEXT NOT NULL,                      
    expiration_date DATE NOT NULL,            
    quantity VARCHAR(255) NOT NULL,                 
    units_per_box INTEGER  NOT NULL             
);
