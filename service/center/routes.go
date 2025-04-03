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
	router.HandleFunc("/getPatients", auth.WithJWTAuth(h.handleGetPatients)).Methods(http.MethodGet)
	router.HandleFunc("/addPatient/{id}", h.handleGetCenters).Methods(http.MethodPost)
	router.HandleFunc("/updatePatient", h.handleUpdatePatient).Methods(http.MethodPatch)
	router.HandleFunc("/deletePatient/{id}", auth.WithJWTAuth(h.handleDeletePatient)).Methods(http.MethodDelete)
	router.HandleFunc("/logout",auth.WithJWTAuth(h.Logout)).Methods("POST")
	router.HandleFunc("/deleteCenter",h.handleDeleteCenter).Methods("POST")
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