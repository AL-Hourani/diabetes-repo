package supervisor

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

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
	router.HandleFunc("/getCityInfo",auth.WithJWTAuth(h.handleGetCentersByCity)).Methods("GET")
	router.HandleFunc("/getSuperInfo",auth.WithJWTAuth(h.handleGetSuperInfo)).Methods("GET")
	router.HandleFunc("/rejectInquiries",auth.WithJWTAuth(h.handleRejectInquiries)).Methods("POST")
	router.HandleFunc("/acceptedInquiries",auth.WithJWTAuth(h.handleAcceptedInquiries)).Methods("POST")
	router.HandleFunc("/CreateDatePatientFile",auth.WithJWTAuth(h.handleGetSuperExcel)).Methods("POST")
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






func (h *Handler) handleRejectInquiries(w http.ResponseWriter, r *http.Request) {
	var queryId types.QueryID

	if err := utils.ParseJSON(r , &queryId); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

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


	err = h.superStore.UpdateInformationStatus(queryId.Query_ID , string(types.InfoStatusCancel))
    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	err = h.superStore.UpdateRecordStatusAndApprovalDate(queryId.Query_ID , string(types.StatusRejected))
    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

		utils.WriteJSON(w , http.StatusOK , map[string]string{
		"message": "query reject successfully",
	})
}








func (h *Handler) handleAcceptedInquiries(w http.ResponseWriter, r *http.Request) {
	var queryAccept types.QueryAccepted

	if err := utils.ParseJSON(r , &queryAccept); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

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


	err = h.superStore.UpdateInformationStatus(queryAccept.Query_ID , string(types.InfoStatusOK))
    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	query , err := h.superStore.GetInformationByID(queryAccept.Query_ID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	err = h.superStore.UpdateRecordStatusAndApprovalDate(query.RequestId, string(types.StatusApproved))
    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
    
	err = h.superStore.UpdateMedicationQuantity(query.RequestId ,query.Quantity )
    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

		utils.WriteJSON(w , http.StatusOK , map[string]string{
		"message": "query Accepted successfully",
	})
}








func (h *Handler) handleGetCentersByCity(w http.ResponseWriter, r *http.Request) {
    cityName := r.URL.Query().Get("city")


	if err := utils.ParseJSON(r , &cityName); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

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


    centers_name_by_city , err := h.superStore.GetCentersByCity(cityName)
    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	nopin , err := h.superStore.GetPatientCountByCity(cityName)
	    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	nopin_lc , err := h.superStore.GetPatientCountByCityLastMonth(cityName)
	    if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	newCityInfo := types.AllCityInfo {
		NumberOfPatientInCity:nopin,
		NumberOfPatientInCityLastMonth:nopin_lc ,
		ActiveCenter:centers_name_by_city ,
	}

	
	utils.WriteJSON(w , http.StatusOK ,newCityInfo )
}









func (h *Handler) handleGetSuperInfo(w http.ResponseWriter, r *http.Request) {

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


	var ( 
		citiesInSystem []string
		nopis          int
		nopisLastMonth int
		firstCenter    types.CenterWithCount
	)


	errCh := make(chan error, 4)

	var wg sync.WaitGroup
	wg.Add(4)


	go func() {
		defer wg.Done()
		cities, err := h.store.GetCities()
		if err != nil {
			errCh <- err
			return
		}
		citiesInSystem = cities
	}()

	go func() {
		defer wg.Done()
		total, err := h.pStore.GetTotalPatientsInSystem()
		if err != nil {
			errCh <- err
			return
		}
		nopis = total
	}()

	go func() {
		defer wg.Done()
		lastMonth, err := h.pStore.GetPatientsLastMonth()
		if err != nil {
			errCh <- err
			return
		}
		nopisLastMonth = lastMonth
	}()


	go func() {
		defer wg.Done()
		center, err := h.superStore.GetCenterWithMostPatients()
		if err != nil {
			errCh <- err
			return
		}
		firstCenter = *center
	}()


	wg.Wait()
	close(errCh)

	
	if err := <-errCh; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	newSystemInfo := types.GenericSuperInfo {
	   NumberOfPatientInSystem:nopis,
       NumberOfPatientInSystemLastMonth:nopisLastMonth ,
	   ActiveCities:citiesInSystem ,
	   FirstCenter: firstCenter,
	}

	
	utils.WriteJSON(w , http.StatusOK ,newSystemInfo )
}








func (h *Handler) handleGetSuperExcel(w http.ResponseWriter, r *http.Request) {
	var date types.MonthDown
	
	if err := utils.ParseJSON(r , &date); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}
    token, ok := r.Context().Value(auth.UserContextKey).(*jwt.Token)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    idSup, err := auth.GetIDFromToken(token)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    user, err := h.pStore.GetLoginByID(idSup)
    if err != nil || user.Role != "supervisor" {
        http.Error(w, "Unauthorized: not supervisor", http.StatusUnauthorized)
        return
    }

    month, year, err := h.superStore.ParseMonthYear(date.MonthDown)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	} else {
		fmt.Println("Month:", month, "Year:", year)
	}
	// المخرجات: Month: 3 Year: 2025


    reviews, err := h.superStore.GetPatientReviewsByMonth(month, year)
    if err != nil {
        http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
        return
    }

    excelFile, err := CreateExcelFile(reviews)
    if err != nil {
        http.Error(w, "Error creating Excel file", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", "attachment; filename=patient_reviews.xlsx")
    w.Header().Set("File-Name", "patient_reviews.xlsx")

    if err := excelFile.Write(w); err != nil {
        http.Error(w, "Error writing Excel file", http.StatusInternalServerError)
        return
    }
}
