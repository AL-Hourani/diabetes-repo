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
		&center.CenterCity,
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
		&center.CenterCity,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}







func (s *Store)	GreateCenter(center types.Center) error {
	_ , err := s.db.Exec("INSERT INTO centers (centerName ,centerPassword , centerEmail , centerCity) VALUES ($1, $2, $3 ,$4)" , center.CenterName , center.CenterPassword , center.CenterEmail , center.CenterCity)
	if err  != nil {
		return err
	}

	return nil
}




//this is not completed

func (s *Store) GetPatientsForCenter(CenterID int) ([]types.CardData , error) {
	rows , err := s.db.Query("SELECT id,fullName,email,date,phone,id_number,isCompleted,sugarType FROM patients WHERE center_id=$1",CenterID)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	cardData := make([]types.CardData,0)
	for rows.Next() {
		cardd , err := scanRowIntoPatientsCard(rows)
		if err != nil {
			return nil , err
		}
		cardData = append(cardData, *cardd)
	}

	return cardData , nil

}
func convertNullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}


func scanRowIntoPatientsCard(rows *sql.Rows) (*types.CardData , error ){
	patient := new(types.CardData)

    var sugarType     sql.NullString
	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Age,
		&patient.Phone,
		&patient.IDNumber,
		&patient.IsCompleted,
		&sugarType,
	)

	
	if err  != nil {
		return nil , err
	}
	patient.SugarType = convertNullStringToPointer(sugarType)

	return patient , nil
}





// get cities 


func (s *Store) GetCities()([]string , error) {
	rows , err := s.db.Query("SELECT DISTINCT centerCity FROM centers")
	if err != nil {
		return nil , err
	}
	defer rows.Close()
	cities := make([]string ,0)
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			log.Println("Error in reading data", err)
			continue
		}
		cities = append(cities, city)
	}

	return cities , nil
	
}

//-----------------------------------------------


func (s *Store) GetCentersByCity(cityName string)([]types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE centerCity=$1" , cityName)
	if err != nil {
		return nil , err
	}

	defer rows.Close()
	
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
		&center.CenterCity,
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
		isCompleted = COALESCE($1, isCompleted),
		gender = COALESCE($2, gender),
		wight = COALESCE($3, wight),
		length_patient = COALESCE($4, length_patient),
		address_patient = COALESCE($5, address_patient),
		bloodSugar = COALESCE($6, bloodSugar),
		hemoglobin = COALESCE($7, hemoglobin),
		bloodPressure = COALESCE($8, bloodPressure),
		sugarType = COALESCE($9, sugarType),
		diseaseDetection = COALESCE($10, diseaseDetection),
		otherDisease = COALESCE($11, otherDisease),
		typeOfMedicine = COALESCE($12, typeOfMedicine),
		urineAcid = COALESCE($13, urineAcid),
		cholesterol = COALESCE($14, cholesterol),
		grease = COALESCE($15, grease),
		historyOfFamilyDisease = COALESCE($16, historyOfFamilyDisease)
	    WHERE id = $17`

	_, err := s.db.Exec(query,
		patient.IsCompleted,
		nullifyString(patient.Gender),
		nullifyString(patient.Weight),
		nullifyString(patient.LengthPatient),
		nullifyString(patient.AddressPatient),
		nullifyString(patient.BloodSugar),
		nullifyString(patient.Hemoglobin),
		nullifyString(patient.BloodPressure),
		nullifyString(patient.SugarType),
		nullifyString(patient.DiseaseDetection),
		nullifyString(patient.OtherDisease),
		nullifyString(patient.TypeOfMedicine),
		nullifyString(patient.UrineAcid),
		nullifyString(patient.Cholesterol),
		nullifyString(patient.Grease),
		nullifyString(patient.HistoryOfFamilyDisease),
		patient.ID,
	)
    
    if err != nil {
        return fmt.Errorf("error updating patient: %v", err)
    }

    return nil
}

func nullifyString(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

// func nullifyBool(b *bool) interface{} {
// 	if b == nil {
// 		return nil
// 	}
// 	return *b
// }



// get number of patient in any center 
func (s *Store) GetPatientCountByCenterName(centerName string) (int, error) {
    var patientCount int
    err := s.db.QueryRow(`
        SELECT COUNT(*)
        FROM patients
        WHERE center_id = (
            SELECT id FROM centers WHERE centerName = $1
        )
    `, centerName).Scan(&patientCount)

    if err != nil {
        return 0, err
    }

    return patientCount, nil
}

func (s *Store) GetCenterProfile(id int) (*types.CenterProfile, error) { 
	row := s.db.QueryRow("SELECT * FROM centers WHERE id=$1",id)
	
	center := new(types.Center)
	err := row.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterEmail,
		&center.CreateAt,
		&center.CenterCity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	patient_number , err := s.GetPatientCountByCenterName(center.CenterName)
	if err != nil {
		return nil, fmt.Errorf("error in get number of patients")
	}

	cenetrProfile := &types.CenterProfile {
		CenterName: center.CenterName,
        CenterEmail: center.CenterEmail,
		CenterCity: center.CenterCity,
		PatientNumber: patient_number,
	}

	return cenetrProfile , nil
}

