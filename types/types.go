package types

import "time"



//postgress strorge 


// patients..............................

type PatientStore interface {
	GetPatientByEmail(email string) (*Patient , error)
	GetPatientById(id int) (*Patient , error)
	GreatePatient(Patient) error
	SetPersonlPatientBasicInfo(BasicPatientInfo) error
	GetPatientsForCenter(CenterID int) ([]int , error) 

}

type Patient struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Password		string       `json:"password"`
	Age             string		 `json:"age"`
	Phone			string        `json:"phone"`
	CenterID		int			 `json:"center_id"`
	CreateAt        time.Time    `json:"createAt"`
}

type RegisterPatientPayload struct {
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	Password		string		 `json:"password" validate:"required"`
	CenterName		string		 `json:"center_name" validate:"required"`
      
}



type BasicPatientInfo struct {
	ID				int			   `json:"id"`
	PatientID		int            `json:"patientID" validate:"required"`
	Weight			string         `json:"weight" validate:"required"`
	Length          string		   `json:"lenght" validate:"required"`
	Address			string         `json:"address" validate:"required"`
	Gender          string         `json:"gender" validate:"required"`
	IDNumber        string         `json:"idNumber" validate:"required"`
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
	IsCompletes     bool         `json:"isCompleted"`
	Patient			[]int        `json:"patient"`
	Token           string      `json:"token"`
}

//end.......................................
//center....................................

type CenterStore interface {
	GetCenterByName(centerName string) (*Center , error) 
	GetCenterByEmail(centerEmail string) (*Center , error)
	GreateCenter(Center) error
	GetPatients(centerID int)([]Patient , error)
	GetCenters()([]Center , error)
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

















