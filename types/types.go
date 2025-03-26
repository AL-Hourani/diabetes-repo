package types

import "time"



//postgress strorge 


// patients..............................

type PatientStore interface {
	GetPatientByEmail(email string) (*Patient , error)
	GetPatientById(id int) (*Patient , error)
	GreatePatient(Patient) error
	SetPersonlPatientBasicInfo(BasicPatientInfo) error
	SetPatientHealthInfo(HealthPatientData) error
	GetPatientsForCenter(CenterID int) ([]CardData , error) 
	GetAllPatientInfo(id int) (*ReaturnAllPatientInfo , error)

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



type BasicPatientInfo struct {
	ID				int			   `json:"id"`
	PatientID		int            `json:"patientID" validate:"required"`
	Weight			string         `json:"weight" validate:"required"`
	Length          string		   `json:"lenght" validate:"required"`
	Address			string         `json:"address" validate:"required"`
	Gender          string         `json:"gender" validate:"required"`
    CreateAt        time.Time      `json:"createAt"`
} 

type BasicPatientInfoPalyoad struct {
	PatientID		int            `json:"patientID" validate:"required"`
	Weight			string         `json:"weight" validate:"required"`
	Length          string		   `json:"lenght" validate:"required"`
	Address			string         `json:"address" validate:"required"`
	Gender          string         `json:"gender" validate:"required"`
	IDNumber        string         `json:"idNumber" validate:"required"`
} 

type LoginPayload struct {
	Email			string		 `json:"email"    validate:"required,email"`
	Password		string		 `json:"password" validate:"required"`
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
	Patient			[]CardData   `json:"patient"`
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
	SugarType		string       `json:"sugarType"`
}

type HealthPatientData struct {
	ID              int            `json:"id" validate:"required"`
	PatientID		int            `json:"patientID" validate:"required"`
	BloodSugar      string		   `json:"booldSugar" validate:"required"`	
	Hemoglobin      string         `json:"hemoglobin" validate:"required"`
	BloodPressure   string         `json:"bloodPressure" validate:"required"`
	SugarType		string         `json:"sugarType" validate:"required"`	
	DiseaseDetection string        `json:"diseaseDetection" validate:"required"`
	OtherDisease     string        `json:"OtherDisease" validate:"required"`
	TypeOfMedicine   string        `json:"typeOfMedicine" validate:"required"`
	UrineAcid        string        `json:"urineAcid" validate:"required"`
	Cholesterol      string        `json:"cholesterol" validate:"required"`
	Grease			 string        `json:"grease" validate:"required"`
	HistoryOfFamilyDisease string  `json:"historyOfFamilyDisease" validate:"required"`
	CreateAt        time.Time      `json:"createAt"`
}


type RegisterHealthPatientData struct {
	PatientID		int            `json:"patientID" validate:"required"`
	BloodSugar      string		   `json:"booldSugar" validate:"required"`	
	Hemoglobin      string         `json:"hemoglobin" validate:"required"`
	BloodPressure   string         `json:"bloodPressure" validate:"required"`
	SugarType		string         `json:"sugarType" validate:"required"`	
	DiseaseDetection string        `json:"diseaseDetection" validate:"required"`
	OtherDisease     string        `json:"OtherDisease" validate:"required"`
	TypeOfMedicine   string        `json:"typeOfMedicine" validate:"required"`
	UrineAcid        string        `json:"urineAcid" validate:"required"`
	Cholesterol      string        `json:"cholesterol" validate:"required"`
	Grease			 string        `json:"grease" validate:"required"`
	HistoryOfFamilyDisease string  `json:"historyOfFamilyDisease" validate:"required"`
}


type ReaturnAllPatientInfo struct {
	FullName		string		   `json:"fullname"`
	Email			string		   `json:"email"`
	Age             string		   `json:"age"`
	Phone			string         `json:"phone"`
	Weight			string         `json:"weight"`
	Length          string		   `json:"lenght" `
	Address			string         `json:"address"`
	Gender          string         `json:"gender" `
	IsCompleted     bool           `json:"isCompleted"`
	IDNumber        string         `json:"idNumber"`
	BloodSugar      string		   `json:"booldSugar" `	
	Hemoglobin      string         `json:"hemoglobin"`
	BloodPressure   string         `json:"bloodPressure" `
	SugarType		string         `json:"sugarType"`	
	DiseaseDetection string        `json:"diseaseDetection"`
	OtherDisease     string        `json:"OtherDisease" `
	TypeOfMedicine   string        `json:"typeOfMedicine"`
	UrineAcid        string        `json:"urineAcid"`
	Cholesterol      string        `json:"cholesterol"`
	Grease			 string        `json:"grease"`
	HistoryOfFamilyDisease string  `json:"historyOfFamilyDisease"`
}


//end.......................................
//center....................................

type CenterStore interface {
	GetCenterByName(centerName string) (*Center , error) 
	GetCenterByEmail(centerEmail string) (*Center , error)
	GreateCenter(Center) error
	GetPatients(centerID int)([]Patient , error)
	GetCenters()([]Center , error)
	DeletePatient(id int) error
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













