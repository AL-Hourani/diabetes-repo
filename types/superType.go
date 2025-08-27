package types






type SuperisorStore interface {
     GetAllCenters() ([]*ReturnCenters, error)
   
}

type ReturnCenters struct {
	ID              int          `json:"id"`
	CenterName		string		 `json:"centerName"`
	CenterEmail	    string       `json:"centerEmail"`
	CenterCity      string       `json:"centerCity"`
	CreateAt        string       `json:"createAt"`
	NumberOfPatient int          `json:"nop"`
}