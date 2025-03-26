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

func  (s *Store) GetPatientsForCenter(CenterID int) ([]types.CardData , error) {
	rows , err := s.db.Query(`SELECT p.id,fullName,email,phone,isCompleted,date,id_number,sugarType 
	 FROM patients p LEFT JOIN patient_health_info h ON p.id = h.patient_id  WHERE p.center_id=$1`,CenterID)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	patientsCard := make([]types.CardData , 0)
	for rows.Next() {
		p , err := scanRowIntoPatients(rows)
		if err != nil {
			return nil , err
		}

		patientsCard = append(patientsCard , *p)
	}

	return patientsCard , nil
}


func scanRowIntoPatients(rows *sql.Rows) (*types.CardData , error ){
	patient := new(types.CardData)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Phone,
		&patient.IsCompleted,
		&patient.Age,
		&patient.IDNumber,
		&patient.SugarType,
        
	)
	
	if err  != nil {
		return nil , err
	}

	return patient , nil
}




// 3
func (s *Store)	GreatePatient(patient types.Patient) error {
	_ , err := s.db.Exec("INSERT INTO patients (fullName , email , password ,phone , date , id_number , isCompleted , center_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)" , patient.FullName , patient.Email , patient.Password,patient.Phone, patient.Age ,patient.IDNumber , patient.IsCompleted , patient.CenterID)
	if err  != nil {
		return err
	}

	return nil
}


func (s *Store) SetPersonlPatientBasicInfo(basicInfo types.BasicPatientInfo) error{
	_ , err := s.db.Exec("INSERT INTO basic_patient_info (patient_id , gender , wight , length_patient ,address_patient) VALUES ($1,$2,$3,$4,$5)" , basicInfo.PatientID , basicInfo.Gender , basicInfo.Weight , basicInfo.Length , basicInfo.Address  )
	if err  != nil {
		return err
	}

	return nil
}


func (s *Store) SetPatientHealthInfo(healthInfo types.HealthPatientData) error {
	_ , err := s.db.Exec("INSERT INTO patient_health_info (patient_id , bloodSugar , hemoglobin , bloodPressure ,sugarType , diseaseDetection , otherDisease , typeOfMedicine , urineAcid  , cholesterol , grease , historyOfFamilyDisease) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)" , healthInfo.PatientID , healthInfo.BloodSugar , healthInfo.Hemoglobin , healthInfo.BloodPressure ,healthInfo.SugarType ,  healthInfo.DiseaseDetection , healthInfo.OtherDisease , healthInfo.TypeOfMedicine , healthInfo.UrineAcid , healthInfo.Cholesterol , healthInfo.Grease , healthInfo.HistoryOfFamilyDisease)
	if err  != nil {
		return err
	}

	return nil
}


func (s *Store) GetAllPatientInfo(id int) (*types.ReaturnAllPatientInfo , error) {
	rows , err := s.db.Query(`SELECT fullName,email,phone,isCompleted,date,wight,length_patient,address_patient,gender,id_number,bloodSugar,hemoglobin,bloodPressure,
	sugarType,diseaseDetection,otherDisease,typeOfMedicine,urineAcid,cholesterol,grease,historyOfFamilyDisease
	FROM patients p INNER JOIN  basic_patient_info b ON p.id = b.patient_id  INNER JOIN 
	patient_health_info h ON p.id = h.patient_id  WHERE p.id=$1`,id)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	allPatentInfo := new(types.ReaturnAllPatientInfo)
	for rows.Next() {
		allPatentInfo , err = scanRowIntoAllPatient(rows)
		if err != nil {
			return nil , err
		}
	}

	return allPatentInfo , nil
}

func scanRowIntoAllPatient(rows *sql.Rows) (*types.ReaturnAllPatientInfo , error ){
	allpatientInfo := new(types.ReaturnAllPatientInfo)

	err := rows.Scan(
		&allpatientInfo.FullName,
		&allpatientInfo.Email,
		&allpatientInfo.Phone,
		&allpatientInfo.IsCompleted,
		&allpatientInfo.Age,
		&allpatientInfo.Weight,
		&allpatientInfo.Length,
		&allpatientInfo.Address,
		&allpatientInfo.Gender,
		&allpatientInfo.IDNumber,
		&allpatientInfo.BloodSugar,
		&allpatientInfo.Hemoglobin,
		&allpatientInfo.BloodPressure,
		&allpatientInfo.SugarType,
		&allpatientInfo.DiseaseDetection,
		&allpatientInfo.OtherDisease,
		&allpatientInfo.TypeOfMedicine,
		&allpatientInfo.UrineAcid,
		&allpatientInfo.Cholesterol,
		&allpatientInfo.Grease,
		&allpatientInfo.HistoryOfFamilyDisease,

	)
	
	if err  != nil {
		return nil , err
	}

	return allpatientInfo , nil
}