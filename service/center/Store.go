package center

import (
	"database/sql"
	"fmt"

	"github.com/AL-Hourani/care-center/types"
)

type Store struct {
	db *sql.DB
}


func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCenterByName(centerName string) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE centerName=?",centerName)
	if err != nil {
		return nil , err
	}

	c := new(types.Center)
	for rows.Next() {
		c , err = scanRowIntoCenter(rows)
		if err != nil {
			return nil , err
		}
	}

	if c.ID == 0 {
		return nil , fmt.Errorf("center not found")
	}

	return c , nil
}


func scanRowIntoCenter(rows *sql.Rows) (*types.Center , error ){
	center := new(types.Center)

	err := rows.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterAddress,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}






func (s *Store)	GreateCenter(center types.Center) error {
	_ , err := s.db.Exec("INSERT INTO centers (centerName ,centerPassword , centerAddress) VALUES (?,?,?)" , center.CenterName , center.CenterPassword , center.CenterAddress)
	if err  != nil {
		return err
	}

	return nil
}




//this is not completed

func (s *Store)	GetPatients(centerId int)([]types.Patient , error) {
	rows , err := s.db.Query("SELECT * FROM patients WHERE center_id=?" ,centerId )
	if err != nil {
		return nil , err
	}

	patients := make([]types.Patient , 0)
	for rows.Next() {
		p , err := scanRowIntoPatients(rows)
		if err != nil {
			return nil , err
		}
		patients = append(patients, *p)
	}

	return patients , nil

}


func scanRowIntoPatients(rows *sql.Rows) (*types.Patient , error ){
	patient := new(types.Patient)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.CenterID,
		&patient.Password,
		&patient.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return patient , nil
}



func (s *Store) GetCenters()([]types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers")
	if err != nil {
		return nil , err
	}

	centers := make([]types.Center , 0)
	for rows.Next() {
		p , err := scanRowIntoCenters(rows)
		if err != nil {
			return nil , err
		}
		centers = append(centers, *p)
	}

	return centers , nil
}


func scanRowIntoCenters(rows *sql.Rows) (*types.Center , error ){
	center := new(types.Center)

	err := rows.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterAddress,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center, nil
}