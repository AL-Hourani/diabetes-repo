package types

import "time"

// patients..............................

type PatientStore interface {
	GetPatientByEmail(email string) (*Patient , error)
	GetPatientById(id int) (*Patient , error)
	GreatePatient(Patient) error

}

type Patient struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Password		string       `json:"password"`
	CenterID		int			 `json:"center_id"`
	CreateAt        time.Time    `json:"createAt"`
}

type RegisterPatientPayload struct {
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Password		string		 `json:"password" validate:"required"`
	CenterName		string		 `json:"center_name" validate:"required"`
}

type LoginPatientPayload struct {
	Email			string		 `json:"email"    validate:"required,email"`
	Password		string		 `json:"password" validate:"required"`
}

//end.......................................
//center....................................

type CenterStore interface {
	GetCenterByName(centerName string) (*Center , error) 
	GreateCenter(Center) error
	GetPatients(centerID int)([]Patient , error)
	GetCenters()([]Center , error)
}

type Center struct {
	ID				int          `json:"centerId"`
    CenterName		string		 `json:"centerName"`
	CenterPassword  string       `json:"centerPassword"`
	CenterAddress	string       `json:"centerAddress"`
	CreateAt        time.Time    `json:"createAt"`
}



type RegisterCenterPayload struct {
	CenterName		string		 `json:"centerName"  validate:"required"`
	CenterPassword  string       `json:"centerPassword"  validate:"required"`
	CenterAddress	string       `json:"centerAddress"  validate:"required"`
	CenterKey		string		 `json:"centerKey" validate:"required"`

}

type LoginCenterPayload struct {
	CenterName		string		 `json:"centerName"  validate:"required"`
	CenterPassword  string       `json:"centerPassword"  validate:"required"`
}

//end....................................................

type HealthOverviewStore interface {
	SetPatientHealthOverview(PatientHealthOverview) error
}


type PatientHealthOverview struct {
	ID 							  int          `json:"id" validate:"required"`
	PatientID		              int          `json:"patient_id"  validate:"required"` 
	Age							  int			 `json:"age"  validate:"required"`
	Gender						  string       `json:"gender"  validate:"required"`
	DiabetesType				  string 		 `json:"diabetes_type"  validate:"required"`
	DiagnosisData				  string       `json:"diagnosis_data"  validate:"required"`
	CurrentBloodSugerLevel		  string 		`json:"current_blood_suger_level"`
	BloodSugerTrends			  string	    `json:"blood_suger_trends"`
}

type RegisterPatientHealthOverviewPayload struct {
	PatientName		              string       `json:"patientName"  validate:"required"` 
	Age							  int		   `json:"age"  validate:"required"`
	Gender						  string       `json:"gender"  validate:"required"`
	DiabetesType				  string 		`json:"diabetes_type"  validate:"required"`
	DiagnosisData				  string       `json:"diagnosis_data"  validate:"required"`
	CurrentBloodSugerLevel		  string 		`json:"current_blood_suger_level"`
	BloodSugerTrends			  string	    `json:"blood_suger_trends"`
}