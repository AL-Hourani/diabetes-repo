package supervisor

import (
	"database/sql"
	"fmt"
	"strconv"
    "time"
    "github.com/xuri/excelize/v2"
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
        WHERE request_id = $2
    `, newStatus, id)

    return err
}




func (s *Store) UpdateMedicationQuantity(id int, newQuantity int) error {
   
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

    // نجمع القديم مع الجديد
    totalQuantity := oldQuantity + newQuantity

    // نحفظه كـ string
    totalQuantityStr := strconv.Itoa(totalQuantity)

    // تحديث الجدول
    _, err = s.db.Exec(`
        UPDATE medications
        SET quantity = $1
        WHERE id = $2
    `, totalQuantityStr, id)

    return err
}






// للجلب اسم المراكز بناء هلى اسم المينة 



func (s *Store) GetCentersByCity(cityName string) ([]string, error) {
    rows, err := s.db.Query(`
        SELECT centerName 
        FROM centers
        WHERE centerCity ILIKE $1
    `, cityName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var centers []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, err
        }
        centers = append(centers, name)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return centers, nil
}



func (s *Store) GetPatientCountByCity(cityName string) (int, error) {
    var count int
    err := s.db.QueryRow(`
        SELECT COUNT(p.id)
        FROM patients p
        INNER JOIN centers c ON p.center_id = c.id
        WHERE c.centerCity ILIKE $1
    `, cityName).Scan(&count)

    if err != nil {
        return 0, err
    }
    return count, nil
}


func (s *Store) GetPatientCountByCityLastMonth(cityName string) (int, error) {
    var count int
    err := s.db.QueryRow(`
        SELECT COUNT(p.id)
        FROM patients p
        INNER JOIN centers c ON p.center_id = c.id
        WHERE c.centerCity ILIKE $1
          AND p.createAt >= NOW() - INTERVAL '1 month'
    `, cityName).Scan(&count)

    if err != nil {
        return 0, err
    }
    return count, nil
}



func (s *Store) GetCenterWithMostPatients() (*types.CenterWithCount, error) {
	query := `
		SELECT c.id, c.centerName, c.centerCity, COUNT(p.id) AS patients_count
		FROM centers c
		LEFT JOIN patients p ON c.id = p.center_id
		GROUP BY c.id, c.centerName, c.centerCity
		ORDER BY patients_count DESC
		LIMIT 1;
	`

	var result types.CenterWithCount

	err := s.db.QueryRow(query).Scan(
		&result.ID,
		&result.CenterName,
		&result.CenterCity,
		&result.PatientsCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no centers found")
		}
		return nil, err
	}

	return &result, nil
}





















func (s *Store) GetPatientReviewsByMonth(month, year int) ([]types.PatientReview, error) {
    query := `
        SELECT 
            r.id AS review_id,
            r.patient_id,
            p.fullName,
            p.email,
            p.phone,
            pm.gender,
            pm.sugarType,
            r.address_patient,
            r.wight,
            r.length_patient,
            r.otherDisease,
            r.hemoglobin,
            r.grease,
            r.urineAcid,
            r.bloodPressure,
            r.cholesterol,
            r.LDL,
            r.HDL,
            r.creatine,
            r.normal_clucose,
            r.clucose_after_meal,
            r.triple_grease,
            r.hba1c,
            COALESCE(NULLIF(e.comments, ''), 'لا يوجد') AS comments,
            r.date_review,

            -- بيانات العيون
            COALESCE(
                CASE WHEN e.has_a_eye_disease THEN 'يوجد مرض' ELSE 'لا يوجد مرض' END, 'لا يوجد'
            ) AS has_a_eye_disease,
            COALESCE(NULLIF(e.in_kind_disease, ''), 'لا يوجد') AS in_kind_disease,
            COALESCE(
                CASE WHEN e.relationship_with_diabetes THEN 'نعم' ELSE 'لا' END, 'لا يوجد'
            ) AS relationship_with_diabetes,
            COALESCE(NULLIF(e.comments, ''), 'لا يوجد') AS comments_eye,

            -- بيانات القلب
            COALESCE(
                CASE WHEN h.has_a_heart_disease THEN 'يوجد مرض' ELSE 'لا يوجد مرض' END, 'لا يوجد'
            ) AS has_a_heart_disease,
            COALESCE(NULLIF(h.heart_disease, ''), 'لا يوجد') AS heart_disease,

            COALESCE(
                CASE WHEN h.relationship_with_diabetes THEN 'نعم' ELSE 'لا' END, 'لا يوجد'
            ) AS relationship_heart_with_diabetes,

            COALESCE(NULLIF(h.comments, ''), 'لا يوجد') AS comments_heart,

            -- بيانات الأعصاب
            COALESCE(
                CASE WHEN n.has_a_nerve_disease THEN 'يوجد مرض' ELSE 'لا يوجد مرض' END, 'لا يوجد'
            ) AS has_a_nerve_disease,
            COALESCE(NULLIF(n.nervous_disease,'') ,'لا يوجد') AS nervous_disease,
            COALESCE(
                CASE WHEN n.relationship_with_diabetes THEN 'نعم' ELSE 'لا' END, 'لا يوجد'
            ) AS relationship_nervous_with_diabetes,
            COALESCE(NULLIF(n.comments,''), 'لا يوجد') AS comments_nervous,

            -- بيانات العظام
            COALESCE(
                CASE WHEN b.has_a_bone_disease THEN 'يوجد مرض' ELSE 'لا يوجد مرض' END, 'لا يوجد'
            ) AS has_a_bone_disease,
            COALESCE(NULLIF(b.nervous_disease,''), 'لا يوجد') AS bone_disease,
            COALESCE(
                CASE WHEN b.relationship_with_diabetes THEN 'نعم' ELSE 'لا' END, 'لا يوجد'
            ) AS relationship_bone_with_diabetes,
            COALESCE(NULLIF(b.comments,'') ,'لا يوجد') AS comments_bone,

            -- بيانات الجهاز البولي
            COALESCE(
                CASE WHEN u.has_a_urinary_disease THEN 'يوجد مرض' ELSE 'لا يوجد مرض' END, 'لا يوجد'
            ) AS has_a_urinary_disease,
            COALESCE(NULLIF(u.nervous_disease,''), 'لا يوجد') AS urinary_disease,
            COALESCE(
                CASE WHEN u.relationship_with_diabetes THEN 'نعم' ELSE 'لا' END, 'لا يوجد'
            ) AS relationship_urinary_with_diabetes,
            COALESCE(NULLIF(u.comments,''), 'لا يوجد') AS comments_urinary,

            COALESCE(t.treatment_type::text, '[]') AS treatment_type

        FROM reviews r
        LEFT JOIN patients p ON p.id = r.patient_id
        LEFT JOIN patient_m pm ON pm.patient_id = r.patient_id
        LEFT JOIN eyes_clinic e ON e.review_id = r.id
        LEFT JOIN heart_clinic h ON h.review_id = r.id
        LEFT JOIN nerve_clinic n ON n.review_id = r.id
        LEFT JOIN bone_clinic b ON b.review_id = r.id
        LEFT JOIN urinary_clinic u ON u.review_id = r.id
        LEFT JOIN treatments t ON t.review_id = r.id
        LEFT JOIN treatment_drugs td ON td.treatment_id = t.id
        WHERE EXTRACT(MONTH FROM r.date_review) = $1
          AND EXTRACT(YEAR FROM r.date_review) = $2
          GROUP BY 
        r.id, r.patient_id, p.fullName, p.email, p.phone, pm.gender, pm.sugarType,
        r.address_patient, r.wight, r.length_patient, r.otherDisease, r.hemoglobin,
        r.grease, r.urineAcid, r.bloodPressure, r.cholesterol, r.LDL, r.HDL, r.creatine,
        r.normal_clucose, r.clucose_after_meal, r.triple_grease, r.hba1c, r.comments, r.date_review,
        e.has_a_eye_disease, e.in_kind_disease, e.relationship_with_diabetes, e.comments,
        h.has_a_heart_disease, h.heart_disease, h.relationship_with_diabetes, h.comments,
        n.has_a_nerve_disease, n.nervous_disease, n.relationship_with_diabetes, n.comments,
        b.has_a_bone_disease, b.bone_disease, b.relationship_with_diabetes, b.comments,
        u.has_a_urinary_disease, u.urinary_disease, u.relationship_with_diabetes, u.comments,
        t.treatment_type
        ORDER BY r.date_review;
    `

    rows, err := s.db.Query(query, month, year)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reviews []types.PatientReview
    for rows.Next() {
        var r types.PatientReview
        err := rows.Scan(
            &r.ReviewID , &r.PatientID,
            &r.PatientFullName, &r.PatientEmail, &r.PatientPhone,
            &r.AddressPatient, &r.Wight, &r.LengthPatient,
            &r.OtherDisease, &r.Hemoglobin, &r.Grease, &r.UrineAcid, &r.BloodPressure,
            &r.Cholesterol, &r.LDL, &r.HDL, &r.Creatine, &r.NormalClucose,
            &r.ClucoseAfterMeal, &r.TripleGrease, &r.Hba1c,&r.Comments, &r.DateReview,

            &r.Has_a_eye_disease, &r.In_kind_disease, &r.Relationship_with_diabetes, &r.Comments_eye,
            &r.Has_a_heart_disease, &r.Heart_disease, &r.Relationship_heart_with_diabetes, &r.Comments_heart,
            &r.Has_a_nerve_disease, &r.Nervous_disease, &r.Relationship_nervous_with_diabetes, &r.Comments_nervous,
            &r.Has_a_bone_disease, &r.Bone_disease, &r.Relationship_bone_with_diabetes, &r.Comments_bone,
            &r.Has_a_urinary_disease, &r.Urinary_disease, &r.Relationship_urinary_with_diabetes, &r.Comments_urinary,
        )
        if err != nil {
            return nil, err
        }
        var treatmentDrugs []types.TreatmentDrugExel

        drugRows, err := s.db.Query(`
            SELECT m.name_arabic, td.dosage_per_day, td.quantity
            FROM treatment_drugs td
            JOIN medications m ON m.id = td.drug_id
            WHERE td.treatment_id = (
                SELECT id FROM treatments WHERE review_id = $1
            )
        `, r.ReviewID)
        if err != nil {
            return nil, err
        }


        for drugRows.Next() {
            var d types.TreatmentDrugExel
            if err := drugRows.Scan(&d.DrugName, &d.DosagePerDay, &d.Quantity); err != nil {
                return nil, err
            }
            treatmentDrugs = append(treatmentDrugs, d)
        }
        defer drugRows.Close()
        r.TreatmentDrugs = treatmentDrugs
        reviews = append(reviews, r)
    }

           
    return reviews, nil
}


func CreateExcelFile(reviews []types.PatientReview) (*excelize.File, error) {
    f := excelize.NewFile()
    sheet := "Patients"
    f.NewSheet(sheet)

    headers := []string{
        "اسم المريض", "البريد الإلكتروني", "الهاتف",
        "الجنس", "نوع السكر",
        "العنوان", "الوزن", "الطول", "أمراض أخرى", "الهيموغلوبين",
        "الدهون", "حمض اليوريك", "ضغط الدم", "الكوليسترول", "LDL", "HDL",
        "الكرياتين", "سكر طبيعي", "سكر بعد الوجبة", "الدهون الثلاثية", "HBA1c", "ملاحظات", "تاريخ المراجعة",
        "أمراض العيون", "نوع المرض بالعيون", "علاقة العين بالسكري", "ملاحظات العيون",
        "أمراض القلب", "نوع المرض بالقلب", "علاقة القلب بالسكري", "ملاحظات القلب",
        "أمراض الأعصاب", "نوع المرض بالأعصاب", "علاقة الأعصاب بالسكري", "ملاحظات الأعصاب",
        "أمراض العظام", "نوع المرض بالعظام", "علاقة العظام بالسكري", "ملاحظات العظام",
        "أمراض الجهاز البولي", "نوع المرض بالجهاز البولي", "علاقة الجهاز البولي بالسكري", "ملاحظات الجهاز البولي",
        "نوع العلاج", "الأدوية (جرعة/يوم و كمية)",
    }

    // تنسيق الهيدر
    headerStyle, _ := f.NewStyle(&excelize.Style{
        Font: &excelize.Font{Bold: true, Family: "Tahoma", Size: 12},
        Fill: excelize.Fill{Type: "pattern", Color: []string{"#D9E1F2"}, Pattern: 1},
        Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
    })

    for i, h := range headers {
        col, _ := excelize.ColumnNumberToName(i + 1)
        cell := col + "1"
        f.SetCellValue(sheet, cell, h)
        f.SetCellStyle(sheet, cell, cell, headerStyle)
    }

    // تعبئة البيانات
    for i, r := range reviews {
        row := i + 2

        // دمج بيانات الأدوية في عمود واحد كنص
        treatmentDrugs := ""
        for _, d := range r.TreatmentDrugs {
            treatmentDrugs += fmt.Sprintf("%s: %s جرعة/يوم، الكمية: %d\n", d.DrugName, d.DosagePerDay, d.Quantity)
        }

        values := []interface{}{
            r.PatientFullName, r.PatientEmail, r.PatientPhone,
            r.Gender, r.SugarType,
            r.AddressPatient, r.Wight, r.LengthPatient, r.OtherDisease, r.Hemoglobin,
            r.Grease, r.UrineAcid, r.BloodPressure, r.Cholesterol, r.LDL, r.HDL,
            r.Creatine, r.NormalClucose, r.ClucoseAfterMeal, r.TripleGrease,
            r.Hba1c, r.Comments, r.DateReview,
            r.Has_a_eye_disease, r.In_kind_disease, r.Relationship_with_diabetes, r.Comments_eye,
            r.Has_a_heart_disease, r.Heart_disease, r.Relationship_heart_with_diabetes, r.Comments_heart,
            r.Has_a_nerve_disease, r.Nervous_disease, r.Relationship_nervous_with_diabetes, r.Comments_nervous,
            r.Has_a_bone_disease, r.Bone_disease, r.Relationship_bone_with_diabetes, r.Comments_bone,
            r.Has_a_urinary_disease, r.Urinary_disease, r.Relationship_urinary_with_diabetes, r.Comments_urinary,
            r.TreatmentType,
            treatmentDrugs,
        }

        for j, v := range values {
            col, _ := excelize.ColumnNumberToName(j + 1)
            f.SetCellValue(sheet, col+strconv.Itoa(row), v)
        }
    }

    // ضبط عرض الأعمدة تلقائياً
    f.SetColWidth(sheet, "A", "AQ", 25)

    return f, nil
}





func (s *Store) ParseMonthYear(input string) (month int, year int, err error) {
 
    t, err := time.Parse("January 2006", input)
    if err != nil {
        return 0, 0, err
    }
    month = int(t.Month())
    year = t.Year()
    return month, year, nil
}