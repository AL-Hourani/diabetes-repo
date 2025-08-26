-- حذف الحقول المطلوبة
ALTER TABLE patients 
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS historyOfFamilyDisease,
    DROP COLUMN IF EXISTS grease,
    DROP COLUMN IF EXISTS cholesterol,
    DROP COLUMN IF EXISTS urineAcid,
    DROP COLUMN IF EXISTS typeOfMedicine,
    DROP COLUMN IF EXISTS otherDisease,
    DROP COLUMN IF EXISTS diseaseDetection,
    DROP COLUMN IF EXISTS sugarType,
    DROP COLUMN IF EXISTS bloodPressure,
    DROP COLUMN IF EXISTS hemoglobin,
    DROP COLUMN IF EXISTS bloodSugar,
    DROP COLUMN IF EXISTS address_patient,
    DROP COLUMN IF EXISTS length_patient,
    DROP COLUMN IF EXISTS wight,
    DROP COLUMN IF EXISTS gender,
    DROP COLUMN IF EXISTS isCompleted;

-- إضافة الحقل الجديد
ALTER TABLE patients
    ADD COLUMN first_login BOOLEAN NOT NULL DEFAULT FALSE;
