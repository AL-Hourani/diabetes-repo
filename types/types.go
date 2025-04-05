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
	GetPatientProfile(id int)(*PatientProfile , error)
	UpdatePatientProfile(patientPayload ParientUpdatePayload) error
	GetSugarTypeStats(centerID int) ([]*SugarTypeStats, error)
	GetUpdatePatientProfile(id int) (*GetPatientUpdateProfile , error) 

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
	City            string       `json:"city"`
}

type RegisterPatientPayload struct {
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	Password		string		 `json:"password" validate:"required"`
	IDNumber		string       `json:"id_number" validate:"required"`
	City            string       `json:"city" validate:"required"`
	CenterName		string		 `json:"center_name" validate:"required"`
      
}

type ParientUpdatePayload struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	IDNumber		string       `json:"id_number" validate:"required"`
	City            string       `json:"city" validate:"required"`
}
type GetPatientUpdateProfile struct {
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	IDNumber		string       `json:"id_number" validate:"required"`
	City            string       `json:"city" validate:"required"`
}

type PatientProfile struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Age             string		 `json:"age"`
	Phone			string       `json:"phone"`
	IDNumber		string       `json:"id_number"`
	IsCompleted     bool         `json:"isCompleted"`
	Gender                *string   `json:"gender"`
	Weight                 *string   `json:"weight"`
	LengthPatient           *string   `json:"length_patient"`
	AddressPatient          *string     `json:"address_patient"`
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
	CenterID		int			 `json:"center_id"`
	City            string       `json:"city"`
}


type PatientDetails struct {
	ID                    int       `json:"id"`
	FullName              *string    `json:"fullName"`
	Email                 *string    `json:"email"`
	Phone                 *string    `json:"phone"`
	Date                  *string    `json:"date"`
	IDNumber              *string    `json:"id_number"`
	IsCompleted           *bool      `json:"isCompleted"`
	Gender                *string   `json:"gender"`
	Weight                 *string   `json:"weight"`
	LengthPatient           *string   `json:"length_patient"`
	AddressPatient          *string     `json:"address_patient"`
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
}

type PatientUpdatePayload struct {
    ID                    int         `json:"id"`
	IsCompleted             bool      `json:"isCompleted"`
	Gender                  *string   `json:"gender"`
	Weight                  *string   `json:"weight"`
	LengthPatient           *string   `json:"length_patient"`
	AddressPatient          *string  `json:"address_patient"`
	BloodSugar              *string  `json:"bloodSugar"`
	Hemoglobin              *string   `json:"hemoglobin"`
	BloodPressure           *string   `json:"bloodPressure"`
	SugarType               *string  `json:"sugarType"`
	DiseaseDetection        *string   `json:"diseaseDetection"`
	OtherDisease            *string   `json:"otherDisease"`
	TypeOfMedicine          *string   `json:"typeOfMedicine"`
	UrineAcid               *string   `json:"urineAcid"`
	Cholesterol             *string   `json:"cholesterol"`
	Grease                  *string   `json:"grease"`
	HistoryOfFamilyDisease  *string `json:"historyOfFamilyDisease"`
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
	ID				int			 `json:"id"`
	Name		    string		 `json:"name"`
	Email			string		 `json:"email"`
	Role 			string       `json:"role"`
	IsCompletes     bool         `json:"isCompleted"`
	Token            string      `json:"token"`
}
type ReturnLoggingCenterData struct {
	ID				int			 `json:"id"`
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
	GetCentersByCity(cityName string)([]Center , error)
	DeletePatient(id int) error
	UpdateIsCompletedPatientField(confirmAcc ConfirmAccount) error
	PatchUpdatePatient(patient *PatientUpdatePayload) error
	GetPatientsForCenter(CenterID int) ([]CardData , error)
	GetCities()([]string , error)
	GetPatientCountByCenterName(centerName string) (int, error)
	GetCenterProfile(id int) (*CenterProfile, error)
	DeleteCenter(id int) error 
	DeleteCenterAndReassignPatients(centerID int, newCenterID int) error
	CenterUpdateCenterProfile(centerUpdate CenterUpdateProfilePayload) error 
     GetCenterUpdateCenterProfile(id int)(*GetCenterUpdateProfile , error)
}

type Center struct {
	ID				int          `json:"centerId"`
    CenterName		string		 `json:"centerName"`
	CenterPassword  string       `json:"centerPassword"`
	CenterEmail	    string       `json:"centerEmail"`
	CenterCity      string       `json:"centerCity"`
	CreateAt        time.Time    `json:"createAt"`
}



type RegisterCenterPayload struct {
	CenterName		string		 `json:"centerName"  validate:"required"`
	CenterPassword  string       `json:"centerPassword"  validate:"required"`
	CenterEmail	    string       `json:"centerEmail"    validate:"required,email"`
	CenterCity      string       `json:"centerCity"    validate:"required"`
	CenterKey		string		 `json:"centerKey" validate:"required"`
}


type CenterProfile struct {
	ID				int          `json:"centerId"`
	CenterName		string		 `json:"centerName"`
	CenterEmail	    string       `json:"centerEmail"`
	CenterCity      string       `json:"centerCity"`
	PatientNumber   int          `json:"patientNumber"`

}

type DeleteCenter struct {
	CenterName		string		 `json:"centerName"`
	CenterNameReassignPatients  string `json:"centerNameReassignPatients"`
}

type CenterUpdateProfilePayload struct {
	ID				int          `json:"centerId"`
	CenterName		string		 `json:"centerName"  validate:"required"`
	CenterEmail	    string       `json:"centerEmail"    validate:"required,email"`
	CenterCity      string       `json:"centerCity"    validate:"required"`
}
type GetCenterUpdateProfile struct {
	CenterName		string		 `json:"centerName"  validate:"required"`
	CenterEmail	    string       `json:"centerEmail"    validate:"required,email"`
	CenterCity      string       `json:"centerCity"    validate:"required"`
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
	OTPCode string `json:"code"`
}


type ConfirmAccount struct {
	ID				int			 `json:"id" validate:"required"`
	IsCompleted     bool         `json:"isCompleted" validate:"required"`
}





// config email 
type EmailRequest struct {
	Sender      SenderInfo    `json:"sender"`
	To          []Recipient   `json:"to"`
	Subject     string        `json:"subject"`
	HTMLContent string        `json:"htmlContent"`
}

// SenderInfo معلومات المرسل
type SenderInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Recipient معلومات المستلم
type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}





// anlisze 

type SugarTypeStats struct {
	SugarType string `json:"sugarType"`
	Total     int    `json:"total"`
}

