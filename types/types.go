package types

import (

	"time"
)

//postgress strorge

// patients..............................

type PatientStore interface {
	GetPatientByEmail(email string) (*Patient , error)
	GetPatientById(id int) (*Patient , error)
	// GreatePatient(Patient) error
	GreatePatient(patient Patient) ( error)
	// GetPatientDetailsByID(patientID int) (*PatientDetails, error)
	GetUserByEmail(email string) (*UserLoginData, error)
	// GetPatientProfile(id int)(*PatientProfile , error)
	UpdatePatientProfile(patientPayload ParientUpdatePayload) error
	GetSugarTypeStats(centerID int) ([]*Statistics, error)
	GetUpdatePatientProfile(id int) (*GetPatientUpdateProfile , error) 
	GetGenderCounts(centerID int) (int, int, error)
	GetTotalPatientsInSystem() (int, error)
	GetSugarTypeAgeRangeStats(centerID int)([]*SugarAgeRangeStat, error)
	GetSugarTypeAgeRangeStatsAllSystem() ([]*SugarAgeRangeStat, error)
    GetBMIStats(centerID int) ([]*BMIStat, error)
	GetCityStats() ([]*CityStat, error)
	 GetUserByEmailRestPassword(email string) error 
	  UpdatePasswordByEmail(email, newPassword string) error 
    GetUserByID(id int) (*UserLoginData, error) 
	GetReviewsByPatientID(patientID int) ([]ReviewResponseForPatient, error) 

	GetNotificationsByUserID(userID int) ([]NotificationTwo, error)
	UpdateIsReadNotifications(userID int) error
	UpdatePatientBasicInfo(p UpdatePatientInfo , id int) (*UpdatePatientInfo, error)
	UpdatePatientCenterInfo(id int, update UpdatePatientCenterInfo)  (UpdatePatientCenterInfo, error) 
    ChangePatientPassword(patientID int, payload ChangePassword) error 



   GetLoginByID(id int) (*Login, error)
   SetFirstLoginTrue(patientID int) error

   GetPatientsLastMonth() (int, error)
}

type Patient struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Password		string       `json:"password"`
	Age             string		 `json:"age"`
	Phone			string       `json:"phone"`
	IDNumber		string       `json:"id_number"`
	CenterID		int			 `json:"center_id"`
	CreateAt        time.Time    `json:"createAt"`
	FirstLogin      bool         `json:"first_login"`

}

type RegisterPatientPayload struct {
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	Password		string		 `json:"password" validate:"required"`
	IDNumber		string       `json:"id_number" validate:"required"`
}

type ParientUpdatePayload struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	IDNumber		string       `json:"id_number" validate:"required"`

}
type GetPatientUpdateProfile struct {
	FullName		string		 `json:"fullname" validate:"required"`
	Email			string		 `json:"email"    validate:"required,email"`
	Age 			string       `json:"age"    validate:"required"`
	Phone           string	     `json:"phone"    validate:"required"`
	IDNumber		string       `json:"id_number" validate:"required"`
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
    ID         int
	Email      string  `json:"email"`
	Password   string  `json:"-"`
	Role       string  `json:"role"`
}



type ReturnLoggingData struct {
	ID				int			 `json:"id"`
	Name		    string		 `json:"name"`
	Email			string		 `json:"email"`
	Role 			string       `json:"role"`
	FirstLogin     bool          `json:"first_login"`
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
	Id               int	     `json:"id"`
    Date             string      `json:"date"`
}

type CardData struct {
	ID				int			 `json:"id"`
	FullName		string		 `json:"fullname"`
	Email			string		 `json:"email"`
	Age             string		 `json:"age"`
	Phone			string       `json:"phone"`
	IDNumber		string       `json:"id_number"`
	SugarType		string      `json:"sugarType"`
	CreateAt         string      `json:"create_At"`
	Reviews         []Review     `json:"reviews"`          
}


type Login struct {
    ID       int
    Email    string
    Password string
	Role     string
}




//end.......................................
//center....................................

type CenterStore interface {
	GetCenterByName(centerName string) (*Center , error) 
	GetCenterByEmail(centerEmail string) (*Center , error)
	 GetCenterByID(id int) (*Center , error)
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
	// FindOrCreateDrugByName(name string ) (int, error) 
	InsertTreatmentDrug(td TreatmentDrug) error 
	GetReviewsByPatientID(patientID int) ([]Review, error)
	DeleteReviewByID(reviewID int) error 
	GreateLoginFailed(center_login InsertLogin) error 
	GetReviewByID(reviewID int) (*ReviewResponse, error) 
	GreateLoginFailedCenter(center_login InsertLogin) error 
	AddArticle(article Article) error 
	GetArticlesForCenter(centerID int) ([]GetArticles , error)
	GetAllArticles() ([]ReturnAllArticle , error)
	AddActivity(article Article) error
	GetActivitiesForCenter(centerID int) ([]GetArticles , error)
	GetAllActivities() ([]ReturnAllArticle , error)

	Addvideo(video Video) error
	GetVideoForCenter(centerID int) ([]GetVideos , error)
	GetAllVideos() ([]ReturnAllvideo , error)
	DeleteArticleByID(id int) error
	DeleteActivityByID(id int) error
	DeleteVidoeByID(id int) error 


	InsertNotification(n NotificationTwo) error
	GetMedicationStats() (MedicationStats, error)
	InsertMedication(m InsertMedication) (int ,error)
	GetAllMedications(centerID int) ([]GeTMedication, error)
	// UpdateMedicationQuantity(id int, newQuantity int) error
	GetLogsByCenterID(centerID int) ([]MedicationLog, error)
	GetReviewMedicationNames(centerID int) ([]GeTMedicationReview, error) 




	GetMedicationByID(id int) (*GeTMedication, error) 
    UpdateMedicationQuantity(id int, decreaseQuantity string) error


	InsertRecord(r InsertRecord) error
	GetRecordsByCenter(centerID int) ([]Record, error)
	InsertInformation(r InsertInformation) error



	InsertMedicationRequest(m InsertRequestMedicine)  error
}



type CheckIsCenter struct {
	SecretKey        string       `json:"secret_key"`
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
	ID                int        `json:"id,string"`
	Quantity          string     `json:"quantity"`
	Dosage_per_day    string     `json:"dosage_per_day"`

}

type Treatment struct {
	Type      []string       	`json:"type"`
	Drugs     []Drug            `json:"druges"`
}


type AddReviwePayload struct {
	
	PatientID		                    int		     `json:"patient_id"`

	Address           					string       `json:"address"`
	Weight            					string       `json:"weight"`
	LengthPatient     					string       `json:"length_patient"`
	SugarType         					string       `json:"sugarType"`
    OtherDisease      					string       `json:"otherDisease"`
    HistoryOfFamilyDisease  			[]string     `json:"historyOfFamilyDisease"`
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
	Relationship_eyes_with_diabetes     bool       `json:"relationship_eyes_with_diabetes"`
	Comments_eyes_clinic                string       `json:"Comments_eyes_clinic"`


	Has_a_heart_disease                 bool         `json:"has_a_heart_disease"`
	Heart_disease                       string       `json:"Heart_disease"`
	Relationship_heart_with_diabetes    bool       `json:"relationship_heart_with_diabetes"`
    Comments_heart_clinic               string       `json:"Comments_heart_clinic"`



    Has_a_nerve_disease                  bool         `json:"has_a_nerve_disease"`
    Nerve_disease                        string       `json:"nerve_disease"`
	Relationship_nerve_with_diabetes     bool       `json:"relationship_nerve_with_diabetes"`
    Comments_nerve_clinic                string       `json:"Comments_nerve_clinic"`




    Has_a_bone_disease                   bool        `json:"has_a_bone_disease"`
    Bone_disease                        string       `json:"bone_disease"`
	Relationship_bone_with_diabetes     bool       `json:"relationship_bone_with_diabetes"`
    Comments_bone_clinic                string       `json:"Comments_bone_clinic"`





    Has_a_urinary_disease                  bool        `json:"has_a_urinary_disease"`
    Urinary_disease                        string       `json:"urinary_disease"`
	Relationship_urinary_with_diabetes     bool       `json:"relationship_urinary_with_diabetes"`
    Comments_urinary_clinic                string       `json:"Comments_urinary_clinic"`


}





type Reviwe struct {
	PatientID		                    int		     
	Address           					string       
	Weight            					string       
	LengthPatient     					string      
	SugarType         					string       
    OtherDisease      					string      
    HistoryOfFamilyDisease  			[]string       
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
	Relationship_eyes_with_diabetes     bool     
	Comments_eyes_clinic                string 

}

type Clinic_heart struct  {

	ReviewID                            int
	Has_a_heart_disease                 bool         
	Heart_disease                       string       
	Relationship_heart_with_diabetes    bool      
    Comments_heart_clinic               string      

}

type Clinic__nerve struct  {
	
	ReviewID                            int
    Has_a_nerve_disease                  bool        
    Nerve_disease                        string 
	Relationship_nerve_with_diabetes     bool       
    Comments_nerve_clinic                string       

}

type Clinic__bone struct  {

	ReviewID                            int
    Has_a_bone_disease                   bool       
    Bone_disease                        string    
	Relationship_bone_with_diabetes     bool    
    Comments_bone_clinic                string   

}


type Clinic__urinary struct  {

	ReviewID                               int
    Has_a_urinary_disease                  bool        
    Urinary_disease                        string  
	Relationship_urinary_with_diabetes     bool      
    Comments_urinary_clinic                string     

}


type TreatmentInsert struct {
	ReviewID  int
	Type      []string            
}


type DrugsName struct {
	
     Name              string 
}

type TreatmentDrug struct {
	TreatmentID   int
	DrugID        int
	DosagePerDay  string
	Quantity      string
}





// login 
type InsertLogin struct {
	
	Email    string 
	Password string
}


type GetReviwe struct {
    ReviewID                            int
	Address           					string       
	Weight            					string       
	LengthPatient     					string      
	SugarType         					string       
    OtherDisease      					string      
    HistoryOfFamilyDisease  			[]string       
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
	Has_a_eye_disease                   bool         
	In_kind_disease                     string      
	Relationship_eyes_with_diabetes     string     
	Comments_eyes_clinic                string 
	Has_a_heart_disease                 bool         
	Heart_disease                       string       
	Relationship_heart_with_diabetes    string      
    Comments_heart_clinic               string  
	Has_a_nerve_disease                  bool        
    Nerve_disease                        string 
	Relationship_nerve_with_diabetes     string       
    Comments_nerve_clinic                string   
	Has_a_bone_disease                   bool       
    Bone_disease                        string    
	Relationship_bone_with_diabetes     string    
    Comments_bone_clinic                string
	Has_a_urinary_disease                  bool        
    Urinary_disease                        string  
	Relationship_urinary_with_diabetes     string      
    Comments_urinary_clinic                string  
	Speed     string
	Type      []string      

}







//response review vvvv
type DrugR struct {
	// ID             int    `json:"id"`
	Name_arabic    string  `json:"name_arabic"`
	Dosage         string  `json:"dosage"`
	Units_per_box  int     `json:"units_per_box"`
	DosagePerDay   string  `json:"dosage_per_day"`
	Quantity       int     `json:"quantity"`
}

type TreatmentR struct {
	Type  []string `json:"type"`
	Drugs []DrugR  `json:"druges"`
}

type ReviewResponse struct {
	Address                        string   `json:"address"`
	Weight                         string   `json:"weight"`
	LengthPatient                  string   `json:"length_patient"`
	SugarType                      string   `json:"sugarType"`
	OtherDisease                   string   `json:"otherDisease"`
	HistoryOfFamilyDisease         []string `json:"historyOfFamilyDisease"`
	HistoryOfDiseaseDetection      string   `json:"historyOfdiseaseDetection"`
	Gender                         string   `json:"gender"`
	Hemoglobin                     string   `json:"hemoglobin"`
	Grease                         string   `json:"grease"`
	UrineAcid                      string   `json:"urineAcid"`
	BloodPressure                  string   `json:"bloodPressure"`
	Cholesterol                    string   `json:"cholesterol"`
	LDL                            string   `json:"ldl"`
	HDL                            string   `json:"hdl"`
	Creatine                       string   `json:"creatine"`
	NormalGlocose                  string   `json:"normal_glocose"`
	GlocoseAfterMeal               string   `json:"Glocose_after_Meal"`
	TripleGrease                   string   `json:"triple_grease"`
	Hba1c                          string   `json:"hba1c"`
	Coments                        string   `json:"coments"`

	Treatments                     TreatmentR `json:"treatments"`

	HasAEyeDisease                 bool     `json:"has_a_eye_disease"`
	InKindDisease                 string   `json:"in_kind_disease"`
	RelationshipEyesWithDiabetes  bool     `json:"relationship_eyes_with_diabetes"`
	CommentsEyesClinic            string   `json:"Comments_eyes_clinic"`

	HasAHeartDisease              bool     `json:"has_a_heart_disease"`
	HeartDisease                  string   `json:"Heart_disease"`
	RelationshipHeartWithDiabetes bool     `json:"relationship_heart_with_diabetes"`
	CommentsHeartClinic           string   `json:"Comments_heart_clinic"`

	HasANerveDisease              bool     `json:"has_a_nerve_disease"`
	NerveDisease                  string   `json:"nerve_disease"`
	RelationshipNerveWithDiabetes bool     `json:"relationship_nerve_with_diabetes"`
	CommentsNerveClinic           string   `json:"Comments_nerve_clinic"`

	HasABoneDisease               bool     `json:"has_a_bone_disease"`
	BoneDisease                   string   `json:"bone_disease"`
	RelationshipBoneWithDiabetes  bool     `json:"relationship_bone_with_diabetes"`
	CommentsBoneClinic            string   `json:"Comments_bone_clinic"`

	HasAUrinaryDisease            bool     `json:"has_a_urinary_disease"`
	UrinaryDisease                string   `json:"urinary_disease"`
	RelationshipUrinaryWithDiabetes bool   `json:"relationship_urinary_with_diabetes"`
	CommentsUrinaryClinic         string   `json:"Comments_urinary_clinic"`
}










// resetpassword 



type Email struct {
	Email      string    `json:"email"`
}

type OTPResetPass struct {
   OTP         string    `json:"otp"`
}
type ResetPassword struct {
   Email               string    `json:"email"`
   NewPassword         string    `json:"newPassword"`
}






// patient app---------------------------------------------------------------------------
// =========================================================================================
//---------------------------------------------------------end -------------------------------

type ChartData struct {
	Date            string       `json:"date"`
	LDL             string       `json:"ldl"`
	HDL             string       `json:"hdl"`
	NormalGlocose   string       `json:"normal_glocose"`
}

type GetPatientHomeData struct {

	FullName		string		 `json:"fullname"`
	Age             string		 `json:"age"`
	IDNumber		string       `json:"id_number"`
	FirstReviewDate string       `json:"firstReviewDate"`
	ChartData       []ChartData  `json:"chartData"`
	MyCenter        string       `json:"center_name"`
	NextReview      string       `json:"nextReview"`

	MyReviews       []Review     `json:"myReviews"`

}

type ReviewResponseForPatient struct {
	ID                             int      `json:"id"`
	Address                        string   `json:"address"`
	Weight                         string   `json:"weight"`
	LengthPatient                  string   `json:"length_patient"`
	SugarType                      string   `json:"sugarType"`
	OtherDisease                   string   `json:"otherDisease"`
	HistoryOfFamilyDisease         []string `json:"historyOfFamilyDisease"`
	HistoryOfDiseaseDetection      string   `json:"historyOfdiseaseDetection"`
	Gender                         string   `json:"gender"`
	Hemoglobin                     string   `json:"hemoglobin"`
	Grease                         string   `json:"grease"`
	UrineAcid                      string   `json:"urineAcid"`
	BloodPressure                  string   `json:"bloodPressure"`
	Cholesterol                    string   `json:"cholesterol"`
	LDL                            string   `json:"ldl"`
	HDL                            string   `json:"hdl"`
	Creatine                       string   `json:"creatine"`
	NormalGlocose                  string   `json:"normal_glocose"`
	GlocoseAfterMeal               string   `json:"Glocose_after_Meal"`
	TripleGrease                   string   `json:"triple_grease"`
	Hba1c                          string   `json:"hba1c"`
	DateReview                     time.Time    `json:"date_review"`
}



















// articles 
type ArticlePayload struct {
	Title     string		 `json:"title"`
	ShortText string         `json:"shortText"` 
	Desc      string         `json:"desc"`
	ImageURL  string         `json:"imageURL"` 
}

type Article struct {
	CenterID  int
	Title     string		 
	Desc      string     
	ImageURL  string         
	ShortText string          
}
type GetArticles struct {
	ID        int            `json:"id"`
	Title     string		 `json:"title"`
	Desc      string         `json:"desc"`
	CreateAt  string         `json:"createAt"`
	ImageURL  string         `json:"imageURL"` 
	ShortText string         `json:"shortText"` 
}
type AllArticles struct {
	ID        int           
	CenterID  int
	Title     string		 
	Desc      string         
	CreateAt  string   
	ImageURL  string         
	ShortText string          

}
type ReturnAllArticle struct {
	ID        int            `json:"id"`
	CenterName string        `json:"centerName"`
	Title      string		 `json:"title"`
	Desc       string         `json:"desc"`
	CreateAt   string         `json:"createAt"`
	ImageURL  string         `json:"imageURL"` 
	ShortText string         `json:"shortText"` 
}






type VideoPayload struct {
	Title     string		 `json:"title"`
	ShortText string         `json:"shortText"` 
	VideoURL  string         `json:"videoURL"` 
}


type Video struct {
	CenterID  int
	Title     string		    
	ShortText string          
	VideoURL  string         
}

type GetVideos struct {
	ID        int            `json:"id"`
	Title     string		 `json:"title"`
	CreateAt  string         `json:"createAt"`
	VideoURL  string         `json:"videoURL"` 
	ShortText string         `json:"shortText"` 
}


type AllVideos struct {
	ID        int       
	CenterID  int
	Title     string		        
	CreateAt  string   
	VideoURL  string         
	ShortText string          

}
type ReturnAllvideo struct {
	ID        int            `json:"id"`
	CenterName string        `json:"centerName"`
	Title      string		 `json:"title"`
	CreateAt   string         `json:"createAt"`
	VideoURL  string         `json:"videoURL"` 
	ShortText string         `json:"shortText"` 
}









// Notification

type Notification struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt  string `json:"created_at"`
	
}

type NotificationPayload struct {
	
    ReceiverID int    `json:"receiver_id"`
    Message    string `json:"message"`
}


type NotificationTwo struct {
    ID        int       `json:"id"`
    SenderID  int       `json:"sender_id"`
    ReceiverID int      `json:"receiver_id"`
    Message   string    `json:"message"`
    IsRead    bool      `json:"is_read"`
    CreatedAt string `json:"created_at"`
}


type V struct {
	CreatedAt time.Time `json:"created_at"`
}







type Medication struct {
    NameArabic     string    `json:"name_arabic"`
    NameEnglish    string    `json:"name_english"`
    MedicationType string    `json:"medication_type"`
    Dosage         string    `json:"dosage"`
    Quantity       int       `json:"quantity"`
    UnitsPerBox    int       `json:"units_per_box"`
}
type InsertMedication struct {
    NameArabic     string    
    NameEnglish    string   
    MedicationType string   
    Dosage         string   
    Quantity       int       
    UnitsPerBox    int     
	CenterID       int  
}

type GeTMedication struct {
	ID             int       `json:"id"`
    NameArabic     string    `json:"name_arabic"`
    NameEnglish    string    `json:"name_english"`
    MedicationType string    `json:"medication_type"`
    Dosage         string    `json:"dosage"`
    Quantity       int        `json:"quantity"`
    UnitsPerBox    int       `json:"units_per_box"`
	CenterID       int       `json:"center_id"`
	
}

type MedicationStats struct {

    TotalQuantity       int `json:"total_quantity"`         
    TotalUniqueMedTypes int `json:"total_unique_med_types"` 

}

type UpdateNewQuantity struct {
	ID          int   `json:"id"`
	NewQuantity int   `json:"new_quantity"`
}


type MedicationLog struct {
	ID             int       `json:"id"`
	NameArabic     string    `json:"name_arabic"`
	Dosage         string    `json:"dosage"`
	MedicationType string    `json:"medication_type"`
	Quantity       int       `json:"quantity"`
	RequestedAt    time.Time `json:"requested_at"`
}











type UpdatePatientInfo struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	IDNumber  string `json:"id_number"`
	Date      string `json:"date"`
}





type UpdatePatientCenterInfo struct {

	City       string  `json:"city"`
	CenterName string    `json:"centerName"`
}

type ChangePassword struct {

	OldPassword       string  `json:"old_password"`
	NewPassword       string  `json:"new_password"`
}


type GeTMedicationReview struct {
    NameArabic     string    `json:"name_arabic"`
    Dosage         string    `json:"dosage"`
    UnitsPerBox    int       `json:"units_per_box"`
}



type Supervisor struct {
	Email			string		 `json:"email"`
	Role 			string       `json:"role"`
	Token           string       `json:"token"`
}













type InsertRecord struct {
    NameArabic     string     
    MedicationType string   
    Dosage         string   
    Quantity       int       
    CreateAt       string
	ApprovalAt     string  
	CenterID       int  
	Status         string
	RequestID      int
}



type InsertRequestMedicine struct {
    NameArabic     string     
    MedicationType string   
    Dosage         string   
    Quantity       int        
	CenterID       int  
	MedicineID      int
}


type Record struct {
    ID               int
    NameArabic       string
    Dosage           string
    MedicationType   string
    RequestedQuantity int
    CenterID         int
    CreatedAt        string
    ApprovalDate     string
    RecordStatus     string
}



type InsertInformation struct {
    NameArabic     string     
    NameEnglish    string     
    Quantity       int       
	CenterID       int  
	Status         string
	RequestId      int
}


type GetAllInformation struct {
	ID             int        `json:"id"`
    NameArabic     string     `json:"name_arabic"`
    NameEnglish    string     `json:"name_english"`
    Quantity       int        `json:"quantity"`
	CenterName       string   `json:"center_name"`
	CenterCity     string     `json:"center_city"`
	Status         string     `json:"status"`
}
	

type InquirieDetails struct {
	NameArabic     string     `json:"name_arabic"`
    NameEnglish    string     `json:"name_english"`
	CenterName       string   `json:"center_name"`
	CenterCity     string     `json:"center_city"`
	RQuantity       int       `json:"r_quantity"`
	CQuantity       int       `json:"c_quantity"`
	Nop             int       `json:"nop"`
	Request_date    string    `json:"request_date"`
	
}