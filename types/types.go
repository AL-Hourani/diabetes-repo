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
	GetSugarTypeStats(centerID int) ([]*Statistics, error)
	GetUpdatePatientProfile(id int) (*GetPatientUpdateProfile , error) 
	GetGenderCounts(centerID int) (int, int, error)
	GetTotalPatientsInSystem() (int, error)
	GetSugarTypeAgeRangeStats(centerID int)([]*SugarAgeRangeStat, error)
	GetSugarTypeAgeRangeStatsAllSystem() ([]*SugarAgeRangeStat, error)
    GetBMIStats(centerID int) ([]*BMIStat, error)
	GetCityStats() ([]*CityStat, error)
   
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

type Review struct {
	Id               int	        `json:"id"`
    CreateAt         string      `json:"create_At"`
}

type CardData struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Age             string		 `json:"age"`
	Phone			string       `json:"phone"`
	IDNumber		string       `json:"id_number"`
	IsCompleted     bool         `json:"isCompleted"`
	SugarType		*string      `json:"sugarType"`
	CreateAt         string      `json:"create_At"`
	Reviews         []Review     `json:"review"`          
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
	InsertReview(reviewdata Reviwe)  (int, error)
    InsertClinicEye(data Clinic_Eye) error
	InsertClinicHeart(data Clinic_heart) error 
	InsertClinicNerve(data Clinic__nerve) error
	InsertClinicBone(data Clinic__bone) error 
	InsertClinicUrinary(data Clinic__urinary) error 
	InsertTreatment(data TreatmentInsert) (int, error)
	FindOrCreateDrugByName(name string) (int, error) 
	InsertTreatmentDrug(td TreatmentDrug) error 
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
	ID				int          `json:"centerId"`
	CenterName		string		 `json:"centerName"  validate:"required"`
	CenterEmail	    string       `json:"centerEmail"    validate:"required,email"`
	CenterCity      string       `json:"centerCity"    validate:"required"`
	PatientNumber   int          `json:"patientNumber"`
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

type Statistics  struct {
	SugarType string `json:"sugarType"`
	Total     int    `json:"total"`



}

type SugarAgeRangeStat struct {
	SugarType string `json:"sugarType"`
	AgeRange  string `json:"ageRange"`
	Total     int    `json:"total"`
}

type BMIStat struct {
	SugarType   string `json:"sugarType"`
	BMICategory string `json:"bmiCategory"` // نحيف، طبيعي، سمنة...
	Total       int    `json:"total"`
}

type CityStat struct {
	City  string `json:"city"`
	Total int    `json:"total"`
}



type SugarStatsResponse struct {
	SugarTypes []*Statistics                    `json:"sugarTypes"`
	MaleCount  int                               `json:"maleCount"`
	FemaleCount int                              `json:"femaleCount"`
	TotalCount  int                              `json:"totalCount"`
	TotalPatientsInSystem int                    `json:"totalPatientsInSystem"`
	SugarAgeRangeStats          []*SugarAgeRangeStat   `json:"sugarAgeRangeStats"`
	TotalSugarAgeRangeStats     []*SugarAgeRangeStat   `json:"totalsugarAgeRangeStats"`
	BMIStats                     []*BMIStat              `json:"bmiStats"`
	CityStats      []*CityStat           `json:"cityStats"`
}





// reiwes .......

// ..........................................


// ............................................................
type Drug struct {
	Name              string     `json:"name"`
	Units             string     `json:"units"`
	Dosage_per_day    string     `json:"dosage_per_day"`

}

type Treatment struct {
	Type      string       	`json:"type"`
    Speed     string       	`json:"speed"`
	Drugs     []Drug        `json:"druges"`
}


type AddReviwePayload struct {
	
	PatientID		                    int		     `json:"patient_id"`



	Address           					string       `json:"title"`
	Weight            					string       `json:"weight"`
	LengthPatient     					string       `json:"length_patient"`
	SugarType         					string       `json:"sugarType"`
    OtherDisease      					string       `json:"otherDisease"`
    HistoryOfFamilyDisease  			string       `json:"historyOfFamilyDisease"`
    HistoryOfDiseaseDetection           string       `json:"historyOfdiseaseDetection"`
    Gender                              string       `json:"gender"`
    Hemoglobin                          string       `json:"hemoglobin"`
    Grease                              string       `json:"grease"`
    UrineAcid                           string       `json:"urineAcid"`
    BloodPressure                       string       `json:"bloodPressure"`
    Cholesterol                         string       `json:"cholesterol"`
    LDL                                 string       `json:"ldl"`
    HDL                                 string       `json:"hdl"`
	Creatine                            string       `json:"creatine"`
	Normal_Glocose                      string       `json:"normal_glocose"`
    Glocose_after_Meal                  string       `json:"Glocose_after_Meal"`
    Triple_Grease                       string       `json:"triple_grease"`
	Hba1c                               string       `json:"hba1c"`
	Coments                             string       `json:"coments"` 

 
	Treatments                           Treatment   `json:"treatments"`


	Has_a_eye_disease                   bool         `json:"has_a_eye_disease"`
	In_kind_disease                     string       `json:"in_kind_disease"`
	Relationship_eyes_with_diabetes     string       `json:"relationship_eyes_with_diabetes"`
	Comments_eyes_clinic                string       `json:"Comments_eyes_clinic"`


	Has_a_heart_disease                 bool         `json:"Has_a_heart_disease"`
	Heart_disease                       string       `json:"Heart_disease"`
	Relationship_heart_with_diabetes    string       `json:"relationship_heart_with_diabetes"`
    Comments_heart_clinic               string       `json:"Comments_heart_clinic"`



    Has_a_nerve_disease                  bool         `json:"Has_a_nerve_disease"`
    Nerve_disease                        string       `json:"nerve_disease"`
	Relationship_nerve_with_diabetes     string       `json:"relationship_nerve_with_diabetes"`
    Comments_nerve_clinic                string       `json:"Comments_nerve_clinic"`




    Has_a_bone_disease                   bool        `json:"Has_a_bone_disease"`
    Bone_disease                        string       `json:"bone_disease"`
	Relationship_bone_with_diabetes     string       `json:"relationship_bone_with_diabetes"`
    Comments_bone_clinic                string       `json:"Comments_bone_clinic"`





    Has_a_urinary_disease                  bool        `json:"Has_a_urinary_disease"`
    Urinary_disease                        string       `json:"urinary_disease"`
	Relationship_urinary_with_diabetes     string       `json:"relationship_urinary_with_diabetes"`
    Comments_urinary_clinic                string       `json:"Comments_urinary_clinic"`



}





type Reviwe struct {
	PatientID		                    int		     
	Address           					string       
	Weight            					string       
	LengthPatient     					string      
	SugarType         					string       
    OtherDisease      					string      
    HistoryOfFamilyDisease  			string       
    HistoryOfDiseaseDetection           string      
    Gender                              string      
    Hemoglobin                          string     
    Grease                              string       
    UrineAcid                           string       
    BloodPressure                       string     
    Cholesterol                         string      
    LDL                                 string   
    HDL                                 string     
	Creatine                            string      
	Normal_Glocose                      string       
    Glocose_after_Meal                  string       
    Triple_Grease                       string       
	Hba1c                               string       
	Coments                             string       

}

type Clinic_Eye struct {
 
	ReviewID                            int
	Has_a_eye_disease                   bool         
	In_kind_disease                     string      
	Relationship_eyes_with_diabetes     string     
	Comments_eyes_clinic                string 

}

type Clinic_heart struct  {

	ReviewID                            int
	Has_a_heart_disease                 bool         
	Heart_disease                       string       
	Relationship_heart_with_diabetes    string      
    Comments_heart_clinic               string      

}

type Clinic__nerve struct  {
	
	ReviewID                            int
    Has_a_nerve_disease                  bool        
    Nerve_disease                        string 
	Relationship_nerve_with_diabetes     string       
    Comments_nerve_clinic                string       

}

type Clinic__bone struct  {

	ReviewID                            int
    Has_a_bone_disease                   bool       
    Bone_disease                        string    
	Relationship_bone_with_diabetes     string    
    Comments_bone_clinic                string   

}


type Clinic__urinary struct  {

	ReviewID                               int
    Has_a_urinary_disease                  bool        
    Urinary_disease                        string  
	Relationship_urinary_with_diabetes     string      
    Comments_urinary_clinic                string     

}


type TreatmentInsert struct {
	ReviewID  int
	Type      string       
    Speed     string      
}


type DrugsName struct {
	
     Name              string 
}

type TreatmentDrug struct {
	TreatmentID   int
	DrugID        int
	DosagePerDay  string
	Units         string
}
