package patient

import (
	"database/sql"
	"fmt"


	"github.com/AL-Hourani/care-center/types"
)

type Store struct {
	db *sql.DB
}


func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
func  (s *Store) GetPatientByEmail(email string) (*types.Patient , error) {
	rows , err := s.db.Query("SELECT id,fullName,email,password,phone,date,id_number,isCompleted,createAt FROM patients WHERE email=$1",email)
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
	rows , err := s.db.Query("SELECT id,fullName,email,password,phone,date,id_number,isCompleted,createAt FROM patients WHERE id=$1",id)
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



// 2




// 3
func (s *Store)	GreatePatient(patient types.Patient) error {
	_ , err := s.db.Exec("INSERT INTO patients (fullName , email , password ,phone , date , id_number , isCompleted , center_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)" , patient.FullName , patient.Email , patient.Password,patient.Phone, patient.Age ,patient.IDNumber , patient.IsCompleted , patient.CenterID)
	if err  != nil {
		return err
	}

	return nil
}

// get all data form patient...

func (s *Store) GetPatientDetailsByID(patientID int) (*types.PatientDetails, error) {
	rows , err := s.db.Query(`SELECT id,fullName, email, phone, date, id_number,
		isCompleted, gender, wight, length_patient, address_patient,
		bloodSugar, hemoglobin, bloodPressure, sugarType, diseaseDetection,
		otherDisease, typeOfMedicine, urineAcid, cholesterol, grease,
		historyOfFamilyDisease, center_id, createAt
	
	FROM patients
	WHERE id=$1`,patientID)

	if err != nil {
		return nil , err
	}
	defer rows.Close()

	patient := new(types.PatientDetails)
	for rows.Next() {
		patient , err = scanRowIntoPatientDeatials(rows)
		if err != nil {
			return nil , err
		}
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


func scanRowIntoPatientDeatials(rows *sql.Rows) (*types.PatientDetails , error ){
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
		&patient.CenterID,
		&patient.CreateAt,
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
	query := `
	(
		SELECT 'patient' AS role, id, fullName AS name, email, password, NULL AS centerName
		FROM patients
		WHERE email = $1
	)
	UNION 
	(
		SELECT 'center' AS role, id, centerName AS name, centerEmail AS email, centerPassword AS password, centerName
		FROM centers
		WHERE centerEmail = $2
	)
	ORDER BY role
	LIMIT 1`

    row := s.db.QueryRow(query, email, email)
	
	user := new(types.UserLoginData)

	err := row.Scan(
		&user.Role,
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CenterName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user , nil
}



