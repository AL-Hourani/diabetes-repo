package patient

import (
	"database/sql"
	"fmt"

	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/AL-Hourani/care-center/types"
)

type Store struct {
	db *sql.DB
}


func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
func  (s *Store) GetPatientByEmail(email string) (*types.Patient , error) {
	rows , err := s.db.Query("SELECT id,fullName,email,password,phone,date,id_number,isCompleted,center_id,createAt FROM patients WHERE email=$1",email)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	p := new(types.Patient)
	for rows.Next() {
		p , err = scanRowIntoPatient(rows)
		if err != nil {
			return nil , err
		}
	}

	if p.ID == 0 {
		return nil , fmt.Errorf("patient not found")
	}

	return p , nil
}



func  (s *Store) GetPatientById(id int) (*types.Patient , error) {
	rows , err := s.db.Query("SELECT id,fullName,email,password,phone,date,id_number,isCompleted,center_id,createAt FROM patients WHERE id=$1",id)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	p := new(types.Patient)
	for rows.Next() {
		p , err = scanRowIntoPatient(rows)
		if err != nil {
			return nil , err
		}
	}

	if p.ID == 0 {
		return nil , fmt.Errorf("patient not found")
	}

	return p , nil
}


func scanRowIntoPatient(rows *sql.Rows) (*types.Patient , error ){
	patient := new(types.Patient)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Password,
		&patient.Phone,
		&patient.Age,
		&patient.IDNumber,
		&patient.IsCompleted,
		&patient.CenterID,
		&patient.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return patient , nil
}




// 3
func (s *Store)	GreatePatient(patient types.Patient) error {
	_ , err := s.db.Exec("INSERT INTO patients (fullName , email , password ,phone , date , id_number , isCompleted , center_id , city) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)" , patient.FullName , patient.Email , patient.Password,patient.Phone, patient.Age ,patient.IDNumber , patient.IsCompleted , patient.CenterID , patient.City)
	if err  != nil {
		return err
	}

	return nil
}





// get all data form patient...

func (s *Store) GetPatientDetailsByID(patientID int) (*types.PatientDetails, error) {
	row := s.db.QueryRow(`SELECT id,fullName, email, phone, date, id_number,
		isCompleted, gender, wight, length_patient, address_patient,
		bloodSugar, hemoglobin, bloodPressure, sugarType, diseaseDetection,
		otherDisease, typeOfMedicine, urineAcid, cholesterol, grease,
		historyOfFamilyDisease
	
	FROM patients
	WHERE id=$1`,patientID)


	patient , err := scanRowIntoPatientDeatials(row)
		if err != nil {
			return nil , err
		}
	

	if patient.ID == 0 {
		return nil , fmt.Errorf("patient not found")
	}

	return patient , nil


}

func convertNullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}


func scanRowIntoPatientDeatials(rows *sql.Row) (*types.PatientDetails , error ){
	patient := new(types.PatientDetails)

	var (
		gender                sql.NullString
		weight                sql.NullString
		lengthPatient         sql.NullString
		addressPatient        sql.NullString
		bloodSugar            sql.NullString
		hemoglobin            sql.NullString
		bloodPressure         sql.NullString
		sugarType             sql.NullString
		diseaseDetection      sql.NullString
		otherDisease          sql.NullString
		typeOfMedicine        sql.NullString
		urineAcid             sql.NullString
		cholesterol           sql.NullString
		grease                sql.NullString
		historyOfFamilyDisease sql.NullString
	)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Phone,
		&patient.Date,
		&patient.IDNumber,
		&patient.IsCompleted,
		&gender,
		&weight,
		&lengthPatient,
		&addressPatient,
		&bloodSugar,
		&hemoglobin,
		&bloodPressure,
		&sugarType,
		&diseaseDetection,
		&otherDisease,
		&typeOfMedicine,
		&urineAcid,
		&cholesterol,
		&grease,
		&historyOfFamilyDisease,
	)
	
	if err  != nil {
		return nil , err
	}

	patient.Gender = convertNullStringToPointer(gender)
	patient.Weight = convertNullStringToPointer(weight)
	patient.LengthPatient = convertNullStringToPointer(lengthPatient)
	patient.AddressPatient = convertNullStringToPointer(addressPatient)
	patient.BloodSugar = convertNullStringToPointer(bloodSugar)
	patient.Hemoglobin = convertNullStringToPointer(hemoglobin)
	patient.BloodPressure = convertNullStringToPointer(bloodPressure)
	patient.SugarType = convertNullStringToPointer(sugarType)
	patient.DiseaseDetection = convertNullStringToPointer(diseaseDetection)
	patient.OtherDisease = convertNullStringToPointer(otherDisease)
	patient.TypeOfMedicine = convertNullStringToPointer(typeOfMedicine)
	patient.UrineAcid = convertNullStringToPointer(urineAcid)
	patient.Cholesterol = convertNullStringToPointer(cholesterol)
	patient.Grease = convertNullStringToPointer(grease)
	patient.HistoryOfFamilyDisease = convertNullStringToPointer(historyOfFamilyDisease)

	return patient , nil
}


func (s *Store) GetUserByEmail(email string) (*types.UserLoginData, error) {

	query := "SELECT * FROM login_serach  WHERE email = $1;"


    row := s.db.QueryRow(query, email)
	
	user := new(types.UserLoginData)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,

	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user , nil
}


func (s *Store) GetUserByID(id int) (*types.UserLoginData, error) {

	query := "SELECT * FROM login_serach  WHERE id = $1;"


    row := s.db.QueryRow(query, id)
	
	user := new(types.UserLoginData)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,

	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user , nil
}






func (s *Store) GetPatientProfile(id int)(*types.PatientProfile , error) {
	row := s.db.QueryRow(`SELECT id,fullName, email, phone, date, id_number,
	isCompleted, gender, wight, length_patient, address_patient,
	bloodSugar, hemoglobin, bloodPressure, sugarType, diseaseDetection,
	otherDisease, typeOfMedicine, urineAcid, cholesterol, grease,
	historyOfFamilyDisease,center_id,city

		FROM patients
		WHERE id=$1`,id)
	patient , err := scanRowIntoPatientProfile(row)
		if err != nil {
			return nil , err
		}
	

	if patient.ID == 0 {
		return nil , fmt.Errorf("patient not found")
	}

	return patient , nil
}


func scanRowIntoPatientProfile(rows *sql.Row) (*types.PatientProfile , error ){
	patient := new(types.PatientProfile)

	var (
		gender                sql.NullString
		weight                sql.NullString
		lengthPatient         sql.NullString
		addressPatient        sql.NullString
		bloodSugar            sql.NullString
		hemoglobin            sql.NullString
		bloodPressure         sql.NullString
		sugarType             sql.NullString
		diseaseDetection      sql.NullString
		otherDisease          sql.NullString
		typeOfMedicine        sql.NullString
		urineAcid             sql.NullString
		cholesterol           sql.NullString
		grease                sql.NullString
		historyOfFamilyDisease sql.NullString
	)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Phone,
		&patient.Age,
		&patient.IDNumber,
		&patient.IsCompleted,
		&gender,
		&weight,
		&lengthPatient,
		&addressPatient,
		&bloodSugar,
		&hemoglobin,
		&bloodPressure,
		&sugarType,
		&diseaseDetection,
		&otherDisease,
		&typeOfMedicine,
		&urineAcid,
		&cholesterol,
		&grease,
		&historyOfFamilyDisease,
		&patient.CenterID,
		&patient.City,
	)
	
	if err  != nil {
		return nil , err
	}

	patient.Gender = convertNullStringToPointer(gender)
	patient.Weight = convertNullStringToPointer(weight)
	patient.LengthPatient = convertNullStringToPointer(lengthPatient)
	patient.AddressPatient = convertNullStringToPointer(addressPatient)
	patient.BloodSugar = convertNullStringToPointer(bloodSugar)
	patient.Hemoglobin = convertNullStringToPointer(hemoglobin)
	patient.BloodPressure = convertNullStringToPointer(bloodPressure)
	patient.SugarType = convertNullStringToPointer(sugarType)
	patient.DiseaseDetection = convertNullStringToPointer(diseaseDetection)
	patient.OtherDisease = convertNullStringToPointer(otherDisease)
	patient.TypeOfMedicine = convertNullStringToPointer(typeOfMedicine)
	patient.UrineAcid = convertNullStringToPointer(urineAcid)
	patient.Cholesterol = convertNullStringToPointer(cholesterol)
	patient.Grease = convertNullStringToPointer(grease)
	patient.HistoryOfFamilyDisease = convertNullStringToPointer(historyOfFamilyDisease)

	return patient , nil
}



func (s *Store) UpdatePatientProfile(patientPayload types.ParientUpdatePayload)error {
	query := `UPDATE Patients
	SET 
	fullName = $1, 
	email = $2,
	phone = $3,
	date = $4,
	id_number = $5,
	city = $6
	WHERE id = $7`
	_, err := s.db.Exec(query,patientPayload.FullName,patientPayload.Email,patientPayload.Phone,patientPayload.Age,patientPayload.IDNumber,patientPayload.City,patientPayload.ID)

	if err != nil {
        return fmt.Errorf("error updating patient: %v", err)
		
    }

    return nil
}




// -------------------------------------------------------



func (s *Store) GetSugarTypeStats(centerID int) ([]*types.Statistics, error) {
	query := `
	SELECT
		sugarType,
		COUNT(*) AS total
	FROM patients
	WHERE center_id = $1
	GROUP BY sugarType
	ORDER BY total DESC;
	`

	rows, err := s.db.Query(query, centerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*types.Statistics

	for rows.Next() {
			var sugarType sql.NullString
			var total int

			err := rows.Scan(&sugarType, &total)
			if err != nil {
				return nil, err
			}

			stat := &types.Statistics{
				SugarType: "غير محدد", // default if NULL
				Total:     total,
			}
			if sugarType.Valid {
				stat.SugarType = sugarType.String
			}

			stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *Store) GetGenderCounts(centerID int) (int, int, error) {
	query := `
	SELECT gender, COUNT(*) FROM patients
	WHERE center_id = $1
	GROUP BY gender;
	`

	rows, err := s.db.Query(query, centerID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var maleCount, femaleCount int

	for rows.Next() {
		var gender sql.NullString
		var count int

		if err := rows.Scan(&gender, &count); err != nil {
			return 0, 0, err
		}

		if gender.Valid {
			switch gender.String {
			case "male", "ذكر":
				maleCount += count
			case "female", "أنثى":
				femaleCount += count
			}
		}
	}

	return maleCount, femaleCount, nil
}

func (s *Store) GetTotalPatientsInSystem() (int, error) {
	query := `SELECT COUNT(*) FROM patients;`

	var total int
	err := s.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (s *Store) GetSugarTypeAgeRangeStats(centerID int) ([]*types.SugarAgeRangeStat, error) {
query := `
SELECT
	sugarType,
	CASE
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) BETWEEN 0 AND 18 THEN '0-18'
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) BETWEEN 19 AND 35 THEN '19-35'
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) BETWEEN 36 AND 50 THEN '36-50'
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) > 50 THEN '51+'
		ELSE 'غير معروف'
	END AS age_range,
	COUNT(*) AS total
FROM patients
WHERE sugarType IS NOT NULL 
  AND date IS NOT NULL 
  AND date != ''
  AND date ~ '^\d{2}-\d{2}-\d{4}$'
  AND center_id = $1
GROUP BY sugarType, age_range
ORDER BY sugarType, total DESC;
`

	rows, err := s.db.Query(query, centerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*types.SugarAgeRangeStat

	for rows.Next() {
		var sugarType, ageRange string
		var total int

		if err := rows.Scan(&sugarType, &ageRange, &total); err != nil {
			return nil, err
		}

		stats = append(stats, &types.SugarAgeRangeStat{
			SugarType: sugarType,
			AgeRange:  ageRange,
			Total:     total,
		})
	}

	return stats, nil
}

func (s *Store) GetSugarTypeAgeRangeStatsAllSystem() ([]*types.SugarAgeRangeStat, error) {
	query := `
SELECT
	sugarType,
	CASE
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) BETWEEN 0 AND 18 THEN '0-18'
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) BETWEEN 19 AND 35 THEN '19-35'
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) BETWEEN 36 AND 50 THEN '36-50'
		WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, TO_DATE(date, 'DD-MM-YYYY'))) > 50 THEN '51+'
		ELSE 'غير معروف'
	END AS age_range,
	COUNT(*) AS total
FROM patients
WHERE sugarType IS NOT NULL 
  AND date IS NOT NULL 
  AND date != ''
  AND date ~ '^\d{2}-\d{2}-\d{4}$'
GROUP BY sugarType, age_range
ORDER BY sugarType, total DESC;

	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*types.SugarAgeRangeStat

	for rows.Next() {
		var sugarType, ageRange string
		var total int

		if err := rows.Scan(&sugarType, &ageRange, &total); err != nil {
			return nil, err
		}

		stats = append(stats, &types.SugarAgeRangeStat{
			SugarType: sugarType,
			AgeRange:  ageRange,
			Total:     total,
		})
	}

	return stats, nil
}


func (s *Store) GetBMIStats(centerID int) ([]*types.BMIStat, error) {
	query := `
		SELECT
			sugarType,
			CASE
				WHEN CAST(length_patient AS FLOAT) > 0 AND CAST(wight AS FLOAT) > 0 THEN
					CASE
						WHEN (CAST(wight AS FLOAT) / POWER(CAST(length_patient AS FLOAT)/100, 2)) < 18.5 THEN 'نحيف'
						WHEN (CAST(wight AS FLOAT) / POWER(CAST(length_patient AS FLOAT)/100, 2)) BETWEEN 18.5 AND 24.9 THEN 'طبيعي'
						WHEN (CAST(wight AS FLOAT) / POWER(CAST(length_patient AS FLOAT)/100, 2)) BETWEEN 25 AND 29.9 THEN 'زيادة وزن'
						ELSE 'سمنة'
					END
				ELSE 'غير معروف'
			END AS bmi_category,
			COUNT(*) AS total
		FROM patients
		WHERE wight IS NOT NULL
		  AND length_patient IS NOT NULL
		  AND wight != ''
		  AND length_patient != ''
		  AND center_id = $1
		GROUP BY sugarType, bmi_category
		ORDER BY sugarType, total DESC;
	`

	rows, err := s.db.Query(query, centerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*types.BMIStat

	for rows.Next() {
		var sugarType, bmiCategory string
		var total int

		if err := rows.Scan(&sugarType, &bmiCategory, &total); err != nil {
			return nil, err
		}

		stats = append(stats, &types.BMIStat{
			SugarType:   sugarType,
			BMICategory: bmiCategory,
			Total:       total,
		})
	}

	return stats, nil
}



func (s *Store) GetCityStats() ([]*types.CityStat, error) {
	query := `
		SELECT
			city,
			COUNT(*) AS total
		FROM patients
		WHERE city IS NOT NULL AND city != ''
		GROUP BY city
		ORDER BY total DESC;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*types.CityStat

	for rows.Next() {
		var city string
		var total int

		if err := rows.Scan(&city, &total); err != nil {
			return nil, err
		}

		stats = append(stats, &types.CityStat{
			City:  city,
			Total: total,
		})
	}

	return stats, nil
}








func (s *Store) GetUpdatePatientProfile(id int) (*types.GetPatientUpdateProfile , error) {
	row := s.db.QueryRow("SELECT fullName,email,phone,date,id_number,city FROM patients WHERE id=$1",id)
	patientProfile := new(types.GetPatientUpdateProfile)
	err := row.Scan(
		&patientProfile.FullName,
		&patientProfile.Email,
		&patientProfile.Phone,
		&patientProfile.Age,
		&patientProfile.IDNumber,
		&patientProfile.City,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("center not found")
		}
		return nil, err
	}

	return patientProfile , nil
}






//reset password 

func  (s *Store) GetUserByEmailRestPassword(email string) error {
	rows , err := s.db.Query("SELECT email FROM login_serach WHERE email=$1",email)
	if err != nil {
		return  err
	}
	defer rows.Close()

	for rows.Next() {
		_ , err = scanRowIntoUsertabele(rows)
		if err != nil {
			return  err
		}
	}


	return nil
}

func scanRowIntoUsertabele(rows *sql.Rows) (*types.Email , error ){
	email := new(types.Email)

	err := rows.Scan(
		&email.Email,
	)
	
	if err  != nil {
		return nil , err
	}

	return email , nil
}




func (s *Store) UpdatePasswordByEmail(email, newPassword string) error {
    // تشفير كلمة المرور الجديدة
 	hashedPassword , err := auth.HashPassword(newPassword)
	if err != nil {
	   return err
	}

    // تحديث كلمة المرور في قاعدة البيانات
    _, err = s.db.Exec(`
        UPDATE login_serach  
        SET password = $1 
        WHERE email = $2
    `, string(hashedPassword), email)

    return err
}