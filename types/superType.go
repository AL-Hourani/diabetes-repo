package types






type SuperisorStore interface {
     GetAllCenters() ([]*Center, error)
   
}

type ReturnCenters struct {
	CenterName		string		 `json:"centerName"`
	CenterEmail	    string       `json:"centerEmail"`
	CenterCity      string       `json:"centerCity"`
	CreateAt        string       `json:"createAt"`
}