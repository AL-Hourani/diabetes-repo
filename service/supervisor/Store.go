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


func (s *Store) GetAllCenters() ([]*types.ReturnCenters, error) {
    rows, err := s.db.Query("SELECT centerName , centerEmail , centerCity , TO_CHAR(createAt, 'DD-MM-YYYY') FROM centers")
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
		&center.CenterName,
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}