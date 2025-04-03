package patient

import (
	"fmt"
	"net/http"
	"strconv"
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
	router.HandleFunc("/getPatient/{id}" , h.handleGetPatient).Methods("GET")
	router.HandleFunc("/getPatientProfile/{id}" , h.handleGetPatientProfile).Methods("GET")
	router.HandleFunc("/getAllPatientInfo/{id}" , h.handleGetAllPatientInfo).Methods("GET")
	router.HandleFunc("/verify-token", h.VerifyTokenHandler).Methods("POST")
	router.HandleFunc("/verifyOtp", h.VerifyOTPHandler).Methods("POST")
	router.HandleFunc("/updatePatientProfile", h.handleUpdatePatientProfile).Methods(http.MethodPatch)

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
	user, err := h.store.GetUserByEmail(LoginPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

		// التحقق من صحة كلمة المرور
		if !auth.ComparePasswords(user.Password, []byte(LoginPayload.Password)) {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
			return
		}

		// إنشاء JWT Token
		secret := []byte(config.Envs.JWTSecret)
		token, err := auth.CreateJWT(secret, user.ID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if user.Role == "patient" {
			returnLoggingData := types.ReturnLoggingData{
				Name:       user.Name,
				Email:      user.Email,
				Role:       "patient",
				IsCompletes: false,
				Token:      token,
			}
			utils.WriteJSON(w, http.StatusOK, returnLoggingData)
		} else {

	
			returnLoggingData := types.ReturnLoggingCenterData{
				Name:    user.Name,
				Email:   user.Email,
				Role:    "center",
				Token:   token,
			}
			utils.WriteJSON(w, http.StatusOK, returnLoggingData)
		}
	

}



var pendingPatients = make(map[string]types.RegisterPatientPayload)

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

	
	// ✅ 4. حفظ بيانات المريض مؤقتًا بانتظار التحقق من OTP
	pendingPatients[patientPayload.Email] = patientPayload

	// ✅ 5. إرسال رسالة انتظار التحقق
	utils.WriteJSON(w, http.StatusAccepted, map[string]string{
		"message": "OTP sent. Please verify to complete registration.",
	})



	
}



// ----------------------------------------------------



func (h *Handler) VerifyOTPHandler(w http.ResponseWriter , r *http.Request) {
	var optCodePayload types.VerifyRequest
	if err := utils.ParseJSON(r , &optCodePayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

		//validate the payoad .....................
		if err := utils.Validate.Struct(optCodePayload);err != nil {
			error := err.(validator.ValidationErrors)
			  utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid payload %v", error) )
			return
		}

		if optCodePayload.OTPCode != "666666" {
			utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid OTP Code"))
			return
		}

			


		patientPayload, exists := pendingPatients[optCodePayload.Email]
		if !exists {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("no registration found for this email"))
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
	    IDNumber: patientPayload.IDNumber,
		IsCompleted: false,
		CenterID: cenetr.ID,
		City: patientPayload.City,
		
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	
	delete(pendingPatients,optCodePayload.Email)


	utils.WriteJSON(w , http.StatusCreated , map[string]string{"message":"successfully Created"})

}










//




func (h *Handler) handleGetPatient (w http.ResponseWriter , r *http.Request) {
	
	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}


	patient , err := h.store.GetPatientById(id)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}



	utils.WriteJSON(w , http.StatusOK , patient)



}


//router agign



// handle get all patient info ..........................................


func (h *Handler) handleGetAllPatientInfo (w http.ResponseWriter , r *http.Request) {

	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}

	patientDetials , err := h.store.GetPatientDetailsByID(id)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}


	utils.WriteJSON(w , http.StatusOK , patientDetials)


}




func (h *Handler)  VerifyTokenHandler(w http.ResponseWriter , r *http.Request) {
	tokenString := auth.GetTokenFromRequest(r)
	if tokenString == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	token, err := auth.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Token expired or invalid", http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(w , http.StatusOK , "Token is Vaild")
}


// update patient profile
func (h *Handler) handleUpdatePatientProfile(w http.ResponseWriter , r *http.Request) {
	var updatePatietPayload types.ParientUpdatePayload
	if err := utils.ParseJSON(r , &updatePatietPayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

	err := h.store.UpdatePatientProfile(updatePatietPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error update patient profile"))
		return
	}
    
	utils.WriteJSON(w , http.StatusOK , map[string]string{
		"message": "patient profile update successfully",
	})

}













	
// secret := []byte(config.Envs.JWTSecret)
// token , err := auth.CreateJWT(secret , center.ID)

// if err != nil {
// 	utils.WriteError(w , http.StatusInternalServerError ,err)
// 	return
// }

// patients , err := h.store.GetPatientsForCenter(center.ID)
// if err != nil {
	
// }


// returnLoggingData := types.ReturnLoggingCenterData {
// 	Name: center.CenterName,
// 	Email: center.CenterEmail,
// 	Role: "center",
// 	Patient: patients,
// 	Token: token,
// }

// utils.WriteJSON(w , http.StatusOK ,returnLoggingData)


func (h *Handler) handleGetPatientProfile(w http.ResponseWriter , r *http.Request) {
	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}
	patientProfile , err := h.store.GetPatientProfile(id)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}



	utils.WriteJSON(w , http.StatusOK , patientProfile)
}