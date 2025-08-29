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
            id ,
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
            &info.ID,
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
        info.CenterCity = center.CenterCity
        
        infos = append(infos, info)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return infos, nil
}



func (s *Store) GetInformationByID(id int) (*types.InsertInformation, error) {
    row := s.db.QueryRow(`
        SELECT 
            name_arabic,
            name_english,
            requested_quantity,
            center_id,
            information_status,
            request_id
        FROM information
        WHERE id = $1
    `, id)

    var info types.InsertInformation
    err := row.Scan(
        &info.NameArabic,
        &info.NameEnglish,
        &info.Quantity,
        &info.CenterID,
        &info.Status,
        &info.RequestId,
  
    )
    if err != nil {
        return nil, err
    }

    return &info, nil
}



func (s *Store) CountPatientsByCenter(centerID int) (int, error) {
    var count int
    err := s.db.QueryRow(`
        SELECT COUNT(*) 
        FROM patients 
        WHERE center_id = $1
    `, centerID).Scan(&count)

    if err != nil {
        return 0, err
    }
    return count, nil
}


func (s *Store) GetMedicationByArabicName(name string , id int) (*types.GetMedicationRow, error) {
    row := s.db.QueryRow(`
        SELECT 
            id,
            name_arabic,
            name_english,
            medication_type,
            dosage,
            quantity,
            units_per_box
        FROM medications
        WHERE name_arabic = $1 and center_id = $2
    `, name  , id)

    var m types.GetMedicationRow
    err := row.Scan(
        &m.ID,
        &m.NameArabic,
        &m.NameEnglish,
        &m.MedicationType,
        &m.Dosage,
        &m.Quantity,
        &m.UnitsPerBox,
    )
    if err != nil {
        return nil, err
    }

    return &m, nil
}




func (s *Store) GetMedicationRequestByID(id int) (*types.MedicationRequest, error) {
    row := s.db.QueryRow(`
        SELECT 
            id,
            name_arabic,
            dosage,
            medication_type,
            requested_quantity,
            center_id,
            requested_at
        FROM medication_requests
        WHERE  medication_id = $1
    `, id)

    var req types.MedicationRequest
    err := row.Scan(
        &req.ID,
        &req.NameArabic,
        &req.Dosage,
        &req.MedicationType,
        &req.RequestedQuantity,
        &req.CenterID,
        &req.RequestedAt,
    )
    if err != nil {
        return nil, err
    }

    return &req, nil
}








func (s *Store) UpdateInformationStatus(id int, newStatus string) error {
    _, err := s.db.Exec(`
        UPDATE information
        SET information_status = $1
        WHERE id = $2
    `, newStatus, id)

    return err
}


func (s *Store) UpdateRecordStatusAndApprovalDate(id int, newStatus string) error {
    _, err := s.db.Exec(`
        UPDATE records
        SET 
            approval_date = CURRENT_DATE,
            record_status = $1
        WHERE id = $2
    `, newStatus, id)

    return err
}


func (s *Store) UpdateMedicationQuantity(id int, newQuantity string) error {
    _, err := s.db.Exec(`
        UPDATE medications
        SET quantity = $1
        WHERE id = $2
    `, newQuantity, id)

    return err
}
