package patient

import (
	"fmt"
	"net/http"

	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/service/auth"
	// "github.com/AL-Hourani/care-center/service/patient"
	"github.com/AL-Hourani/care-center/types"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.PatientStore
	storeCenter types.CenterStore
}

func NewHandler(store types.PatientStore , centerStore types.CenterStore ) *Handler {
	return &Handler{store: store , storeCenter: centerStore }
}

func (h *Handler) RegisterPatientRoutes(router *mux.Router) {
	router.HandleFunc("/Login", h.handleLogin).Methods("POST")
	router.HandleFunc("/patientRegister", h.handlePatientRegister).Methods("POST")
}




// handle login for both centres or patients
//----------------------------------------------------------------------------



func (h *Handler) handleLogin(w http.ResponseWriter , r *http.Request) {
	//get json payload
	var LoginPayload types.LoginPayload
	if err := utils.ParseJSON(r , &LoginPayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

	//validate the payoad .....................
	if err := utils.Validate.Struct(LoginPayload);err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid payload %v", error) )
		return
	}




	// find patient .............................................................



    patient , errLogin := h.store.GetPatientByEmail(LoginPayload.Email)
	if errLogin == nil {
	    if !auth.ComparePasswords(patient.Password , [] byte(LoginPayload.Password)) {
				utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("not found invalid password"))
			return
	    }

		secret := []byte(config.Envs.JWTSecret)
		token , err := auth.CreateJWT(secret , patient.ID)
		
		if err != nil {
			utils.WriteError(w , http.StatusInternalServerError ,err)
			return
		}
	
		utils.WriteJSON(w , http.StatusOK ,map[string]string{"toke":token})
		
	} else  {
		center , err2 := h.storeCenter.GetCenterByEmail(LoginPayload.Email)
		if err2 != nil {
			utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("not found invalid  email or password"))
			return 
		}
		if !auth.ComparePasswords(center.CenterPassword , [] byte(LoginPayload.Password)) {
			utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("not found invalid password"))
			return
		}
	
		secret := []byte(config.Envs.JWTSecret)
		token , err := auth.CreateJWT(secret , center.ID)
		
		if err != nil {
			utils.WriteError(w , http.StatusInternalServerError ,err)
			return
		}
	
		utils.WriteJSON(w , http.StatusOK ,map[string]string{"toke":token})
	}
	

}




func (h *Handler) handlePatientRegister(w http.ResponseWriter , r *http.Request) {

	//get json payload
		var patientPayload types.RegisterPatientPayload
		if err := utils.ParseJSON(r , &patientPayload); err != nil {
			utils.WriteError(w , http.StatusBadRequest , err)
			return
		}

    //validate the payoad .....................
	if err := utils.Validate.Struct(patientPayload);err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid payload %v", error) )
		return
	}


	// check if user exists

	_ , err := h.store.GetPatientByEmail(patientPayload.Email)
	if err == nil {
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("patient with email %s already exists" , patientPayload.Email))
		return 
	}

	hashedPassword , err := auth.HashPassword(patientPayload.Password)
	if err != nil {
		utils.WriteError(w , http.StatusInternalServerError , err)
	}

	//get center 

	cenetr , err := h.storeCenter.GetCenterByName(patientPayload.CenterName)
	if err != nil {
		utils.WriteError(w , http.StatusInternalServerError , err)
	}
	//if it dosen't we create the new user
	err = h.store.GreatePatient(types.Patient{
		FullName: patientPayload.FullName,
		Email: patientPayload.Email,
		Password: hashedPassword,
		Age: patientPayload.Age,
		Phone:patientPayload.Phone,
		CenterID: cenetr.ID,
		
	})

	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}

	utils.WriteJSON(w , http.StatusCreated , map[string]string{"message":"successfully Created"})
}

