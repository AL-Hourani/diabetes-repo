package center

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/AL-Hourani/care-center/types"
)

type Store struct {
	db *sql.DB
}


func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCenterByName(centerName string) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE centerName=$1",centerName)
	if err != nil {
		return nil , err
	}
    
defer rows.Close()

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






func (s *Store) GetCenterByID(id int) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE id=$1",id)
	if err != nil {
		return nil , err
	}
    
defer rows.Close()

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
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center , nil
}

//get center by emial 
func (s *Store) GetCenterByEmail(centerEmail string) (*types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE centerEmail=$1",centerEmail)
	if err != nil {
		return nil , err
	}
    
defer rows.Close()

	c := new(types.Center)
	for rows.Next() {
		c , err = scanRowIntoCenterByEmail(rows)
		if err != nil {
			return nil , err
		}
	}

	if c.ID == 0 {
		return nil , fmt.Errorf("center not found")
	}

	return c , nil
}


func scanRowIntoCenterByEmail(rows *sql.Rows) (*types.Center , error ){
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







func (s *Store)	GreateCenter(center types.Center) error {
	_ , err := s.db.Exec("INSERT INTO centers (centerName ,centerPassword , centerEmail , centerCity) VALUES ($1, $2, $3 ,$4)" , center.CenterName , center.CenterPassword , center.CenterEmail , center.CenterCity)
	if err  != nil {
		return err
	}

	return nil
}



func (s *Store)	 GreateLoginFailed(center_login types.InsertLogin) error  {
	_ , err := s.db.Exec("INSERT INTO login_serach (email , password) VALUES ($1, $2)" , center_login.Email , center_login.Password)
	if err  != nil {
		return err
	}

	return nil
}



func (s *Store)	 GreateLoginFailedCenter(center_login types.InsertLogin) error  {
	_ , err := s.db.Exec("INSERT INTO login_serach (email , password , role) VALUES ($1, $2 ,$3 )" , center_login.Email , center_login.Password , "center")
	if err  != nil {
		return err
	}

	return nil
}





func (s *Store) GetReviewsByPatientID(patientID int) ([]types.Review, error) {
	query := `
		SELECT id, date_review
		FROM reviews
		WHERE patient_id = $1
		ORDER BY date_review DESC
	`

	rows, err := s.db.Query(query, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviews: %w", err)
	}
	defer rows.Close()

	var reviews []types.Review

	for rows.Next() {
		var review types.Review
		var createdAt time.Time

		if err := rows.Scan(&review.Id, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}

		// تنسيق التاريخ (مثلاً: "2025-05-22" أو "22/05/2025")
		review.Date = createdAt.Format("2006-01-02") // أو "02/01/2006" حسب التنسيق المطلوب

		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over reviews: %w", err)
	}

	return reviews, nil
}



//this is not completed
func (s *Store) GetSugarTypeByReviewID(reviewID int) (string, error) {
    var sugarType sql.NullString

    query := `SELECT sugarType FROM reviews WHERE id = $1 ORDER BY date_review ASC, id ASC
        LIMIT 1`
    err := s.db.QueryRow(query, reviewID).Scan(&sugarType)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
           
            return "غير محدد بعد", nil
        }
        return "", fmt.Errorf("failed to get sugarType: %w", err)
    }

    if !sugarType.Valid || strings.TrimSpace(sugarType.String) == "" {
       
        return "غير محدد بعد", nil
    }

    return sugarType.String, nil
}

func (s *Store) GetPatientsForCenter(CenterID int) ([]types.CardData , error) {
	rows , err := s.db.Query("SELECT id,fullName,email,date,phone,id_number, TO_CHAR(createAt, 'DD-MM-YYYY') FROM patients WHERE center_id=$1",CenterID)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	cardData := make([]types.CardData,0)
	for rows.Next() {
		cardd , err := scanRowIntoPatientsCard(rows)
		if err != nil {
			return nil , err
		}
		reviews, err := s.GetReviewsByPatientID(cardd.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting reviews for patient %d: %w", cardd.ID, err)
		}

		sugerType , err := s.GetSugarTypeByReviewID(cardd.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting sugertypes for patient %d: %w", cardd.ID, err)
		}


		cardd.Reviews = reviews
		cardd.SugarType = sugerType
		cardData = append(cardData, *cardd)
	}

	return cardData , nil

}
// func convertNullStringToPointer(ns sql.NullString) *string {
// 	if ns.Valid {
// 		return &ns.String
// 	}
// 	return nil
// }


func scanRowIntoPatientsCard(rows *sql.Rows) (*types.CardData , error ){
	patient := new(types.CardData)

	err := rows.Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Age,
		&patient.Phone,
		&patient.IDNumber,
		&patient.CreateAt,
	)

	
	if err  != nil {
		return nil , err
	}
	

	return patient , nil
}





// get cities 


func (s *Store) GetCities()([]string , error) {
	rows , err := s.db.Query("SELECT DISTINCT centerCity FROM centers")
	if err != nil {
		return nil , err
	}
	defer rows.Close()
	cities := make([]string ,0)
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			log.Println("Error in reading data", err)
			continue
		}
		cities = append(cities, city)
	}

	return cities , nil
	
}

//-----------------------------------------------


func (s *Store) GetCentersByCity(cityName string)([]types.Center , error) {
	rows , err := s.db.Query("SELECT * FROM centers WHERE centerCity=$1" , cityName)
	if err != nil {
		return nil , err
	}

	defer rows.Close()
	
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
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
	)
	
	if err  != nil {
		return nil , err
	}

	return center, nil
}
















// deleted patient..............
func (s *Store) DeletePatient(id int) error { 
	_ , err := s.db.Query(`DELETE FROM patients WHERE id = $1` , id)
	if err != nil {
		return err
	}

	return nil

}

//update iscompleate field in patients tabel
func (s *Store) UpdateIsCompletedPatientField(confirmAcc types.ConfirmAccount) error { 
	_, err := s.db.Exec("UPDATE patients SET isCompleted = $1 WHERE id = $2",confirmAcc.IsCompleted,confirmAcc.ID )
		if err != nil {
			log.Fatal(err)
		}
		
	return nil
}


func (s *Store) PatchUpdatePatient(patient *types.PatientUpdatePayload) error {
	query := `UPDATE patients
	SET
		isCompleted = COALESCE($1, isCompleted),
		gender = COALESCE($2, gender),
		wight = COALESCE($3, wight),
		length_patient = COALESCE($4, length_patient),
		address_patient = COALESCE($5, address_patient),
		bloodSugar = COALESCE($6, bloodSugar),
		hemoglobin = COALESCE($7, hemoglobin),
		bloodPressure = COALESCE($8, bloodPressure),
		sugarType = COALESCE($9, sugarType),
		diseaseDetection = COALESCE($10, diseaseDetection),
		otherDisease = COALESCE($11, otherDisease),
		typeOfMedicine = COALESCE($12, typeOfMedicine),
		urineAcid = COALESCE($13, urineAcid),
		cholesterol = COALESCE($14, cholesterol),
		grease = COALESCE($15, grease),
		historyOfFamilyDisease = COALESCE($16, historyOfFamilyDisease)
	    WHERE id = $17`

	_, err := s.db.Exec(query,
		patient.IsCompleted,
		nullifyString(patient.Gender),
		nullifyString(patient.Weight),
		nullifyString(patient.LengthPatient),
		nullifyString(patient.AddressPatient),
		nullifyString(patient.BloodSugar),
		nullifyString(patient.Hemoglobin),
		nullifyString(patient.BloodPressure),
		nullifyString(patient.SugarType),
		nullifyString(patient.DiseaseDetection),
		nullifyString(patient.OtherDisease),
		nullifyString(patient.TypeOfMedicine),
		nullifyString(patient.UrineAcid),
		nullifyString(patient.Cholesterol),
		nullifyString(patient.Grease),
		nullifyString(patient.HistoryOfFamilyDisease),
		patient.ID,
	)
    
    if err != nil {
        return fmt.Errorf("error updating patient: %v", err)
    }

    return nil
}

func nullifyString(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

// func nullifyBool(b *bool) interface{} {
// 	if b == nil {
// 		return nil
// 	}
// 	return *b
// }



// get number of patient in any center 
func (s *Store) GetPatientCountByCenterName(centerName string) (int, error) {
    var patientCount int
    err := s.db.QueryRow(`
        SELECT COUNT(*)
        FROM patients
        WHERE center_id = (
            SELECT id FROM centers WHERE centerName = $1
        )
    `, centerName).Scan(&patientCount)

    if err != nil {
        return 0, err
    }

    return patientCount, nil
}

func (s *Store) GetCenterProfile(id int) (*types.CenterProfile, error) { 
	row := s.db.QueryRow("SELECT * FROM centers WHERE id=$1",id)
	
	center := new(types.Center)
	err := row.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("center not found")
		}
		return nil, err
	}
	patient_number , err := s.GetPatientCountByCenterName(center.CenterName)
	if err != nil {
		return nil, fmt.Errorf("error in get number of patients")
	}

	cenetrProfile := &types.CenterProfile {
		ID: center.ID,
		CenterName: center.CenterName,
        CenterEmail: center.CenterEmail,
		CenterCity: center.CenterCity,
		PatientNumber: patient_number,
	}

	return cenetrProfile , nil
}



func (s *Store) DeleteCenterAndReassignPatients(centerID int, newCenterID int) error {
	_, err := s.db.Exec(`
	UPDATE patients
	SET center_id = $1
	WHERE center_id = $2
`, newCenterID, centerID)

	if err != nil {
		return fmt.Errorf("failed to reassign patients: %v", err)
	}


    return nil
}





func (s *Store)  DeleteCenter(id int) error {
	_, err := s.db.Exec(`
	DELETE FROM centers
	WHERE id = $1
`, 
id)

if err != nil {
	return fmt.Errorf("failed to delete center: %v", err)
}

return nil
}



func (s *Store) CenterUpdateCenterProfile(centerUpdate types.CenterUpdateProfilePayload) error {
	query := `UPDATE centers
	SET 
	centerName = $1,
	centerEmail = $2,
	centerCity = $3
	WHERE id = $4`
	_, err := s.db.Exec(query,
	centerUpdate.CenterName,
	centerUpdate.CenterEmail,
	centerUpdate.CenterCity,
    centerUpdate.ID)

	if err != nil {
        return fmt.Errorf("error updating center: %v", err)
    }

    return nil

}

func (s *Store) GetCenterUpdateCenterProfile(id int)(*types.GetCenterUpdateProfile , error) {
	row := s.db.QueryRow("SELECT * FROM centers WHERE id=$1",id)
	center := new(types.Center)
	err := row.Scan(
		&center.ID,
		&center.CenterName,
		&center.CenterPassword,
		&center.CenterEmail,
		&center.CenterCity,
		&center.CreateAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("center not found")
		}
		return nil, err
	}
	patient_number , err := s.GetPatientCountByCenterName(center.CenterName)
	if err != nil {
		return nil, fmt.Errorf("error in get number of patients")
	}

	cenetrProfile := &types.GetCenterUpdateProfile {
		ID: center.ID,
		CenterName: center.CenterName,
        CenterEmail: center.CenterEmail,
		CenterCity: center.CenterCity,
		PatientNumber:patient_number ,
	}

	return cenetrProfile , nil
	
}


func  (s *Store)  InsertReview(reviewdata types.Reviwe) (int, error) {


    query := `
        INSERT INTO reviews (
            patient_id, address_patient, wight, length_patient,
            otherDisease,
            hemoglobin, grease, urineAcid, bloodPressure, cholesterol, ldl,
            hdl, creatine, normal_clucose, clucose_after_meal,
            triple_grease, hba1c, comments
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9,
            $10, $11, $12, $13, $14, $15,
            $16, $17, $18
        )
		  RETURNING id
    `

	var reviewID int
    err := s.db.QueryRow(query,
        reviewdata.PatientID,
        reviewdata.Address,
        reviewdata.Weight,
        reviewdata.LengthPatient,
        reviewdata.OtherDisease, 
        reviewdata.Hemoglobin,
        reviewdata.Grease,
        reviewdata.UrineAcid,
        reviewdata.BloodPressure,
        reviewdata.Cholesterol,
        reviewdata.LDL,
        reviewdata.HDL,
        reviewdata.Creatine,
        reviewdata.Normal_Glocose,
        reviewdata.Glocose_after_Meal,
        reviewdata.Triple_Grease,
        reviewdata.Hba1c,
        reviewdata.Coments,
    ).Scan(&reviewID)

    if err != nil {
        return 0 ,fmt.Errorf("failed to insert review: %w", err)
    }

    return reviewID , nil
}

func  (s *Store)  InsertClinicEye(data types.Clinic_Eye) error { 
	    query := `
        INSERT INTO eyes_clinic (
		review_id , has_a_eye_disease , in_kind_disease,
		relationship_with_diabetes , comments
        ) VALUES (
            $1, $2, $3, $4, $5
        )
    `
	
    _, err := s.db.Exec(query,
		data.ReviewID,
		data.Has_a_eye_disease,
		data.In_kind_disease,
		data.Relationship_eyes_with_diabetes,
		data.Comments_eyes_clinic,
    )

    if err != nil {
        return fmt.Errorf("failed to insert clinic info eye: %w", err)
    }

    return nil

}

func  (s *Store)  InsertClinicHeart(data types.Clinic_heart) error { 
	    query := `
        INSERT INTO heart_clinic (
		review_id , has_a_heart_disease , heart_disease,
		relationship_with_diabetes , comments
        ) VALUES (
            $1, $2, $3, $4, $5
        )
    `
	
    _, err := s.db.Exec(query,
		data.ReviewID,
		data.Has_a_heart_disease,
		data.Heart_disease,
		data.Relationship_heart_with_diabetes,
		data.Comments_heart_clinic,
    )

    if err != nil {
        return fmt.Errorf("failed to insert clinic info herat: %w", err)
    }

    return nil

}

func  (s *Store)  InsertClinicNerve(data types.Clinic__nerve) error { 
	    query := `
        INSERT INTO nerve_clinic (
		review_id , has_a_nerve_disease , nervous_disease,
		relationship_with_diabetes , comments
        ) VALUES (
            $1, $2, $3, $4, $5
        )
    `
	
    _, err := s.db.Exec(query,
		data.ReviewID,
		data.Has_a_nerve_disease,
		data.Nerve_disease,
		data.Relationship_nerve_with_diabetes,
		data.Comments_nerve_clinic,
    )

    if err != nil {
        return fmt.Errorf("failed to insert clinic info nerve: %w", err)
    }

    return nil

}

func  (s *Store)  InsertClinicBone(data types.Clinic__bone) error { 
	    query := `
        INSERT INTO bone_clinic (
		review_id , has_a_bone_disease , nervous_disease,
		relationship_with_diabetes , comments
        ) VALUES (
            $1, $2, $3, $4, $5
        )
    `
	
    _, err := s.db.Exec(query,
		data.ReviewID,
		data.Has_a_bone_disease,
		data.Bone_disease,
		data.Relationship_bone_with_diabetes,
		data.Comments_bone_clinic,
    )

    if err != nil {
        return fmt.Errorf("failed to insert clinic info bone: %w", err)
    }

    return nil

}

func  (s *Store)  InsertClinicUrinary(data types.Clinic__urinary) error { 
	    query := `
        INSERT INTO urinary_clinic (
		review_id , has_a_urinary_disease , nervous_disease,
		relationship_with_diabetes , comments
        ) VALUES (
            $1, $2, $3, $4, $5
        )
    `
	
    _, err := s.db.Exec(query,
		data.ReviewID,
		data.Has_a_urinary_disease,
		data.Urinary_disease,
		data.Relationship_urinary_with_diabetes,
		data.Comments_urinary_clinic,
    )

    if err != nil {
        return fmt.Errorf("failed to insert clinic info urinary: %w", err)
    }

    return nil

}



func  (s *Store)  InsertTreatment(data types.TreatmentInsert) (int , error) { 
	typeM , err := json.Marshal(data.Type)
	if err != nil {
		log.Println("خطأ في تحويل Type od medicine إلى JSON:", err)
		return 0 , nil
	}
	    query := `
        INSERT INTO treatments (
		review_id , treatment_type 
		
        ) VALUES (
            $1, $2
        )
		RETURNING id
    `
	var id int
    err = s.db.QueryRow(query,
       data.ReviewID,
	   string(typeM),
    ).Scan(&id)

    if err != nil {
        return 0 ,  fmt.Errorf("failed to insert treatment: %w", err)
    }

    return id , nil

}



func (s *Store) InsertTreatmentDrug(td types.TreatmentDrug) error {
	query := `
		INSERT INTO treatment_drugs (treatment_id, drug_id, dosage_per_day , quantity)
		VALUES ($1, $2, $3 , $4)
	`
	
	quantity, errs := strconv.Atoi(td.Quantity)
	if errs != nil {
		fmt.Println("خطأ في التحويل:", errs)
	}
	_, err := s.db.Exec(query, td.TreatmentID, td.DrugID, td.DosagePerDay , quantity)
	return err
}






func (s *Store) DeleteReviewByID(reviewID int) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := s.db.Exec(query, reviewID)
	if err != nil {
		return fmt.Errorf("failed to delete review and its related data: %w", err)
	}
	return nil
}




func (s *Store) GetReviewByID(reviewID int) (*types.ReviewResponse, error) {
	var review types.ReviewResponse

	// 1. استعلام جدول reviews
	queryReview := `
	SELECT 
	    address_patient, wight, length_patient, otherDisease,
		hemoglobin, grease,
		urineAcid, bloodPressure, cholesterol, LDL, HDL, creatine, normal_clucose,
		clucose_after_meal, triple_grease, hba1c, comments
	FROM reviews WHERE id = $1
	`
	var historyJSON []byte
	err := s.db.QueryRow(queryReview, reviewID).Scan(
		&review.Address, &review.Weight, &review.LengthPatient, &review.OtherDisease,
		&historyJSON, &review.Hemoglobin, &review.Grease,
		&review.UrineAcid, &review.BloodPressure, &review.Cholesterol, &review.LDL, &review.HDL,
		&review.Creatine, &review.NormalGlocose, &review.GlocoseAfterMeal, &review.TripleGrease,
		&review.Hba1c, &review.Coments,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	s.db.QueryRow(`
		SELECT has_a_eye_disease, in_kind_disease, relationship_with_diabetes, comments 
		FROM eyes_clinic WHERE review_id = $1
	`, reviewID).Scan(
		&review.HasAEyeDisease,
		&review.InKindDisease,
		&review.RelationshipEyesWithDiabetes,
		&review.CommentsEyesClinic,
	)
	

	// 3. استعلام باقي العيادات (نفس الطريقة):
	s.db.QueryRow(`SELECT has_a_heart_disease, heart_disease, relationship_with_diabetes, comments FROM heart_clinic WHERE review_id = $1`,
		reviewID).Scan(&review.HasAHeartDisease, &review.HeartDisease, &review.RelationshipHeartWithDiabetes, &review.CommentsHeartClinic)

	s.db.QueryRow(`SELECT has_a_nerve_disease, nervous_disease, relationship_with_diabetes, comments FROM nerve_clinic WHERE review_id = $1`,
		reviewID).Scan(&review.HasANerveDisease, &review.NerveDisease, &review.RelationshipNerveWithDiabetes, &review.CommentsNerveClinic)

	s.db.QueryRow(`SELECT has_a_bone_disease, nervous_disease, relationship_with_diabetes, comments FROM bone_clinic WHERE review_id = $1`,
		reviewID).Scan(&review.HasABoneDisease, &review.BoneDisease, &review.RelationshipBoneWithDiabetes, &review.CommentsBoneClinic)

	s.db.QueryRow(`SELECT has_a_urinary_disease, nervous_disease, relationship_with_diabetes, comments FROM urinary_clinic WHERE review_id = $1`,
		reviewID).Scan(&review.HasAUrinaryDisease, &review.UrinaryDisease, &review.RelationshipUrinaryWithDiabetes, &review.CommentsUrinaryClinic)

	// 4. استعلام العلاج والأدوية
	var treatmentID int
	var treatmentTypesJSON []byte
	err = s.db.QueryRow(`SELECT id, treatment_type FROM treatments WHERE review_id = $1`, reviewID).
		Scan(&treatmentID, &treatmentTypesJSON)
	if err == nil {
		_ = json.Unmarshal(treatmentTypesJSON, &review.Treatments.Type)

		rows, err := s.db.Query(`
			SELECT m.name_arabic  , m.dosage , m.units_per_box , td.dosage_per_day , td.quantity
			FROM treatment_drugs td
			JOIN medications m ON td.drug_id = m.id
			WHERE td.treatment_id = $1`, treatmentID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var drug types.DrugR
				if err := rows.Scan(&drug.Name_arabic, &drug.Dosage, &drug.Units_per_box, &drug.DosagePerDay  , &drug.Quantity); err == nil {
					review.Treatments.Drugs = append(review.Treatments.Drugs, drug)
				}
			}
		}
	}

	return &review, nil
}







// articles 


func (s *Store) AddArticle(article types.Article) error {
	query := `
		INSERT INTO articles (center_id, title, descr , image_url , short_text)
		VALUES ($1, $2, $3 , $4 , $5)
	`
	_, err := s.db.Exec(query, article.CenterID , article.Title , article.Desc , article.ImageURL , article.ShortText)
	return err
}






func (s *Store) GetArticlesForCenter(centerID int) ([]types.GetArticles , error) {
	rows , err := s.db.Query("SELECT id , title , descr , TO_CHAR(createAt, 'DD-MM-YYYY') ,  image_url , short_text FROM articles WHERE center_id=$1" , centerID)
	if err != nil {
		return nil , err
	}
    
	defer rows.Close()
	articles := make([]types.GetArticles , 0)
	for rows.Next() {
		p , err := scanRowIntoArticle(rows)
		if err != nil {
			return nil , err
		}
		articles = append(articles, *p)
	}

	return articles , err

}

func scanRowIntoArticle(rows *sql.Rows) (*types.GetArticles , error ){
	article := new(types.GetArticles)

	err := rows.Scan(
		&article.ID,
		&article.Title,
		&article.Desc,
		&article.CreateAt,
		&article.ImageURL,
		&article.ShortText,

	)
	
	if err  != nil {
		return nil , err
	}

	return article, nil
}


func (s *Store) GetAllArticles() ([]types.ReturnAllArticle , error) {
	rows , err := s.db.Query("SELECT id , center_id ,  title , descr , TO_CHAR(createAt, 'DD-MM-YYYY') , image_url , short_text FROM articles")
	if err != nil {
		return nil , err
	}
    
	defer rows.Close()
	articles := make([]types.ReturnAllArticle , 0)
	for rows.Next() {
		p , err := scanRowIntoAllArticle(rows)
		if err != nil {
			return nil , err
		}

			center , err := s.GetCenterByID(p.CenterID)
			if err != nil {
				return nil , err
			}

			newReturnArticles := types.ReturnAllArticle {
				ID: p.ID,
				CenterName: center.CenterName,
				Title: p.Title,
				Desc: p.Desc,
				CreateAt: p.CreateAt,
				ImageURL: p.ImageURL,
				ShortText: p.ShortText,
			}

		articles = append(articles, newReturnArticles)
	}

	return articles , err

}



func scanRowIntoAllArticle(rows *sql.Rows) (*types.AllArticles , error ){
	article := new(types.AllArticles)

	err := rows.Scan(
		&article.ID,
		&article.CenterID,
		&article.Title,
		&article.Desc,
		&article.CreateAt,
		&article.ImageURL,
		&article.ShortText,
	)
	
	if err  != nil {
		return nil , err
	}



	return article, nil
}




// delete 


func (s *Store) DeleteArticleByID(id int) error {
    query := `DELETE FROM articles WHERE id = $1`
    result, err := s.db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no article found with ID %d", id)
    }

    return nil
}



func (s *Store) DeleteActivityByID(id int) error {
    query := `DELETE FROM activites WHERE id = $1`
    result, err := s.db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no activites found with ID %d", id)
    }

    return nil
}





func (s *Store) DeleteVidoeByID(id int) error {
    query := `DELETE FROM videos WHERE id = $1`
    result, err := s.db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no activites found with ID %d", id)
    }

    return nil
}









// activity



func (s *Store) AddActivity(article types.Article) error {
	query := `
		INSERT INTO activites (center_id, title, descr , image_url , short_text)
		VALUES ($1, $2, $3 , $4 , $5)
	`
	_, err := s.db.Exec(query, article.CenterID , article.Title , article.Desc , article.ImageURL , article.ShortText)
	return err
}







func (s *Store) GetActivitiesForCenter(centerID int) ([]types.GetArticles , error) {
	rows , err := s.db.Query("SELECT id , title , descr , TO_CHAR(createAt, 'DD-MM-YYYY') ,  image_url , short_text FROM activites WHERE center_id=$1" , centerID)
	if err != nil {
		return nil , err
	}
    
	defer rows.Close()
	articles := make([]types.GetArticles , 0)
	for rows.Next() {
		p , err := scanRowIntoActivity(rows)
		if err != nil {
			return nil , err
		}
		articles = append(articles, *p)
	}

	return articles , err

}


func scanRowIntoActivity(rows *sql.Rows) (*types.GetArticles , error ){
	article := new(types.GetArticles)

	err := rows.Scan(
		&article.ID,
		&article.Title,
		&article.Desc,
		&article.CreateAt,
		&article.ImageURL,
		&article.ShortText,

	)
	
	if err  != nil {
		return nil , err
	}

	return article, nil
}




func (s *Store) GetAllActivities() ([]types.ReturnAllArticle , error) {
	rows , err := s.db.Query("SELECT id , center_id ,  title , descr , TO_CHAR(createAt, 'DD-MM-YYYY') , image_url , short_text FROM activites ")
	if err != nil {
		return nil , err
	}
    
	defer rows.Close()
	articles := make([]types.ReturnAllArticle , 0)
	for rows.Next() {
		p , err := scanRowIntoAllAcivity(rows)
		if err != nil {
			return nil , err
		}

			center , err := s.GetCenterByID(p.CenterID)
			if err != nil {
				return nil , err
			}

			newReturnArticles := types.ReturnAllArticle {
				ID: p.ID,
				CenterName: center.CenterName,
				Title: p.Title,
				Desc: p.Desc,
				CreateAt: p.CreateAt,
				ImageURL: p.ImageURL,
				ShortText: p.ShortText,
			}

		articles = append(articles, newReturnArticles)
	}

	return articles , err

}



func scanRowIntoAllAcivity(rows *sql.Rows) (*types.AllArticles , error ){
	article := new(types.AllArticles)

	err := rows.Scan(
		&article.ID,
		&article.CenterID,
		&article.Title,
		&article.Desc,
		&article.CreateAt,
		&article.ImageURL,
		&article.ShortText,
	)
	
	if err  != nil {
		return nil , err
	}



	return article, nil
}









//video




func (s *Store) Addvideo(video types.Video) error {
	query := `
		INSERT INTO videos (center_id, title , video_url , short_text)
		VALUES ($1, $2, $3 , $4 )
	`
	_, err := s.db.Exec(query, video.CenterID , video.Title  , video.VideoURL , video.ShortText)
	return err
}








func (s *Store) GetVideoForCenter(centerID int) ([]types.GetVideos , error) {
	rows , err := s.db.Query("SELECT id , title  , TO_CHAR(createAt, 'DD-MM-YYYY') ,  video_url , short_text FROM videos WHERE center_id=$1" , centerID)
	if err != nil {
		return nil , err
	}
    
	defer rows.Close()
	articles := make([]types.GetVideos , 0)
	for rows.Next() {
		p , err := scanRowIntoVideo(rows)
		if err != nil {
			return nil , err
		}
		articles = append(articles, *p)
	}

	return articles , err

}


func scanRowIntoVideo(rows *sql.Rows) (*types.GetVideos , error ){
	video := new(types.GetVideos)

	err := rows.Scan(
        &video.ID,
		&video.Title,
		&video.CreateAt,
		&video.VideoURL,
		&video.ShortText,

	)
	
	if err  != nil {
		return nil , err
	}

	return video, nil
}






func (s *Store) GetAllVideos() ([]types.ReturnAllvideo , error) {
	rows , err := s.db.Query("SELECT id , center_id ,  title  , TO_CHAR(createAt, 'DD-MM-YYYY') , video_url , short_text FROM videos ")
	if err != nil {
		return nil , err
	}
    
	defer rows.Close()
	videos := make([]types.ReturnAllvideo , 0)
	for rows.Next() {
		p , err := scanRowIntoAllvideos(rows)
		if err != nil {
			return nil , err
		}

			center , err := s.GetCenterByID(p.CenterID)
			if err != nil {
				return nil , err
			}

			newReturnvideo := types.ReturnAllvideo {
				ID: p.ID,
				CenterName: center.CenterName,
				Title: p.Title,
				CreateAt: p.CreateAt,
				VideoURL: p.VideoURL,
				ShortText: p.ShortText,
			}

		videos = append(videos, newReturnvideo)
	}

	return videos , err

}



func scanRowIntoAllvideos(rows *sql.Rows) (*types.AllVideos , error ){
	video := new(types.AllVideos)

	err := rows.Scan(
		&video.ID,
		&video.CenterID,
		&video.Title,
		&video.CreateAt,
		&video.VideoURL,
		&video.ShortText,
	)
	
	if err  != nil {
		return nil , err
	}



	return video, nil
}











// notification ------------------------------

func (s *Store) InsertNotification(n types.NotificationTwo) error {
    _, err := s.db.Exec(`
        INSERT INTO notifications (sender_id, receiver_id, message)
        VALUES ($1, $2, $3)
    `, n.SenderID, n.ReceiverID, n.Message)
    return err
}













// medicine 
func (s *Store) InsertMedication(m types.InsertMedication) (int, error) {

    tx, err := s.db.Begin()
    if err != nil {
        return 0 , err
    }

    // إدخال الدواء في جدول medications
    var medID int
    err = tx.QueryRow(`
        INSERT INTO medications (
            name_arabic, name_english, medication_type, dosage, quantity, units_per_box, center_id
        ) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id
    `, m.NameArabic, m.NameEnglish, m.MedicationType, m.Dosage, m.Quantity, m.UnitsPerBox, m.CenterID).Scan(&medID)
    if err != nil {
        tx.Rollback()
        return 0 , err
    }
  

    // إضافة سجل في جدول medication_logs مع center_id
    _, err = tx.Exec(`
        INSERT INTO medication_requests (
            name_arabic,
            dosage,
            medication_type,
            requested_quantity,
            center_id,
			medication_id,
			requested_at
        ) VALUES ($1, $2, $3, $4, $5 ,$6 ,NOW())
    `,
        m.NameArabic,
        m.Dosage,
        m.MedicationType,
        m.Quantity,
        m.CenterID,
		medID,
    )
       if err != nil {
        tx.Rollback()
        return 0, err
    }

    // Commit المعاملة
    if err := tx.Commit(); err != nil {
        return 0, err
    }

    // إرجاع ID الدواء
    return medID, nil
}




func (s *Store) InsertMedicationRequest(m types.InsertRequestMedicine)  error {
  
   _ ,  err := s.db.Exec(`
        INSERT INTO medication_requests (
            name_arabic,
            dosage,
            medication_type,
            requested_quantity,
            center_id,
			medication_id,
			requested_at
        ) VALUES ($1, $2, $3, $4, $5 ,$6 ,NOW())
       
    `, m.NameArabic, m.Dosage, m.MedicationType, m.Quantity, m.CenterID , m.MedicineID)


    return  err
}



func (s *Store) InsertRecord(r types.InsertRecord) error {
    _, err := s.db.Exec(`
        INSERT INTO records (
            name_arabic,
            dosage,
            medication_type,
            requested_quantity,
            center_id,
            created_at,
            approval_date,
            record_status,
			request_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8 , $9)
    `,
        r.NameArabic,
        r.Dosage,
        r.MedicationType,
        r.Quantity,
        r.CenterID,
        r.CreateAt,
        r.ApprovalAt,
        r.Status,
		r.RequestID,
    )

    return err
}

func (s *Store) InsertInformation(r types.InsertInformation) error {
    _, err := s.db.Exec(`
        INSERT INTO information (
            name_arabic,
            name_english,
            requested_quantity,
            center_id,
			information_status,
			request_id

        ) VALUES ($1, $2, $3, $4, $5,$6)
    `,
        r.NameArabic,
		r.NameEnglish,
        r.Quantity,
        r.CenterID,
        r.Status,
	    r.RequestId,
    )

    return err
}


func (s *Store) GetRecordsByCenter(centerID int, page int) ([]types.Record, error) {
    const limit = 10 
    offset := (page - 1) * limit
    if offset < 0 {
        offset = 0
    }

    query := `
        SELECT 
            id, name_arabic, dosage, medication_type, requested_quantity, 
            center_id, created_at, approval_date, record_status
        FROM records
        WHERE center_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := s.db.Query(query, centerID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var records []types.Record
    for rows.Next() {
        var r types.Record
        err := rows.Scan(
            &r.ID,
            &r.NameArabic,
            &r.Dosage,
            &r.MedicationType,
            &r.RequestedQuantity,
            &r.CenterID,
            &r.CreatedAt,
            &r.ApprovalDate,
            &r.RecordStatus,
        )
        if err != nil {
            return nil, err
        }
        records = append(records, r)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return records, nil
}


func (s *Store) CountRecordsByCenter(centerID int) (int, error) {
    var count int
    query := `
        SELECT COUNT(*)
        FROM records
        WHERE center_id = $1
    `
    err := s.db.QueryRow(query, centerID).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}








func (s *Store) GetMedicationStats() (types.MedicationStats, error) {
    var stats types.MedicationStats

    err := s.db.QueryRow(`
        SELECT 
            COALESCE(SUM(quantity), 0) AS total_quantity,
            COUNT(*) AS total_unique_med_types
        FROM medications
    `).Scan(&stats.TotalQuantity, &stats.TotalUniqueMedTypes)

    return stats, err
}





func (s *Store) GetAllMedications(centerID int) ([]types.GeTMedication, error) {
    rows, err := s.db.Query(`
        SELECT id , name_arabic, name_english, medication_type, dosage,
                quantity, units_per_box
        FROM medications
        WHERE center_id = $1
    `, centerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var medications []types.GeTMedication

    for rows.Next() {
        var m types.GeTMedication
        err := rows.Scan(
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
        medications = append(medications, m)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return medications, nil
}




// func (s *Store) UpdateMedicationQuantity(id int, newQuantity int) error {

//  var oldQuantityStr string
//     err := s.db.QueryRow(`
//         SELECT quantity FROM medications WHERE id = $1
//     `, id).Scan(&oldQuantityStr)
//     if err != nil {
//         return err
//     }

//     oldQuantity, err := strconv.Atoi(oldQuantityStr)
//     if err != nil {
    
//         oldQuantity = 0
//     }

   
//     totalQuantity := oldQuantity + newQuantity

    
//     totalQuantityStr := strconv.Itoa(totalQuantity)

//     _, err = s.db.Exec(`
//         UPDATE medications
//         SET quantity = $1
//         WHERE id = $2
//     `, totalQuantityStr, id)

//     return err
// }


func (s *Store) GetLogsByCenterID(centerID int) ([]types.MedicationLog, error) {
	rows, err := s.db.Query(`
        SELECT id, name_arabic, dosage, medication_type, requested_quantity, requested_at
        FROM medication_requests
        WHERE center_id = $1
        ORDER BY requested_at DESC
    `, centerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []types.MedicationLog
	for rows.Next() {
		var log types.MedicationLog
		err := rows.Scan(
			&log.ID,
			&log.NameArabic,
			&log.Dosage,
			&log.MedicationType,
			&log.Quantity,
			&log.RequestedAt,
		
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}





func (s *Store) GetReviewMedicationNames(centerID int) ([]types.GeTMedicationReview, error) {
	rows, err := s.db.Query(`
		SELECT  name_arabic , dosage , units_per_box
		FROM medications
		WHERE center_id = $1
	`, centerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewMedicine []types.GeTMedicationReview
	for rows.Next() {
		var rm types.GeTMedicationReview
		if err := rows.Scan(&rm.NameArabic , &rm.Dosage , &rm.UnitsPerBox); err != nil {
			return nil, err
		}
		reviewMedicine = append(reviewMedicine, rm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviewMedicine, nil
}














func (s *Store) GetMedicationByID(id int) (*types.GeTMedication, error) {
    var m types.GeTMedication

    err := s.db.QueryRow(`
        SELECT id, name_arabic, name_english, medication_type, dosage, quantity, units_per_box , center_id
        FROM medications
        WHERE id = $1
    `, id).Scan(
        &m.ID,
        &m.NameArabic,
        &m.NameEnglish,
        &m.MedicationType,
        &m.Dosage,
        &m.Quantity,
        &m.UnitsPerBox,
		&m.CenterID,
    )

    if err != nil {
        return nil, err
    }

    return &m, nil
}









func (s *Store) UpdateMedicationQuantity(id int, decreaseQuantity string) error {
    var oldQuantityStr string
    err := s.db.QueryRow(`
        SELECT quantity FROM medications WHERE id = $1
    `, id).Scan(&oldQuantityStr)
    if err != nil {
        return err
    }

    oldQuantity, err := strconv.Atoi(oldQuantityStr)
    if err != nil {
        oldQuantity = 0
    }
    dQuan, err := strconv.Atoi(decreaseQuantity)
    if err != nil {
        oldQuantity = 0
    }

    
    totalQuantity := oldQuantity - dQuan

  
    // if totalQuantity < 0 {
    //     totalQuantity = 0
    // }

    totalQuantityStr := strconv.Itoa(totalQuantity)

    _, err = s.db.Exec(`
        UPDATE medications
        SET quantity = $1
        WHERE id = $2
    `, totalQuantityStr, id)

    return err
}
