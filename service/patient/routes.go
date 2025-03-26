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
	router.HandleFunc("/setPatientHealthInfo",h.handleSetPatientHealthInfo).Methods("POST")
	router.HandleFunc("/setPatientPersonalInfo",h.handleSetPatientPersonalInfo).Methods("POST")
	router.HandleFunc("/getPatient/{id}" , h.handleGetPatient).Methods("GET")
	router.HandleFunc("/getAllPatientInfo/{id}" , h.handleGetAllPatientInfo).Methods("GET")
	router.HandleFunc("/Logout" , h.handleGetAllPatientInfo).Methods("GET")
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

		returnLoggingData := types.ReturnLoggingData {
			Name: patient.FullName,
			Email: patient.Email,
			Role: "patient",
			IsCompletes: false,
			Token: token,
		}
		

	    
		utils.WriteJSON(w , http.StatusOK ,returnLoggingData)
		
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

		patients , err := h.store.GetPatientsForCenter(center.ID)
		if err != nil {
			
		}


		returnLoggingData := types.ReturnLoggingCenterData {
			Name: center.CenterName,
			Email: center.CenterEmail,
			Role: "center",
			IsCompletes: true,
			Patient: patients,
			Token: token,
		}
	
		utils.WriteJSON(w , http.StatusOK ,returnLoggingData)
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
	    IDNumber: patientPayload.IDNumber,
		IsCompleted: false,
		CenterID: cenetr.ID,
		
	})

	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}

	utils.WriteJSON(w , http.StatusCreated , map[string]string{"message":"successfully Created"})
}





//


func (h *Handler) handleSetPatientPersonalInfo(w http.ResponseWriter , r *http.Request) {
	var personalPatientInfo types.BasicPatientInfoPalyoad
		//get json payload
		if err := utils.ParseJSON(r , &personalPatientInfo); err != nil {
			utils.WriteError(w , http.StatusBadRequest , err)
			return
		}

	//validate the payoad .....................
	if err := utils.Validate.Struct(personalPatientInfo);err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid payload %v", error) )
		return
	}



	_ , err := h.store.GetPatientById(personalPatientInfo.PatientID)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("patient with id %d is not exists" ,personalPatientInfo.PatientID ))
		return 
	}


	err = h.store.SetPersonlPatientBasicInfo(types.BasicPatientInfo{
        PatientID: personalPatientInfo.PatientID,
		Weight: personalPatientInfo.Weight,
		Length: personalPatientInfo.Length,
		Address: personalPatientInfo.Address,
		Gender: personalPatientInfo.Gender,
		
	})

	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}


    utils.WriteJSON(w , http.StatusCreated , map[string]string{"message":"successfully add basic info"})


}



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
func (h *Handler) handleSetPatientHealthInfo(w http.ResponseWriter , r *http.Request) {
	var patientHealtPayload  types.RegisterHealthPatientData
	//get json payload
			if err := utils.ParseJSON(r , &patientHealtPayload); err != nil {
				utils.WriteError(w , http.StatusBadRequest , err)
				return
			}

    

	//validate the payoad .....................

	if err := utils.Validate.Struct(patientHealtPayload);err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid payload %v", error) )
		return
	}

	//check if patient is exists 

	_ , err := h.store.GetPatientById(patientHealtPayload.PatientID)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("patient with id %d is not exists" ,patientHealtPayload.PatientID ))
		return 
	}

	err = h.store.SetPatientHealthInfo(types.HealthPatientData {
		PatientID: patientHealtPayload.PatientID,
		BloodSugar: patientHealtPayload.BloodSugar,
		Hemoglobin: patientHealtPayload.Hemoglobin,
		BloodPressure: patientHealtPayload.BloodPressure,
		SugarType: patientHealtPayload.SugarType,
		DiseaseDetection: patientHealtPayload.DiseaseDetection,
		OtherDisease: patientHealtPayload.OtherDisease,
		TypeOfMedicine: patientHealtPayload.TypeOfMedicine,
		UrineAcid: patientHealtPayload.UrineAcid,
		Cholesterol: patientHealtPayload.Cholesterol,
		Grease: patientHealtPayload.Grease,
		HistoryOfFamilyDisease: patientHealtPayload.HistoryOfFamilyDisease,
	})

	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}


    utils.WriteJSON(w , http.StatusCreated , map[string]string{"message":"successfully add health info"})


}




// handle get all patient info ..........................................



func (h *Handler) handleGetAllPatientInfo (w http.ResponseWriter , r *http.Request) {


	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}

	AllpatientInfo , err := h.store.GetAllPatientInfo(id)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}



	utils.WriteJSON(w , http.StatusOK ,AllpatientInfo)


}


