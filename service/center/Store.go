package center

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AL-Hourani/care-center/types"
)

type Store struct {
	db *sql.DB
}


func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCenterByName(centerName string) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE centerName=$1",centerName)
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
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}

//get center by emial 
func (s *Store) GetCenterByEmail(centerEmail string) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE centerEmail=$1",centerEmail)
	if err != nil {
		return nil , err
	}
    
defer rows.Close()

	c := new(types.Center)
	for rows.Next() {
		c , err = scanRowIntoCenterByEmail(rows)
		if err != nil {
			return nil , err
		}
	}

	if c.ID == 0 {
		return nil , fmt.Errorf("center not found")
	}

	return c , nil
}


func scanRowIntoCenterByEmail(rows *sql.Rows) (*types.Center , error ){
	center := new(types.Center)

	err := rows.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterEmail,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}







func (s *Store)	GreateCenter(center types.Center) error {
	_ , err := s.db.Exec("INSERT INTO centers (centerName ,centerPassword , centerEmail) VALUES ($1, $2, $3)" , center.CenterName , center.CenterPassword , center.CenterEmail)
	if err  != nil {
		return err
	}

	return nil
}




//this is not completed

func (s *Store)	GetPatients(centerId int)([]types.Patient , error) {
	rows , err := s.db.Query("SELECT * FROM patients WHERE center_id=?" ,centerId )
	if err != nil {
		return nil , err
	}

	patients := make([]types.Patient , 0)
	for rows.Next() {
		p , err := scanRowIntoPatients(rows)
		if err != nil {
			return nil , err
		}
		patients = append(patients, *p)
	}

	return patients , nil

}


func scanRowIntoPatients(rows *sql.Rows) (*types.Patient , error ){
	patient := new(types.Patient)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.IDNumber,
		&patient.IsCompleted,
		&patient.CenterID,
		&patient.Password,
		&patient.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return patient , nil
}



func (s *Store) GetCenters()([]types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers")
	if err != nil {
		return nil , err
	}

	centers := make([]types.Center , 0)
	for rows.Next() {
		p , err := scanRowIntoCenters(rows)
		if err != nil {
			return nil , err
		}
		centers = append(centers, *p)
	}

	return centers , nil
}


func scanRowIntoCenters(rows *sql.Rows) (*types.Center , error ){
	center := new(types.Center)

	err := rows.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterEmail,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center, nil
}
















// deleted patient..............
func (s *Store) DeletePatient(id int) error { 
	_ , err := s.db.Query(`DELETE FROM patients WHERE id = $1` , id)
	if err != nil {
		return err
	}

	return nil

}

//update iscompleate field in patients tabel
func (s *Store) UpdateIsCompletedPatientField(confirmAcc types.ConfirmAccount) error { 
	_, err := s.db.Exec("UPDATE patients SET isCompleted = $1 WHERE id = $2",confirmAcc.IsCompleted,confirmAcc.ID )
		if err != nil {
			log.Fatal(err)
		}
		
	return nil
}


func (s *Store) PatchUpdatePatient(patient *types.PatientUpdatePayload) error {
	query := `UPDATE patients
	SET
		fullName = COALESCE($1, fullName),
		email = COALESCE($2, email),
		phone = COALESCE($3, phone),
		date = COALESCE($4, date),
		id_number = COALESCE($5, id_number),
		gender = COALESCE($6, gender),
		wight = COALESCE($7, wight),
		length_patient = COALESCE($8, length_patient),
		address_patient = COALESCE($9, address_patient),
		bloodSugar = COALESCE($10, bloodSugar),
		hemoglobin = COALESCE($11, hemoglobin),
		bloodPressure = COALESCE($12, bloodPressure),
		sugarType = COALESCE($13, sugarType),
		diseaseDetection = COALESCE($14, diseaseDetection),
		otherDisease = COALESCE($15, otherDisease),
		typeOfMedicine = COALESCE($16, typeOfMedicine),
		urineAcid = COALESCE($17, urineAcid),
		cholesterol = COALESCE($18, cholesterol),
		grease = COALESCE($19, grease),
		historyOfFamilyDisease = COALESCE($20, historyOfFamilyDisease)
	WHERE id = $21`

	_, err := s.db.Exec(query,
        patient.FullName, patient.Email, patient.Phone,
        patient.Age, patient.IDNumber, patient.Gender, patient.Weight, patient.LengthPatient,
        patient.AddressPatient, patient.BloodSugar, patient.Hemoglobin, patient.BloodPressure,
        patient.SugarType, patient.DiseaseDetection, patient.OtherDisease, patient.TypeOfMedicine,
        patient.UrineAcid, patient.Cholesterol, patient.Grease, patient.HistoryOfFamilyDisease,
        patient.ID)
    
    if err != nil {
        return fmt.Errorf("error updating patient: %v", err)
    }

    return nil
}