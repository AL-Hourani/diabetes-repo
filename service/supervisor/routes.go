package supervisor

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/AL-Hourani/care-center/types"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.CenterStore
    pStore types.PatientStore
    superStore types.SuperisorStore
}

func NewHandler(store types.CenterStore , patientStore types.PatientStore , supervisorStore types.SuperisorStore) *Handler {
	return &Handler {
		store: store,
		pStore: patientStore,
	    superStore: supervisorStore,
	}
}




func (h *Handler) RegisterSuperVisorRoutes(router *mux.Router) {
	router.HandleFunc("/getSupervisorCenters",auth.WithJWTAuth(h.handleGetAllCentersData)).Methods("GET")
	router.HandleFunc("/getInquiries",auth.WithJWTAuth(h.handleGetInquiries)).Methods("GET")
	router.HandleFunc("/getInquiriesDetails/{id}",auth.WithJWTAuth(h.handleGetInquiriesDetails)).Methods("GET")
	router.HandleFunc("/superLogin",h.handleLoginSupervisor).Methods("POST")
}







func (h *Handler) handleGetAllCentersData(w http.ResponseWriter, r *http.Request) {
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

	user , err := h.pStore.GetLoginByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if user.Role != "supervisor" {
		http.Error(w, "Unauthorized: You are not supervisor", http.StatusUnauthorized)
		return
	}

	centerData , err := h.superStore.GetAllCenters()
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w , http.StatusOK , centerData	)

}




func (h *Handler) handleGetInquiries(w http.ResponseWriter, r *http.Request) {
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

	user , err := h.pStore.GetLoginByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if user.Role != "supervisor" {
		http.Error(w, "Unauthorized: You are not supervisor", http.StatusUnauthorized)
		return
	}

    InquiriesData , err := h.superStore.GetAllInformation()
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w , http.StatusOK , InquiriesData	)

}

func (h *Handler) handleGetInquiriesDetails(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(auth.UserContextKey).(*jwt.Token)
	if !ok {
		http.Error(w, "Unauthorized: No token found", http.StatusUnauthorized)
		return
	}

	idSup, err := auth.GetIDFromToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user , err := h.pStore.GetLoginByID(idSup)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if user.Role != "supervisor" {
		http.Error(w, "Unauthorized: You are not supervisor", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}



    Inquirie , err := h.superStore.GetInformationByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	center , err := h.superStore.GetCenterByID(Inquirie.CenterID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	nop , err := h.superStore.CountPatientsByCenter(center.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	cQuentitiy , err  := h.superStore.GetMedicationByArabicName(Inquirie.NameArabic , center.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	currentquentitiy , err  := strconv.Atoi(cQuentitiy.Quantity)
	if err != nil {
        fmt.Println("خطأ في التحويل:", err)
        return
    }

	RequestDetails , err := h.superStore.GetMedicationRequestByID(Inquirie.RequestId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	newQueryDetails := types.InquirieDetails {
		NameArabic: Inquirie.NameArabic,
		NameEnglish: Inquirie.NameEnglish,
		CenterName: center.CenterName,
		CenterCity: center.CenterCity,
		RQuantity: Inquirie.Quantity,
		CQuantity: currentquentitiy,
		Nop: nop,
		Request_date: RequestDetails.RequestedAt,

	}

	utils.WriteJSON(w , http.StatusOK , newQueryDetails)

}




func (h *Handler) handleLoginSupervisor(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.pStore.GetUserByEmail(LoginPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

		
		if !auth.ComparePasswords(user.Password, []byte(LoginPayload.Password)) {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
			return
		}


		
		if user.Role == "supervisor" {
				secret := []byte(config.Envs.JWTSecret)
				token, err := auth.CreateJWT(secret, user.ID)
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, err)
					return
				}

				ruturnData := types.Supervisor {
					Email: user.Email,
					Role: user.Role,
					Token: token,
				}

				utils.WriteJSON(w, http.StatusOK, ruturnData)
				return

			}else {
				utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: You are not supervisor"))
				return
			}


}