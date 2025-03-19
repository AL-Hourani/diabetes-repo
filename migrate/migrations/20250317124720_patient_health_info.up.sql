CREATE TABLE IF NOT EXISTS health_info (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `age` INT NOT NULL,
    `gender` VARCHAR(255)    NOT NULL,
    `diabetes_type` VARCHAR(255) NOT NULL,
    `patient_id` INT UNSIGNED NOT NULL,
    `diagnosis_data`  VARCHAR(255) NOT NULL,
    `current_blood_suger_level` VARCHAR(255) NOT NULL,
    `blood_suger_trends` VARCHAR(255) NOT NULL,
    `createAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY(id),
    FOREIGN KEY (`patient_id`) REFERENCES patients(`id`) ON DELETE CASCADE
);
