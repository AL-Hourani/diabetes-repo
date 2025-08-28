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



func (s *Store) GetCenterByID(id int) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE id=$1",id)
	if err != nil {
		return nil , err
	}
    
defer rows.Close()

	c := new(types.Center)
	for rows.Next() {
		c , err = scanRowIntoGetCenter(rows)
		if err != nil {
			return nil , err
		}
	}

	if c.ID == 0 {
		return nil , fmt.Errorf("center not found")
	}

	return c , nil
}


func scanRowIntoGetCenter(rows *sql.Rows) (*types.Center , error ){
	center := new(types.Center)

	err := rows.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}

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





func (s *Store) GetAllInformation() ([]types.GetAllInformation, error) {
    rows, err := s.db.Query(`
        SELECT 
            name_arabic,
            name_english,
            requested_quantity,
            center_id,
            information_status
        FROM information
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var infos []types.GetAllInformation
    for rows.Next() {
        var info types.GetAllInformation
        var centerID int
        err := rows.Scan(
            &info.NameArabic,
            &info.NameEnglish,
            &info.Quantity,
            &centerID,
            &info.Status,
        )
        if err != nil {
            return nil, err
        }

        center , err := s.GetCenterByID(centerID)
        if err != nil {
            return nil , err
        }
        
        info.CenterName = center.CenterName
        
        infos = append(infos, info)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return infos, nil
}
