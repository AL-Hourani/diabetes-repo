package types






type SuperisorStore interface {
     GetAllCenters() ([]*ReturnCenters, error)
	 GetAllInformation() ([]GetAllInformation, error)
	 GetInformationByID(id int) (*InsertInformation, error) 
	 GetCenterByID(id int) (*Center , error)
	 CountPatientsByCenter(centerID int) (int, error)
	 GetMedicationByArabicName(name string , id int) (*GetMedicationRow, error)
	 GetMedicationRequestByID(id int) (*MedicationRequest, error)
     UpdateInformationStatus(id int, newStatus string) error
     UpdateRecordStatusAndApprovalDate(id int, newStatus string) error
     UpdateMedicationQuantity(id int, newQuantity int) error
     GetCentersByCity(cityName string) ([]CenterWithID, error)
     GetPatientCountByCity(cityName string) (int, error)
     GetPatientCountByCityLastMonth(cityName string) (int, error) 
     GetCenterWithMostPatients() (*CenterWithCount, error) 
     ParseMonthYear(input string) (month int, year int, err error) 
     GetPatientReviewsByMonth(month, year int) ([]PatientReview, error) 

      GetPatientCountByCenterLastMonth(id int) (int, error) 
      GetMaleCountByCenter(id int) (int, error)

      GetFemaleCountByCenter(id int) (int, error)
      GetPatientCountByCenterName(id int) (int, error)

}

type ReturnCenters struct {
	ID              int          `json:"id"`
	CenterName		string		 `json:"centerName"`
	CenterEmail	    string       `json:"centerEmail"`
	CenterCity      string       `json:"centerCity"`
	CreateAt        string       `json:"createAt"`
	NumberOfPatient int          `json:"nop"`
}

type GetMedicationRow struct {
    ID            int
    NameArabic    string
    NameEnglish   string
    MedicationType string
    Dosage        string
    Quantity      string
    UnitsPerBox   int
}


type MedicationRequest struct {
    ID               int
    NameArabic       string
    Dosage           string
    MedicationType   string
    RequestedQuantity int
    CenterID         int
    RequestedAt      string  // أو time.Time إذا عمودك DATE
}


type QueryID struct {
    Query_ID    int    `json:"query_id"`
}
type QueryAccepted struct {
    Query_ID    int    `json:"query_id"`
}

type City struct {
    City   string    `json:"city_name"`
}



type CenterWithID struct {
    ID       int      `json:"id"`
    CenterName string  `json:"centerName"`
}

type AllCityInfo struct {
    NumberOfPatientInCity          int `json:"nopic"`
    NumberOfPatientInCityLastMonth int `json:"nopic_lm"`
    ActiveCenter                   []CenterWithID `json:"active_center"`
}
type GenericSuperInfo struct {
    NumberOfPatientInSystem          int `json:"nopis"`
    NumberOfPatientInSystemLastMonth int `json:"nopis_lm"`
    ActiveCities                   []string `json:"active_cities"`
    FirstCenter                    CenterWithCount  `json:"first_center"`
}



type AboutCenterInfo struct {
    NumberOfPatientInCenter          int `json:"nop_in_center"`
    RegisterPatientLastMonth         int  `json:"nop_in_center_lm"`
    NumberOfMale                     int `json:"number_of_male"`
    NumberOfFemale                     int `json:"number_of_Female"`
}
type CenterWithCount struct {
	ID            int
	CenterName    string
	CenterCity    string
	PatientsCount int
}

type TreatmentDrugExel struct {
    DrugName     string
    DosagePerDay string
    Quantity     int
}

type PatientReview struct {
    ReviewID        int
    PatientID       int
    PatientFullName string
    PatientEmail    string
    PatientPhone    string
    Gender          string   // جديد
    SugarType       string   // جديد
    AddressPatient  string
    Wight           string
    LengthPatient   string
    OtherDisease    string
    Hemoglobin      string
    Grease          string
    UrineAcid       string
    BloodPressure   string
    Cholesterol     string
    LDL             string
    HDL             string
    Creatine        string
    NormalClucose   string
    ClucoseAfterMeal string
    TripleGrease    string
    Hba1c           string
    Comments        string
    DateReview      string

    Has_a_eye_disease             string
    In_kind_disease               string
    Relationship_with_diabetes    string
    Comments_eye                  string

    Has_a_heart_disease           string
    Heart_disease                 string
    Relationship_heart_with_diabetes string
    Comments_heart                string

    Has_a_nerve_disease           string
    Nervous_disease               string
    Relationship_nervous_with_diabetes string
    Comments_nervous              string

    Has_a_bone_disease            string
    Bone_disease                  string
    Relationship_bone_with_diabetes string
    Comments_bone                 string

    Has_a_urinary_disease         string
    Urinary_disease               string
    Relationship_urinary_with_diabetes string
    Comments_urinary              string

    TreatmentType  string 
    TreatmentDrugs []TreatmentDrugExel
}


type MonthDown struct {
    MonthDown   string    `json:"date"`
}




type CenterAboutInfor struct {
    
}