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

func (s *Store) GetPatientByEmail(email string) (*types.Patient , error) {
	rows , err := s.db.Query("SELECT * FROM patients WHERE email=$1", email)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	p := new(types.Patient)
	for rows.Next() {
		p , err = scanRowIntoPatientBy(rows)
		if err != nil {
			return nil , err
		}
	}

	if p.ID == 0 {
		return nil , fmt.Errorf("patient not found")
	}

	return p , nil
}

func scanRowIntoPatientBy(rows *sql.Rows) (*types.Patient , error ){
	patient := new(types.Patient)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Password,
		&patient.CenterID,
		&patient.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return patient , nil
}

func  (s *Store) GetPatientById(id int) (*types.Patient , error) {
	rows , err := s.db.Query("SELECT * FROM patients WHERE id=$1",id)
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
	_ , err := s.db.Exec("INSERT INTO patients (fullName , email , password ,phone , date , center_id) VALUES ($1,$2,$3,$4,$5,$6)" , patient.FullName , patient.Email , patient.Password,patient.Phone, patient.Age , patient.CenterID)
	if err  != nil {
		return err
	}

	return nil
}


func (s *Store) SetPersonlPatientBasicInfo(basicInfo types.BasicPatientInfo) error{
	_ , err := s.db.Exec("INSERT INTO basic_patient_info (patient_id , gender , wight , length_patient ,address_patient , id_number) VALUES ($1,$2,$3,$4,$5,$6)" , basicInfo.PatientID , basicInfo.Gender , basicInfo.Weight , basicInfo.Length , basicInfo.Address , basicInfo.IDNumber )
	if err  != nil {
		return err
	}

	return nil
}