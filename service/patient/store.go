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

func (s *Store) GetPatientsForCenter(CenterID int) ([]types.CardData , error) {
	return nil , nil
}



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


func scanRowIntoPatientDeatials(rows *sql.Rows) (*types.PatientDetails , error ){
	patient := new(types.PatientDetails)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Phone,
		&patient.Date,
		&patient.IDNumber,
		&patient.IsCompleted,
		&patient.Gender,
		&patient.Weight,
		&patient.LengthPatient,
		&patient.AddressPatient,
		&patient.BloodSugar,
		&patient.Hemoglobin,
		&patient.BloodPressure,
		&patient.SugarType,
		&patient.DiseaseDetection,
		&patient.OtherDisease,
		&patient.TypeOfMedicine,
		&patient.UrineAcid,
		&patient.Cholesterol,
		&patient.Grease,
		&patient.HistoryOfFamilyDisease,
		&patient.CenterID,
		&patient.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return patient , nil
}

