package center

import (
	"fmt"
	"net/http"

	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/AL-Hourani/care-center/types"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
    
)

type Handler struct {
	store types.CenterStore
}

func NewHandler(store types.CenterStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterCenterRoutes(router *mux.Router) {
	router.HandleFunc("/centerLogin", h.handleCenterLogin).Methods("POST")
	router.HandleFunc("/centerRegister", h.handleCenterRegister).Methods("POST")
	router.HandleFunc("/getPatients", h.handleGetPatients).Methods(http.MethodGet)
	router.HandleFunc("/getCenters", h.handleGetCenters).Methods(http.MethodGet)
}



func (h *Handler) handleCenterLogin(w http.ResponseWriter , r *http.Request) {
		//get json payload
		var centerPayload types.LoginCenterPayload
		if err := utils.ParseJSON(r , &centerPayload); err != nil {
			utils.WriteError(w , http.StatusBadRequest , err)
		}

		//validate the payoad .....................
	if err := utils.Validate.Struct(centerPayload);err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid payload %v", error) )
		return
	}

	//find center
	center , err := h.store.GetCenterByName(centerPayload.CenterName)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("not found invalid center name or center password"))
		return 
	}

	if !auth.ComparePasswords(center.CenterPassword , [] byte(centerPayload.CenterPassword)) {
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
        _ , err := h.store.GetCenterByName(centerPayload.CenterName)
		if err == nil {
			utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("center with name %s already exists" , centerPayload.CenterName))
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
		CenterAddress: centerPayload.CenterAddress,

	})

	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}

	utils.WriteJSON(w , http.StatusCreated ,  map[string]string{"message":"successfully Created"})
}

func (h *Handler) handleGetPatients(w http.ResponseWriter , r *http.Request) {
	    var centerName string
	
		if err := utils.ParseJSON(r , &centerName); err != nil {
			utils.WriteError(w , http.StatusBadRequest , err)
		}
		// check if center exists
        center , err := h.store.GetCenterByName(centerName)
		if err != nil {
			utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("center with name %s is not exists" , centerName))
			return 
		}

	patientsList , err := h.store.GetPatients(center.ID)
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