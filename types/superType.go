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
     UpdateMedicationQuantity(id int, newQuantity string) error
   
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
    Quantity    int    `json:"quantity"`
}