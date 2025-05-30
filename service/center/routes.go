package center

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/AL-Hourani/care-center/types"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/golang-jwt/jwt/v5"

	// "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.CenterStore
    pStore types.PatientStore
}

func NewHandler(store types.CenterStore , patientStore types.PatientStore) *Handler {
	return &Handler {
		store: store,
		pStore: patientStore,
	}
}

func (h *Handler) RegisterCenterRoutes(router *mux.Router) {
	// router.HandleFunc("/centerLogin", h.handleCenterLogin).Methods("POST")
	router.HandleFunc("/centerRegister", h.handleCenterRegister).Methods("POST")
	router.HandleFunc("/confirmAccount", h.handleConfirmPatientAccount).Methods("POST")
	router.HandleFunc("/getCenters/{city}",h.handleGetCenters).Methods("GET")
	router.HandleFunc("/getCities",h.handleGetCities).Methods("GET")
	router.HandleFunc("/getCenterProfile/{id}",h.handleGetCenetrProfile).Methods("GET")
	router.HandleFunc("/updateCenterProfile", h.handleUpdateCenterProfile).Methods(http.MethodPatch)
	router.HandleFunc("/getPatients", auth.WithJWTAuth(h.handleGetPatients)).Methods(http.MethodGet)
	router.HandleFunc("/addPatient/{id}", h.handleGetCenters).Methods(http.MethodPost)
	router.HandleFunc("/updatePatient", h.handleUpdatePatient).Methods(http.MethodPatch)
	router.HandleFunc("/deletePatient/{id}", auth.WithJWTAuth(h.handleDeletePatient)).Methods(http.MethodDelete)
	router.HandleFunc("/logout",auth.WithJWTAuth(h.Logout)).Methods("POST")
	router.HandleFunc("/deleteCenter",h.handleDeleteCenter).Methods(http.MethodDelete)
	router.HandleFunc("/addReviewe",h.handleAddReviewe).Methods("POST")
	router.HandleFunc("/reviewdelete/{id}", h.handleDeleteReview).Methods("DELETE")


}




func (h *Handler) handleCenterRegister(w http.ResponseWriter , r *http.Request) {
	//get json payload
		var centerPayload types.RegisterCenterPayload
		if err := utils.ParseJSON(r , &centerPayload); err != nil {
			utils.WriteError(w , http.StatusBadRequest , err)
		}

	// check from secret center key

	   if centerPayload.CenterKey != config.GetEnv("CENTER_KEY" , ""){
		   utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid center Key"))
		   return
	   }

		
	// check if center exists
        // _ , err := h.store.GetCenterByName(centerPayload.CenterName)
		// if err == nil {
		// 	utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("center with name %s already exists" , centerPayload.CenterName))
		// 	return 
		// }
	
	//check if the center email is uniqe 
	 _ , err := h.store.GetCenterByEmail(centerPayload.CenterEmail)
	if err == nil {
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("center with email %s already exists" , centerPayload.CenterEmail))
		return 
	}
	

		hashedPassword , err := auth.HashPassword(centerPayload.CenterPassword)
		if err != nil {
			utils.WriteError(w , http.StatusInternalServerError , err)
		}

		
	//                verfiy email 
	
	//1- send otp code to the user


	//2- check if user otp corccert 




	

	//if it dosen't we create the new center
	err = h.store.GreateCenter(types.Center{
		CenterName: centerPayload.CenterName,
		CenterPassword:hashedPassword,
		CenterEmail:centerPayload.CenterEmail,
        CenterCity: centerPayload.CenterCity,
	})
  


	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}

	newLoginFailed := types.InsertLogin {
		Email:centerPayload.CenterEmail ,
		Password:hashedPassword ,
	}
	
	err = h.store.GreateLoginFailed(newLoginFailed)
	

	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}

	utils.WriteJSON(w , http.StatusCreated ,  map[string]string{"message":"successfully Created"})
}

func (h *Handler) handleGetPatients(w http.ResponseWriter , r *http.Request) {

	token, ok := r.Context().Value(auth.UserContextKey).(*jwt.Token)
	if !ok {
		http.Error(w, "Unauthorized: No token found", http.StatusUnauthorized)
		return
	}

	id, err := auth.GetIDFromToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	patientsList , err := h.store.GetPatientsForCenter(id)
	if err != nil {
		utils.WriteError(w , http.StatusInternalServerError , err)
		return
	}

	utils.WriteJSON(w , http.StatusOK , patientsList)
}



func (h *Handler) handleGetCenters(w http.ResponseWriter , r *http.Request) {
	vars := mux.Vars(r)
	cityName, ok := vars["city"]
	if !ok || cityName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("city name is required"))
		return
	}

	centerList , err := h.store.GetCentersByCity(cityName)
	if err != nil {
		utils.WriteError(w , http.StatusNotFound , err)
		return
	}
	utils.WriteJSON(w , http.StatusOK , centerList)
}



func (h *Handler) handleDeletePatient(w http.ResponseWriter , r *http.Request) {
	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}

	err = h.store.DeletePatient(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest , err)
	}

	utils.WriteJSON(w , http.StatusOK ,  map[string]string{"message":"successfully Deleted"})
}


// handle confirm patient ....


func (h *Handler) handleConfirmPatientAccount(w http.ResponseWriter , r *http.Request) {
	var confrimAccout types.ConfirmAccount
	
	if err := utils.ParseJSON(r , &confrimAccout); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
	}

	err := h.store.UpdateIsCompletedPatientField(confrimAccout)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,err)
	}

	utils.WriteJSON(w , http.StatusOK ,  map[string]string{"message":"successfully Confirm Account"})

}






// update patient

func (h *Handler) handleUpdatePatient(w http.ResponseWriter , r *http.Request) { 
	var udpatePayload types.PatientUpdatePayload
		
	if err := utils.ParseJSON(r , &udpatePayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

	if udpatePayload.ID == 0 {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("patient ID is required"))
        return
    }

	
    // تحديث بيانات المريض جزئيًا باستخدام دالة PATCH
    err := h.store.PatchUpdatePatient(&udpatePayload)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError,fmt.Errorf("faild updated"))
        return
    }

	patient , err := h.pStore.GetPatientDetailsByID(udpatePayload.ID)
	if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }

    utils.WriteJSON(w, http.StatusOK,patient)
}





//logout

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}



func (h *Handler) handleGetCities(w http.ResponseWriter, r *http.Request) {

   cities_list , err := h.store.GetCities()
   if err != nil {
	utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error get cities"))
	return
   }

   utils.WriteJSON(w, http.StatusOK, cities_list)
}



// get center profile 

func (h *Handler) handleGetCenetrProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}

	centerProfile ,err := h.store.GetCenterProfile(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error get center profile"))
	}

	utils.WriteJSON(w, http.StatusOK, centerProfile)

}


func (h *Handler) handleDeleteCenter(w http.ResponseWriter, r *http.Request) {
	var deleteCenter types.DeleteCenter
	if err := utils.ParseJSON(r , &deleteCenter); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

	getCenter , err := h.store.GetCenterByName(deleteCenter.CenterName)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("no center for this name"))
	}
	toCenter , err := h.store.GetCenterByName(deleteCenter.CenterNameReassignPatients)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("no center for this name"))
	}

	numberPatinets , err := h.store.GetPatientCountByCenterName(getCenter.CenterName)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error get number of center"))
	}
	if numberPatinets > 0 {
		err := h.store.DeleteCenterAndReassignPatients(getCenter.ID, toCenter.ID) 
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error delete Center And Reassign Patients"))
		}
	}

	err = h.store.DeleteCenter(getCenter.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error delete Center"))
	}

	utils.WriteJSON(w , http.StatusOK , map[string]string{
		"message": "delete and And Reassign Patients successfully",
	})

}





//update center profile 
func (h *Handler) handleUpdateCenterProfile(w http.ResponseWriter, r *http.Request) {
	var centerUpdatePayload types.CenterUpdateProfilePayload
	if err := utils.ParseJSON(r , &centerUpdatePayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

	err := h.store.CenterUpdateCenterProfile(centerUpdatePayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error update center profile"))
		return
	}

	updateProfileInfo , err := h.store.GetCenterUpdateCenterProfile(centerUpdatePayload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error get center profile"))
		return
	}

	utils.WriteJSON(w , http.StatusOK , updateProfileInfo)
}






// 




























// handle with reviws
func (h *Handler) handleAddReviewe (w http.ResponseWriter, r *http.Request) { 
    var  AddReviewsPayload types.AddReviwePayload


	if err := utils.ParseJSON(r , &AddReviewsPayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}



	newReviewe := types.Reviwe{
		PatientID: AddReviewsPayload.PatientID,
		Address: AddReviewsPayload.Address,
		Weight: AddReviewsPayload.Weight,
		LengthPatient: AddReviewsPayload.LengthPatient,
		SugarType: AddReviewsPayload.SugarType,
		OtherDisease: AddReviewsPayload.OtherDisease,
		HistoryOfFamilyDisease: AddReviewsPayload.HistoryOfFamilyDisease,
		HistoryOfDiseaseDetection: AddReviewsPayload.HistoryOfDiseaseDetection,
		Gender: AddReviewsPayload.Gender,
		Hemoglobin: AddReviewsPayload.Hemoglobin,
		Grease: AddReviewsPayload.Grease,
		UrineAcid: AddReviewsPayload.UrineAcid,
		BloodPressure: AddReviewsPayload.BloodPressure,
		Cholesterol: AddReviewsPayload.Cholesterol,
		LDL: AddReviewsPayload.LDL,
		HDL: AddReviewsPayload.HDL,
		Creatine: AddReviewsPayload.Creatine,
		Normal_Glocose: AddReviewsPayload.Normal_Glocose,
		Glocose_after_Meal: AddReviewsPayload.Glocose_after_Meal,
		Triple_Grease: AddReviewsPayload.Triple_Grease,
		Hba1c: AddReviewsPayload.Hba1c,
		Coments: AddReviewsPayload.Coments,
		
	}

	// add review 
	Review_id , err := h.store.InsertReview(newReviewe)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error in adding review info : %v" , err))
		return
	}


	// get id for this review


	// add treatment
	newTreatment := types.TreatmentInsert {
	    ReviewID: Review_id,
		Speed: AddReviewsPayload.Treatments.Speed,
		Type: AddReviewsPayload.Treatments.Type,

	}
	treatmentID ,err := h.store.InsertTreatment(newTreatment)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error in adding Tratment"))
		return
	}


	for _, drug := range AddReviewsPayload.Treatments.Drugs {
	drugID, err := h.store.FindOrCreateDrugByName(drug.Name)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add/find drug: %w", err))
		return
	}

	// ربط الدواء بالعلاج
	err = h.store.InsertTreatmentDrug(types.TreatmentDrug{
		TreatmentID:   treatmentID,
		DrugID:        drugID,
		DosagePerDay:  drug.Dosage_per_day,
		Units:         drug.Units,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to insert treatment-drug: %w", err))
		return
	}
}




	if AddReviewsPayload.Has_a_eye_disease {

      newClininEyeInfo := types.Clinic_Eye {
		ReviewID: Review_id,
        Has_a_eye_disease: AddReviewsPayload.Has_a_eye_disease,
		In_kind_disease: AddReviewsPayload.In_kind_disease,
		Relationship_eyes_with_diabetes: AddReviewsPayload.Relationship_eyes_with_diabetes,
		Comments_eyes_clinic: AddReviewsPayload.Comments_eyes_clinic,
	  }

	  	err := h.store.InsertClinicEye(newClininEyeInfo)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error in adding clinic eye info"))
				return
			}
	}

	// 2

	if AddReviewsPayload.Has_a_heart_disease {
		
      newClininHeartInfo := types.Clinic_heart {
		ReviewID: Review_id,
		Has_a_heart_disease: AddReviewsPayload.Has_a_heart_disease,
		Heart_disease: AddReviewsPayload.Heart_disease,
		Relationship_heart_with_diabetes: AddReviewsPayload.Relationship_heart_with_diabetes,
		Comments_heart_clinic: AddReviewsPayload.Comments_heart_clinic,
	  }
	  	  	err := h.store.InsertClinicHeart(newClininHeartInfo)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error in adding clinic heart info"))
				return
			}

	}



	// 3

	if AddReviewsPayload.Has_a_nerve_disease {
				
      newClininNerveInfo := types.Clinic__nerve {
		ReviewID: Review_id,
		Has_a_nerve_disease: AddReviewsPayload.Has_a_nerve_disease,
		Nerve_disease: AddReviewsPayload.Nerve_disease,
		Relationship_nerve_with_diabetes: AddReviewsPayload.Relationship_nerve_with_diabetes,
		Comments_nerve_clinic: AddReviewsPayload.Comments_nerve_clinic,
	  }
	  	  	err := h.store.InsertClinicNerve(newClininNerveInfo)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error in adding clinic nerve info"))
				return
			}

	}


	// 4
	
	if AddReviewsPayload.Has_a_bone_disease {
				
      newClininBoneInfo := types.Clinic__bone {
		ReviewID: Review_id,
	    Has_a_bone_disease: AddReviewsPayload.Has_a_bone_disease,
		Bone_disease: AddReviewsPayload.Bone_disease,
		Relationship_bone_with_diabetes: AddReviewsPayload.Relationship_bone_with_diabetes,
		Comments_bone_clinic: AddReviewsPayload.Comments_bone_clinic,
	  }
	  	  	err := h.store.InsertClinicBone(newClininBoneInfo)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error in adding clinic bone info"))
				return
			}

	}

	// 5

		
	if AddReviewsPayload.Has_a_urinary_disease {
				
      newClininUrinaryInfo := types.Clinic__urinary {
		ReviewID: Review_id,
		Has_a_urinary_disease: AddReviewsPayload.Has_a_urinary_disease,
		Urinary_disease: AddReviewsPayload.Urinary_disease,
		Relationship_urinary_with_diabetes: AddReviewsPayload.Relationship_urinary_with_diabetes,
		Comments_urinary_clinic: AddReviewsPayload.Comments_urinary_clinic,
	  }
	  	  	err := h.store.InsertClinicUrinary(newClininUrinaryInfo)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error in adding clinic urinary info"))
				return
			}

	}




	utils.WriteJSON(w , http.StatusOK , map[string]string{
		"message": "Review added successfully",
	})



}










func (h *Handler) handleDeleteReview(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("review ID is required"))
		return
	}

	// تحويل ID من نص إلى عدد صحيح
	reviewID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid review ID"))
		return
	}

	// تنفيذ عملية الحذف
	err = h.store.DeleteReviewByID(reviewID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete review: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Review deleted successfully",
	})
}
