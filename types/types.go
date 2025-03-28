package types

import (

	"time"
)

//postgress strorge

// patients..............................

type PatientStore interface {
	GetPatientByEmail(email string) (*Patient , error)
	GetPatientById(id int) (*Patient , error)
	GreatePatient(Patient) error
	GetPatientDetailsByID(patientID int) (*PatientDetails, error)
	GetUserByEmail(email string) (*UserLoginData, error)

}

type Patient struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Password		string       `json:"password"`
	Age             string		 `json:"age"`
	Phone			string       `json:"phone"`
	IDNumber		string       `json:"id_number"`
	IsCompleted     bool         `json:"isCompleted"`
	CenterID		int			 `json:"center_id"`
	CreateAt        time.Time    `json:"createAt"`
}

type RegisterPatientPayload struct {
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	Password		string		 `json:"password" validate:"required"`
	IDNumber		string       `json:"id_number" validate:"required"`
	CenterName		string		 `json:"center_name" validate:"required"`
      
}


type PatientDetails struct {
	ID                    int       `json:"id"`
	FullName              string    `json:"fullName"`
	Email                 string    `json:"email"`
	Phone                 string    `json:"phone"`
	Date                  string    `json:"date"`
	IDNumber              string    `json:"id_number"`
	IsCompleted           bool      `json:"isCompleted"`
	Gender                *string   `json:"gender"`
	Weight                *string   `json:"weight"`
	LengthPatient          *string   `json:"length_patient"`
	AddressPatient       *string     `json:"address_patient"`
	BloodSugar             *string  `json:"bloodSugar"`
	Hemoglobin             *string   `json:"hemoglobin"`
	BloodPressure          *string   `json:"bloodPressure"`
	SugarType              *string  `json:"sugarType"`
	DiseaseDetection       *string   `json:"diseaseDetection"`
	OtherDisease           *string   `json:"otherDisease"`
	TypeOfMedicine         *string   `json:"typeOfMedicine"`
	UrineAcid             *string   `json:"urineAcid"`
	Cholesterol            *string   `json:"cholesterol"`
	Grease                *string   `json:"grease"`
	HistoryOfFamilyDisease  *string `json:"historyOfFamilyDisease"`
	CreateAt              time.Time `json:"createAt"`
}

type PatientUpdatePayload struct
{
	ID				        int			 `json:"id"`
	FullName		        string		 `json:"fullname"`
	Email			        string		 `json:"email"`
	Age                     string		 `json:"age"`
	Phone			        string       `json:"phone"`
	IDNumber		        string       `json:"id_number"`
	IsCompleted             bool      `json:"isCompleted"`
	Gender                  string   `json:"gender"`
	Weight                  string   `json:"weight"`
	LengthPatient           string   `json:"length_patient"`
	AddressPatient          string  `json:"address_patient"`
	BloodSugar              string  `json:"bloodSugar"`
	Hemoglobin              string   `json:"hemoglobin"`
	BloodPressure           string   `json:"bloodPressure"`
	SugarType               string  `json:"sugarType"`
	DiseaseDetection        string   `json:"diseaseDetection"`
	OtherDisease            string   `json:"otherDisease"`
	TypeOfMedicine          string   `json:"typeOfMedicine"`
	UrineAcid               string   `json:"urineAcid"`
	Cholesterol             string   `json:"cholesterol"`
	Grease                  string   `json:"grease"`
	HistoryOfFamilyDisease  string `json:"historyOfFamilyDisease"`
	CreateAt                time.Time    `json:"createAt"`
}





type LoginPayload struct {
	Email			string		 `json:"email"    validate:"required,email"`
	Password		string		 `json:"password" validate:"required"`
}

//البيانات المشتركة 
type UserLoginData struct {
	Role       string  `json:"role"`
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Password   string  `json:"-"`
	CenterName *string `json:"centerName,omitempty"`
}



type ReturnLoggingData struct {
	Name		    string		 `json:"name"`
	Email			string		 `json:"email"`
	Role 			string       `json:"role"`
	IsCompletes     bool         `json:"isCompleted"`
	Token            string      `json:"token"`
}
type ReturnLoggingCenterData struct {
	Name		    string		 `json:"name"`
	Email			string		 `json:"email"`
	Role 			string       `json:"role"`
	Token           string       `json:"token"`
}

type CardData struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Age             string		 `json:"age"`
	Phone			string       `json:"phone"`
	IDNumber		string       `json:"id_number"`
	IsCompleted     bool         `json:"isCompleted"`
	SugarType		*string       `json:"sugarType"`
}






//end.......................................
//center....................................

type CenterStore interface {
	GetCenterByName(centerName string) (*Center , error) 
	GetCenterByEmail(centerEmail string) (*Center , error)
	GreateCenter(Center) error
	GetCenters()([]Center , error)
	DeletePatient(id int) error
	UpdateIsCompletedPatientField(confirmAcc ConfirmAccount) error
	PatchUpdatePatient(patient *PatientUpdatePayload) error
	GetPatientsForCenter(CenterID int) ([]CardData , error)
}

type Center struct {
	ID				int          `json:"centerId"`
    CenterName		string		 `json:"centerName"`
	CenterPassword  string       `json:"centerPassword"`
	CenterEmail	    string       `json:"centerEmail"`
	CreateAt        time.Time    `json:"createAt"`
}



type RegisterCenterPayload struct {
	CenterName		string		 `json:"centerName"  validate:"required"`
	CenterPassword  string       `json:"centerPassword"  validate:"required"`
	CenterEmail	    string       `json:"centerEmail"    validate:"required,email"`
	CenterKey		string		 `json:"centerKey" validate:"required"`

}

// type LoginCenterPayload struct {
// 	CenterEmail		string		 `json:"centerEmail"    validate:"required,email"`
// 	CenterPassword  string       `json:"centerPassword"  validate:"required"`
// }

//end....................................................



type OTPRequest struct {
	Email string `json:"email"`
}


type VerifyRequest struct {
	Email   string `json:"email"`
	OTPCode string `json:"otp"`
}


type ConfirmAccount struct {
	ID				int			 `json:"id" validate:"required"`
	IsCompleted     bool         `json:"isCompleted" validate:"required"`
}







