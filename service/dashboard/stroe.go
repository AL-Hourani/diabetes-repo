package dashboard


import (
	"database/sql"


	"github.com/AL-Hourani/care-center/types"
)

type Store struct {
	db *sql.DB
}


func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}



func (s *Store)	SetPatientHealthOverview(patientHealth types.PatientHealthOverview) error {
	_ , err := s.db.Exec("INSERT INTO health () VALUES (?,?,?)" , )
	if err  != nil {
		return err
	}

	return nil
}
