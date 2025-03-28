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
}

func NewHandler(store types.CenterStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterCenterRoutes(router *mux.Router) {
	// router.HandleFunc("/centerLogin", h.handleCenterLogin).Methods("POST")
	router.HandleFunc("/centerRegister", h.handleCenterRegister).Methods("POST")
	router.HandleFunc("/confirmAccount", h.handleConfirmPatientAccount).Methods("POST")
	router.HandleFunc("/getCenters",h.handleGetCenters).Methods("GET")
	router.HandleFunc("/getPatients", auth.WithJWTAuth(h.handleGetPatients)).Methods(http.MethodGet)
	router.HandleFunc("/addPatient/{id}", h.handleGetCenters).Methods(http.MethodPost)
	router.HandleFunc("/updatePatient", h.handleUpdatePatient).Methods(http.MethodPatch)
	router.HandleFunc("/deletePatient/{id}", auth.WithJWTAuth(h.handleDeletePatient)).Methods(http.MethodDelete)
	router.HandleFunc("/logout",auth.WithJWTAuth(h.Logout)).Methods("POST")
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

	//if it dosen't we create the new center
	err = h.store.GreateCenter(types.Center{
		CenterName: centerPayload.CenterName,
		CenterPassword:hashedPassword,
		CenterEmail:centerPayload.CenterEmail,

	})
  
	

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
	centerList , err := h.store.GetCenters()
	if err != nil {
		utils.WriteError(w , http.StatusInternalServerError , err)
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
		
	if err := utils.ParseJSONUpdate(r , &udpatePayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
	}

	if udpatePayload.ID == 0 {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("patient ID is required"))
        return
    }

	
    // تحديث بيانات المريض جزئيًا باستخدام دالة PATCH
    err := h.store.PatchUpdatePatient(&udpatePayload)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }

    utils.WriteJSON(w, http.StatusOK,  map[string]string{"message":"Patient updated successfully"})
}





//logout

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}
