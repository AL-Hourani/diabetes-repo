package patient

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
    "log"
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
	rows , err := s.db.Query("SELECT id,fullName,email,password,phone,date,id_number,center_id,createAt,first_login FROM patients WHERE email=$1",email)
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
	rows , err := s.db.Query("SELECT id,fullName,email,password,phone,date,id_number,center_id,createAt,first_login FROM patients WHERE id=$1",id)
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
		&patient.CenterID,
		&patient.CreateAt,
		&patient.FirstLogin,
	)
	
	if err  != nil {
		return nil , err
	}

	return patient , nil
}



func (s *Store) GetLoginByID(id int) (*types.Login, error) {
 
    row := s.db.QueryRow(`
        SELECT id, email, password , role
        FROM login_serach
        WHERE id=$1
    `, id)

    l := new(types.Login)
    err := row.Scan(&l.ID, &l.Email, &l.Password , &l.Role)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }

    return l, nil
}



// 3
// func (s *Store)	GreatePatient(patient types.Patient) (int , error ) {
// 	_ , err := s.db.Exec("INSERT INTO patients (fullName , email , password ,phone , date , id_number , isCompleted , center_id , city) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)" , patient.FullName , patient.Email , patient.Password,patient.Phone, patient.Age ,patient.IDNumber , patient.IsCompleted , patient.CenterID , patient.City)
// 	if err  != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *Store) CreatePatient(patient types.Patient) (int, error) {
    var id int
    err := s.db.QueryRow(`
        INSERT INTO patients 
        (fullName, email, password, phone, date, id_number, center_id) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `, patient.FullName, patient.Email, patient.Password, patient.Phone, patient.Age, patient.IDNumber, patient.CenterID).Scan(&id)

    if err != nil {
        return 0, err
    }

    return id, nil
}

func (s *Store) CreatePatientM(patientM types.PatientM) error {
	historyJSON, errM := json.Marshal(patientM.HistoryOfFamilyDisease)
	if errM != nil {
		log.Println("خطأ في تحويل historyOfFamilyDisease إلى JSON:", errM)
		return errM
	}
     _ ,err := s.db.Exec(`
        INSERT INTO patient_m
        (patient_id, gender, sugarType, historyOfFamilyDisease, diseaseDetection)
        VALUES ($1, $2, $3, $4, $5)
    `,
        patientM.PatientID,
        patientM.Gender,
        patientM.SugarType,
        string(historyJSON),
        patientM.DiseaseDetection,
    )

    if err != nil {
        return err
    }

    return nil
}




func (s *Store) SetFirstLoginTrue(patientID int) error {
    query := `
        UPDATE patients
        SET first_login = FALSE
        WHERE id = $1
    `
    _, err := s.db.Exec(query, patientID)
    return err
}



// get all data form patient...

// func (s *Store) GetPatientDetailsByID(patientID int) (*types.PatientDetails, error) {
// 	row := s.db.QueryRow(`SELECT id,fullName, email, phone, date, id_number,
// 		isCompleted, gender, wight, length_patient, address_patient,
// 		bloodSugar, hemoglobin, bloodPressure, sugarType, diseaseDetection,
// 		otherDisease, typeOfMedicine, urineAcid, cholesterol, grease,
// 		historyOfFamilyDisease
	
// 	FROM patients
// 	WHERE id=$1`,patientID)


// 	patient , err := scanRowIntoPatientDeatials(row)
// 		if err != nil {
// 			return nil , err
// 		}
	

// 	if patient.ID == 0 {
// 		return nil , fmt.Errorf("patient not found")
// 	}

// 	return patient , nil


// }

// func convertNullStringToPointer(ns sql.NullString) *string {
// 	if ns.Valid {
// 		return &ns.String
// 	}
// 	return nil
// }


// func scanRowIntoPatientDeatials(rows *sql.Row) (*types.PatientDetails , error ){
// 	patient := new(types.PatientDetails)

// 	var (
// 		gender                sql.NullString
// 		weight                sql.NullString
// 		lengthPatient         sql.NullString
// 		addressPatient        sql.NullString
// 		bloodSugar            sql.NullString
// 		hemoglobin            sql.NullString
// 		bloodPressure         sql.NullString
// 		sugarType             sql.NullString
// 		diseaseDetection      sql.NullString
// 		otherDisease          sql.NullString
// 		typeOfMedicine        sql.NullString
// 		urineAcid             sql.NullString
// 		cholesterol           sql.NullString
// 		grease                sql.NullString
// 		historyOfFamilyDisease sql.NullString
// 	)

// 	err := rows.Scan(
// 		&patient.ID,
// 		&patient.FullName,
// 		&patient.Email,
// 		&patient.Phone,
// 		&patient.Date,
// 		&patient.IDNumber,
// 		&patient.IsCompleted,
// 		&gender,
// 		&weight,
// 		&lengthPatient,
// 		&addressPatient,
// 		&bloodSugar,
// 		&hemoglobin,
// 		&bloodPressure,
// 		&sugarType,
// 		&diseaseDetection,
// 		&otherDisease,
// 		&typeOfMedicine,
// 		&urineAcid,
// 		&cholesterol,
// 		&grease,
// 		&historyOfFamilyDisease,
// 	)
	
// 	if err  != nil {
// 		return nil , err
// 	}

// 	patient.Gender = convertNullStringToPointer(gender)
// 	patient.Weight = convertNullStringToPointer(weight)
// 	patient.LengthPatient = convertNullStringToPointer(lengthPatient)
// 	patient.AddressPatient = convertNullStringToPointer(addressPatient)
// 	patient.BloodSugar = convertNullStringToPointer(bloodSugar)
// 	patient.Hemoglobin = convertNullStringToPointer(hemoglobin)
// 	patient.BloodPressure = convertNullStringToPointer(bloodPressure)
// 	patient.SugarType = convertNullStringToPointer(sugarType)
// 	patient.DiseaseDetection = convertNullStringToPointer(diseaseDetection)
// 	patient.OtherDisease = convertNullStringToPointer(otherDisease)
// 	patient.TypeOfMedicine = convertNullStringToPointer(typeOfMedicine)
// 	patient.UrineAcid = convertNullStringToPointer(urineAcid)
// 	patient.Cholesterol = convertNullStringToPointer(cholesterol)
// 	patient.Grease = convertNullStringToPointer(grease)
// 	patient.HistoryOfFamilyDisease = convertNullStringToPointer(historyOfFamilyDisease)

// 	return patient , nil
// }


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






// func (s *Store) GetPatientProfile(id int)(*types.PatientProfile , error) {
// 	row := s.db.QueryRow(`SELECT id,fullName, email, phone, date, id_number,
// 	isCompleted, gender, wight, length_patient, address_patient,
// 	bloodSugar, hemoglobin, bloodPressure, sugarType, diseaseDetection,
// 	otherDisease, typeOfMedicine, urineAcid, cholesterol, grease,
// 	historyOfFamilyDisease,center_id,city

// 		FROM patients
// 		WHERE id=$1`,id)
// 	patient , err := scanRowIntoPatientProfile(row)
// 		if err != nil {
// 			return nil , err
// 		}
	

// 	if patient.ID == 0 {
// 		return nil , fmt.Errorf("patient not found")
// 	}

// 	return patient , nil
// }


// func scanRowIntoPatientProfile(rows *sql.Row) (*types.PatientProfile , error ){
// 	patient := new(types.PatientProfile)

// 	var (
// 		gender                sql.NullString
// 		weight                sql.NullString
// 		lengthPatient         sql.NullString
// 		addressPatient        sql.NullString
// 		bloodSugar            sql.NullString
// 		hemoglobin            sql.NullString
// 		bloodPressure         sql.NullString
// 		sugarType             sql.NullString
// 		diseaseDetection      sql.NullString
// 		otherDisease          sql.NullString
// 		typeOfMedicine        sql.NullString
// 		urineAcid             sql.NullString
// 		cholesterol           sql.NullString
// 		grease                sql.NullString
// 		historyOfFamilyDisease sql.NullString
// 	)

// 	err := rows.Scan(
// 		&patient.ID,
// 		&patient.FullName,
// 		&patient.Email,
// 		&patient.Phone,
// 		&patient.Age,
// 		&patient.IDNumber,
// 		&patient.IsCompleted,
// 		&gender,
// 		&weight,
// 		&lengthPatient,
// 		&addressPatient,
// 		&bloodSugar,
// 		&hemoglobin,
// 		&bloodPressure,
// 		&sugarType,
// 		&diseaseDetection,
// 		&otherDisease,
// 		&typeOfMedicine,
// 		&urineAcid,
// 		&cholesterol,
// 		&grease,
// 		&historyOfFamilyDisease,
// 		&patient.CenterID,
// 		&patient.City,
// 	)
	
// 	if err  != nil {
// 		return nil , err
// 	}

// 	patient.Gender = convertNullStringToPointer(gender)
// 	patient.Weight = convertNullStringToPointer(weight)
// 	patient.LengthPatient = convertNullStringToPointer(lengthPatient)
// 	patient.AddressPatient = convertNullStringToPointer(addressPatient)
// 	patient.BloodSugar = convertNullStringToPointer(bloodSugar)
// 	patient.Hemoglobin = convertNullStringToPointer(hemoglobin)
// 	patient.BloodPressure = convertNullStringToPointer(bloodPressure)
// 	patient.SugarType = convertNullStringToPointer(sugarType)
// 	patient.DiseaseDetection = convertNullStringToPointer(diseaseDetection)
// 	patient.OtherDisease = convertNullStringToPointer(otherDisease)
// 	patient.TypeOfMedicine = convertNullStringToPointer(typeOfMedicine)
// 	patient.UrineAcid = convertNullStringToPointer(urineAcid)
// 	patient.Cholesterol = convertNullStringToPointer(cholesterol)
// 	patient.Grease = convertNullStringToPointer(grease)
// 	patient.HistoryOfFamilyDisease = convertNullStringToPointer(historyOfFamilyDisease)

// 	return patient , nil
// }



func (s *Store) UpdatePatientProfile(patientPayload types.ParientUpdatePayload , id int)error {
	query := `UPDATE Patients
	SET 
	fullName = $1, 
	email = $2,
	phone = $3,
	date = $4,
	id_number = $5
	WHERE id = $6`
	_, err := s.db.Exec(query,patientPayload.FullName,patientPayload.Email,patientPayload.Phone,patientPayload.Age,patientPayload.IDNumber,id)

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
	row := s.db.QueryRow("SELECT fullName,email,phone,date,id_number FROM patients WHERE id=$1",id)
	patientProfile := new(types.GetPatientUpdateProfile)
	err := row.Scan(
		&patientProfile.FullName,
		&patientProfile.Email,
		&patientProfile.Phone,
		&patientProfile.Age,
		&patientProfile.IDNumber,
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
  
 	hashedPassword , err := auth.HashPassword(newPassword)
	if err != nil {
	   return err
	}

 
    _, err = s.db.Exec(`
        UPDATE login_serach  
        SET password = $1 
        WHERE email = $2
    `, string(hashedPassword), email)

    if err != nil {
        return fmt.Errorf("failed to update login_serach: %w", err)
    }

	_, err = s.db.Exec(`
		UPDATE patients
		SET password = $1
		WHERE email = $2
	`,  string(hashedPassword), email)
   if err != nil {
        return fmt.Errorf("failed to update patients: %w", err)
    }
	
    return err
}









// app patirnt..............
func (s *Store) GetReviewsByPatientID(patientID int) ([]types.ReviewResponseForPatient, error) {
	var reviews []types.ReviewResponseForPatient

	rows, err := s.db.Query(`
		SELECT 
			id, address_patient, wight, length_patient, otherDisease,
			hemoglobin, grease,
			urineAcid, bloodPressure, cholesterol, LDL, HDL, creatine, normal_clucose,
			clucose_after_meal, triple_grease, hba1c ,date_review 
		FROM reviews 
		WHERE patient_id = $1
		ORDER BY date_review DESC
	`, patientID)
	if err != nil {
		return nil, fmt.Errorf("query reviews failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var review types.ReviewResponseForPatient
		
		err := rows.Scan(
			&review.ID, &review.Address, &review.Weight, &review.LengthPatient, &review.OtherDisease,
			&review.Hemoglobin, &review.Grease,
			&review.UrineAcid, &review.BloodPressure, &review.Cholesterol, &review.LDL, &review.HDL,
			&review.Creatine, &review.NormalGlocose, &review.GlocoseAfterMeal, &review.TripleGrease,
			&review.Hba1c,&review.DateReview,
		)
		if err != nil {
			continue // يمكن تسجيل الخطأ
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}





func (s *Store) GetLastReviewByPatientID(patientID int) (*types.ReviewResponseForPatient, error) {
	var review types.ReviewResponseForPatient

	err := s.db.QueryRow(`
		SELECT 
			id, address_patient, wight, length_patient, otherDisease,
			hemoglobin, grease,
			urineAcid, bloodPressure, cholesterol, LDL, HDL, creatine, normal_clucose,
			clucose_after_meal, triple_grease, hba1c, date_review 
		FROM reviews 
		WHERE patient_id = $1
		ORDER BY date_review DESC
		LIMIT 1
	`, patientID).Scan(
		&review.ID, &review.Address, &review.Weight, &review.LengthPatient, &review.OtherDisease,
		&review.Hemoglobin, &review.Grease,
		&review.UrineAcid, &review.BloodPressure, &review.Cholesterol, &review.LDL, &review.HDL,
		&review.Creatine, &review.NormalGlocose, &review.GlocoseAfterMeal, &review.TripleGrease,
		&review.Hba1c, &review.DateReview,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // لا توجد مراجعات
		}
		return nil, fmt.Errorf("query last review failed: %w", err)
	}

	return &review, nil
}


func (s *Store) GetTreatmentTypeByReviewID(reviewID int) (json.RawMessage, error) {
    var treatmentType json.RawMessage

    err := s.db.QueryRow(`
        SELECT treatment_type
        FROM treatments
        WHERE review_id = $1
        LIMIT 1
    `, reviewID).Scan(&treatmentType)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // لا توجد بيانات للعلاج
        }
        return nil, fmt.Errorf("query treatment_type failed: %w", err)
    }

    return treatmentType, nil
}

















//get notifications .....

func (s *Store) GetNotificationsByUserID(userID int) ([]types.NotificationTwo, error) {
    rows, err := s.db.Query(`
        SELECT id, sender_id, receiver_id, message, is_read, created_at
        FROM notifications
        WHERE receiver_id = $1
        ORDER BY created_at DESC
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifs []types.NotificationTwo
    for rows.Next() {
        var (
            n           types.NotificationTwo
            createdTime time.Time
        )

        err := rows.Scan(&n.ID, &n.SenderID, &n.ReceiverID, &n.Message, &n.IsRead, &createdTime)
        if err != nil {
            return nil, err
        }

        n.CreatedAt = FormatRelativeTime(createdTime) 
        notifs = append(notifs, n)
    }
    return notifs, nil
}




func  FormatRelativeTime(t time.Time) string {
    duration := time.Since(t)

    seconds := int(duration.Seconds())
    minutes := int(duration.Minutes())
    hours := int(duration.Hours())
    days := hours / 24

    switch {
    case seconds < 60:
        if seconds <= 1 {
            return "منذ ثانية واحدة"
        } else if seconds <= 2 {
            return "منذ ثانيتين"
        } else if seconds < 11 {
            return fmt.Sprintf("منذ %d ثوانٍ", seconds)
        } else {
            return fmt.Sprintf("منذ %d ثانية", seconds)
        }

    case minutes < 60:
        if minutes == 1 {
            return "منذ دقيقة واحدة"
        } else if minutes == 2 {
            return "منذ دقيقتين"
        } else if minutes < 11 {
            return fmt.Sprintf("منذ %d دقائق", minutes)
        } else {
            return fmt.Sprintf("منذ %d دقيقة", minutes)
        }

    case hours < 24:
        if hours == 1 {
            return "منذ ساعة واحدة"
        } else if hours == 2 {
            return "منذ ساعتين"
        } else if hours < 11 {
            return fmt.Sprintf("منذ %d ساعات", hours)
        } else {
            return fmt.Sprintf("منذ %d ساعة", hours)
        }

    case days < 7:
        if days == 1 {
            return "منذ يوم واحد"
        } else if days == 2 {
            return "منذ يومين"
        } else {
            return fmt.Sprintf("منذ %d أيام", days)
        }

    default:
        return t.Format("02-01-2006") // مثل 27-06-2025
    }
}







func (s *Store) UpdateIsReadNotifications(userID int) error {
    query := `UPDATE notifications SET is_read = true WHERE receiver_id = $1 AND is_read = false`
    _, err := s.db.Exec(query, userID)
    return err 
}










func (s *Store) UpdatePatientBasicInfo(p types.UpdatePatientInfo , id int) (*types.UpdatePatientInfo, error){

    tx, err := s.db.Begin()
    if err != nil {
        return &types.UpdatePatientInfo{}, err
    }

    // جلب البريد الحالي من جدول المرضى
    var oldEmail string
    err = tx.QueryRow(`SELECT email FROM patients WHERE id = $1`, id).Scan(&oldEmail)
    if err != nil {
        tx.Rollback()
        return &types.UpdatePatientInfo{}, fmt.Errorf("failed to fetch existing email: %w", err)
    }

    // تحديث بيانات المريض
    _, err = tx.Exec(`
        UPDATE patients
        SET fullName = $1, email = $2, phone = $3, id_number = $4, date = $5
        WHERE id = $6
    `, p.FullName, p.Email, p.Phone, p.IDNumber, p.Date, id)
    if err != nil {
        tx.Rollback()
        return &types.UpdatePatientInfo{}, err
    }

    // تحديث البريد في جدول login_serach باستخدام البريد القديم
    _, err = tx.Exec(`
        UPDATE login_serach
        SET email = $1
        WHERE email = $2
    `, p.Email, oldEmail)
    if err != nil {
        tx.Rollback()
        return &types.UpdatePatientInfo{}, err
    }

    if err := tx.Commit(); err != nil {
        return &types.UpdatePatientInfo{}, err
    }

    return &p, nil
}




func (s *Store) GetCenterIDByName(name string) (int, error) {
	var id int
	err := s.db.QueryRow(`SELECT id FROM centers WHERE centerName = $1`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}







func (s *Store) GetCenterByID(id int) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE id=$1",id)
	if err != nil {
		return nil , err
	}
    
defer rows.Close()

	c := new(types.Center)
	for rows.Next() {
		c , err = scanRowIntoCenter(rows)
		if err != nil {
			return nil , err
		}
	}

	if c.ID == 0 {
		return nil , fmt.Errorf("center not found")
	}

	return c , nil
}


func scanRowIntoCenter(rows *sql.Rows) (*types.Center , error ){
	center := new(types.Center)

	err := rows.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}


func (s *Store) UpdatePatientCenterInfo(id int, centerName string) (types.UpdatePatientCenterInfo, error) {
	centerID, err := s.GetCenterIDByName(centerName)
	if err != nil {
		return types.UpdatePatientCenterInfo{}, fmt.Errorf("center not found: %w", err)
	}

	_, err = s.db.Exec(`
		UPDATE patients
		SET  center_id = $1
		WHERE id = $2
	`,  centerID, id)
	if err != nil {
		return types.UpdatePatientCenterInfo{}, err
	}

    center , err := s.GetCenterByID(centerID)
	if err != nil {
		return types.UpdatePatientCenterInfo{} , err
	}

	result := types.UpdatePatientCenterInfo{
		City:       center.CenterCity,
		CenterName: center.CenterName,
	}

	return result, nil
	
}








func (s *Store) ChangePatientPassword(patientID int, payload types.ChangePassword) error {

	var email, storedHashedPassword string
	err := s.db.QueryRow(`SELECT email, password FROM patients WHERE id = $1`, patientID).Scan(&email, &storedHashedPassword)
	if err != nil {
		return fmt.Errorf("failed to get patient credentials: %w", err)
	}

	
	if !auth.ComparePasswords(storedHashedPassword,[]byte(payload.OldPassword)) {
		return fmt.Errorf("old password is incorrect")
	}

	// تشفير كلمة المرور الجديدة
	hashedNewPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// تحديث كلمة المرور الجديدة في جدول login_serach
	_, err = s.db.Exec(`UPDATE login_serach SET password = $1 WHERE email = $2`, hashedNewPassword, email)
	if err != nil {
		return fmt.Errorf("failed to update password in login_serach: %w", err)
	}

	// تحديث كلمة المرور في جدول patients أيضاً (إذا مخزنة هناك)
	_, err = s.db.Exec(`UPDATE patients SET password = $1 WHERE id = $2`, hashedNewPassword, patientID)
	if err != nil {
		return fmt.Errorf("failed to update password in patients: %w", err)
	}

	return nil
}













// count alll patient


//عدد المرضى الكلي في النظام لدي

func (s *Store) GetTotalPatients() (int, error) {
    var total int
    err := s.db.QueryRow(`SELECT COUNT(*) FROM patients`).Scan(&total)
    if err != nil {
        return 0, err
    }
    return total, nil
}

// عدد المرضى المسجلين ف الشهر الاخير

func (s *Store) GetPatientsLastMonth() (int, error) {
    var total int
    err := s.db.QueryRow(`
        SELECT COUNT(*) 
        FROM patients 
        WHERE createAt >= NOW() - INTERVAL '1 month'
    `).Scan(&total)

    if err != nil {
        return 0, err
    }
    return total, nil
}
