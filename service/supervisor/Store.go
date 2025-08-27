package supervisor

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
// 

func (s *Store) GetAllCenters() ([]*types.ReturnCenters, error) {
    rows, err := s.db.Query(`SELECT c.id, c.centerName, c.centerEmail, c.centerCity,TO_CHAR(c.createAt, 'DD-MM-YYYY'), COUNT(p.id) as patient_count
        FROM centers c
        LEFT JOIN patients p ON p.center_id = c.id
        GROUP BY c.id, c.centerName, c.centerEmail, c.centerCity, c.createAt`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var centers []*types.ReturnCenters

    for rows.Next() {
        c, err := scanRowIntoCenter(rows)
        if err != nil {
            return nil, err
        }
        centers = append(centers, c)
    }

  
    if len(centers) == 0 {
        return nil, fmt.Errorf("no centers found")
    }

    return centers, nil
}


func scanRowIntoCenter(rows *sql.Rows) (*types.ReturnCenters , error ){
	center := new(types.ReturnCenters)

	err := rows.Scan(
        &center.ID,
		&center.CenterName,
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
        &center.NumberOfPatient,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}