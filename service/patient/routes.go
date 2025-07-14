package patient

import (
	"fmt"
	"net/http"
	"time"

	"strconv"

	"github.com/AL-Hourani/care-center/config"

	"github.com/golang-jwt/jwt/v5"

	// "github.com/AL-Hourani/care-center/mail"
	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/AL-Hourani/care-center/service/notifications"
	"github.com/AL-Hourani/care-center/service/session"

	// "github.com/AL-Hourani/care-center/service/patient"
	"github.com/AL-Hourani/care-center/types"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.PatientStore
	storeCenter types.CenterStore
	SessionManager *session.Manager
	NotifHub      *notifications.Hub
}

func NewHandler(store types.PatientStore , centerStore types.CenterStore , sessionManager session.Manager  , notifHub *notifications.Hub) *Handler {
	return &Handler{store: store , storeCenter: centerStore , SessionManager: &sessionManager , NotifHub:       notifHub, }
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
	router.HandleFunc("/CenterStatistics/{id}",h.handleStatisticsSugerType).Methods("GET")
	router.HandleFunc("/sendEmail",h.handleVerifyEmail).Methods("POST")
	router.HandleFunc("/verfiyOTPResetPassword",h.handleVerifyOTP).Methods("POST")
	router.HandleFunc("/resetPassword",h.handleResetPassword).Methods("POST")
	router.HandleFunc("/gethomePatient" ,  auth.WithJWTAuth(h.handleGethomePatient)).Methods("GET")

	router.HandleFunc("/getNotifications" ,  auth.WithJWTAuth(h.handleGetPatientNotifications)).Methods("GET")
	router.HandleFunc("/notifications/mark-all-read" ,  auth.WithJWTAuth(h.MarkAllNotificationsAsRead)).Methods("PUT")
   
	router.HandleFunc("/UpdatePatientProfile" ,  auth.WithJWTAuth(h.handleUpdatePatient)).Methods("POST")

	

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

		
		if !auth.ComparePasswords(user.Password, []byte(LoginPayload.Password)) {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
			return
		}



		if user.Role == "patient" {
		   patient , err  := h.store.GetPatientByEmail(user.Email)
		   if err != nil {
			   utils.WriteError(w, http.StatusInternalServerError, err)
			   return
		   }
		   
				// إنشاء JWT Token
				secret := []byte(config.Envs.JWTSecret)
				token, err := auth.CreateJWT(secret, patient.ID)
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, err)
					return
				}

			returnLoggingData := types.ReturnLoggingData{
				ID:          patient.ID,
				Name:        patient.FullName,
				Email:       user.Email,
				Role:        user.Role,
				IsCompletes:  false,
				Token:        token,
			}
			utils.WriteJSON(w, http.StatusOK, returnLoggingData)
			return 
		} else {

		  center , err  := h.storeCenter.GetCenterByEmail(user.Email)
		   if err != nil {
			   utils.WriteError(w, http.StatusInternalServerError, err)
			   return
		   }
	   			
				secret := []byte(config.Envs.JWTSecret)
				token, err := auth.CreateJWT(secret, center.ID)
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, err)
					return
				}

	
			returnLoggingData := types.ReturnLoggingCenterData{
			    ID:      center.ID,      
				Name:    center.CenterName,
				Email:   user.Email,
				Role:    user.Role,
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

	// err = mail.Mailer(patientPayload.Email , patientPayload.FullName)
	// 	if err != nil {
	// 		utils.WriteError(w, http.StatusInternalServerError, err)
	// 		return
	// 	}

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


			


		patientPayload, exists := pendingPatients[optCodePayload.Email]
		if !exists {
			utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("no email registered"))
			return
		}
		// if !auth.VerifyOTP(optCodePayload.Email , optCodePayload.Email) {
		// 	utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid OTP Code"))
		// 	return
		// }

		if optCodePayload.OTPCode != "666666" {
			utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid OTP Code"))
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
	id , err := h.store.GreatePatient(types.Patient{
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

	newLoginFailed := types.InsertLogin {
		Email:patientPayload.Email ,
		Password:hashedPassword ,
	}
	
	err = h.storeCenter.GreateLoginFailed(newLoginFailed)
	

	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}

	message := fmt.Sprintf("مرحبًا %s، أهلًا بك في %s.\nيسعدنا انضمامك إلينا، صحتك أمانة في أعيننا. ", patientPayload.FullName, patientPayload.CenterName)

		h.NotifHub.Broadcast <- types.Notification{
		SenderID:  1  , 
		ReceiverID: id,
		Message:    message,
		IsRead: false,
		CreatedAt: FormatRelativeTime(time.Now()),

	}

	
	_ = h.storeCenter.InsertNotification(types.NotificationTwo{
		SenderID:   1,
		ReceiverID: id,
		Message:    message,
		})



	
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

	
	id, err := auth.GetIDFromToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	_ , err = h.store.GetPatientById(id)
		if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	_ , err = h.storeCenter.GetCenterByID(id)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
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
    
	
	updateProfileInfo , err := h.store.GetPatientProfile(updatePatietPayload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest ,  fmt.Errorf("error get patient profile"))
		return
	}

	
	utils.WriteJSON(w , http.StatusOK , updateProfileInfo)
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




// handleStatisticsSugerType
func (h *Handler) handleStatisticsSugerType(w http.ResponseWriter , r *http.Request) {
	vars := mux.Vars(r)
	id , err := strconv.Atoi(vars["id"])
	if err != nil  {
       utils.WriteError(w, http.StatusBadRequest , fmt.Errorf("invalid ID"))
       return
	}
	statisticSuger , err := h.store.GetSugarTypeStats(id)
	if err != nil {
		utils.WriteError(w , http.StatusBadRequest ,err)
		return 
	}
		// Get gender counts
	maleCount, femaleCount, err := h.store.GetGenderCounts(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	
	totalSystemPatients, err := h.store.GetTotalPatientsInSystem()
	if err != nil {
		utils.WriteError(w , http.StatusInternalServerError ,err)
		return 
	}

	sugarAgeStats, err := h.store.GetSugarTypeAgeRangeStats(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	totoalsugarAgeStats, err := h.store.GetSugarTypeAgeRangeStatsAllSystem()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	bmiStats, err := h.store.GetBMIStats(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	cityStats, err := h.store.GetCityStats()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}


	response := types.SugarStatsResponse{
		SugarTypes:  statisticSuger,
		MaleCount:   maleCount,
		FemaleCount: femaleCount,
		TotalCount: maleCount + femaleCount,
		TotalPatientsInSystem: totalSystemPatients,
		SugarAgeRangeStats:     sugarAgeStats,
		TotalSugarAgeRangeStats: totoalsugarAgeStats,
		BMIStats: bmiStats,
	    CityStats: cityStats,
	}

	utils.WriteJSON(w, http.StatusOK, response)

}





//reste password

func (h *Handler) handleVerifyEmail(w http.ResponseWriter , r *http.Request) {
	var emailPayload types.Email
	if err := utils.ParseJSON(r , &emailPayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

	err := h.store.GetUserByEmailRestPassword(emailPayload.Email)
	if err != nil {
		utils.WriteError(w , http.StatusNotFound , err)
	}

	// send otp 


	err = h.SessionManager.SetValue(w, r, "reset-session", "resetEmail", emailPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}


	utils.WriteJSON(w, http.StatusAccepted, map[string]string{
		"message": "OTP sent. Please verify to complete Reset Password.",
	})


}

func (h *Handler) handleVerifyOTP(w http.ResponseWriter , r *http.Request) {
	var otpPayload types.OTPResetPass
	if err := utils.ParseJSON(r , &otpPayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}

	if otpPayload.OTP != "666666" {
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid otp"))
		return
	}


	// send otp 

	utils.WriteJSON(w, http.StatusAccepted, map[string]string{
		"message": "OTP verify successfully  ('_')",
	})


}

func (h *Handler) handleResetPassword(w http.ResponseWriter , r *http.Request) {
	var resetPasswordPayload types.ResetPassword
	if err := utils.ParseJSON(r , &resetPasswordPayload); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}


	// rest password......
	err := h.store.UpdatePasswordByEmail(resetPasswordPayload.Email, resetPasswordPayload.NewPassword)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }


	
   
    utils.WriteJSON(w, http.StatusOK, map[string]string{
     "message": "Password updated successfully.",
    })   


}






// patient page app ...


func (h *Handler) handleGethomePatient(w http.ResponseWriter , r *http.Request) {
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

	patient , err  := h.store.GetPatientById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		 return
	}
    

	reviews , err := h.store.GetReviewsByPatientID(id)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		 return
	}

	var chartData []types.ChartData
    var myReviews []types.Review
	for _, r := range reviews {
	// 1. تعبئة ChartData
	chartData = append(chartData, types.ChartData{
		Date:          r.DateReview.Format("02-01-2006"),
		LDL:           r.LDL,
		HDL:           r.HDL,
		NormalGlocose: r.NormalGlocose,
	})

	// 2. تعبئة MyReviews
	myReviews = append(myReviews, types.Review{
		Id:   r.ID,
		Date: r.DateReview.Format("02-01-2006"),
	})
}



    var next_Reviwe string
    var first_Reviwe string
	if len(reviews) > 0 {
		// نأخذ أول مراجعة (يفترض أنها مرتبة تنازليًا)
		first_Reviwe =  reviews[0].DateReview.Format("02-01-2006")
		latest := reviews[0].DateReview
		next_Reviwe = latest.AddDate(0, 1, 0).Format("02-01-2006")
		
	} else {
    // لا توجد مراجعات، نستخدم تاريخ اليوم
    first_Reviwe = "لا يوجد اي مراجعة يرجى زيازة المركز  لاكمال الحساب"
    next_Reviwe = time.Now().Format("02-01-2006")
  
}


    center , err := h.storeCenter.GetCenterByID(patient.CenterID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		 return
	}


	newGetPatientHomeData := types.GetPatientHomeData {
		FullName: patient.FullName,
		Age: patient.Age,
		IDNumber: patient.IDNumber,
		FirstReviewDate: first_Reviwe,
		ChartData: chartData,
		MyCenter: center.CenterName ,
		NextReview: next_Reviwe ,
		MyReviews:myReviews,

	}



		utils.WriteJSON(w, http.StatusOK, newGetPatientHomeData)
}














func (h *Handler) handleGetPatientNotifications(w http.ResponseWriter , r *http.Request) {
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


	notifications, err := h.store.GetNotificationsByUserID(id)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }

    utils.WriteJSON(w, http.StatusOK, notifications)
}






func (h *Handler) MarkAllNotificationsAsRead(w http.ResponseWriter , r *http.Request) {
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


	 err = h.store.UpdateIsReadNotifications(id)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }

      utils.WriteJSON(w, http.StatusOK, map[string]string{
     "message": "update success",
    })   

}







func (h *Handler) handleUpdatePatient(w http.ResponseWriter, r *http.Request) {
	var updateData types.UpdatePatientInfo
	token, ok := r.Context().Value(auth.UserContextKey).(*jwt.Token)
	if !ok {
		http.Error(w, "Unauthorized: No token found", http.StatusUnauthorized)
		return
	}

	patientID, err := auth.GetIDFromToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	
	if err := utils.ParseJSON(r , &updateData); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
		return
	}


	updatedPatient, err := h.store.UpdatePatientBasicInfo(updateData , patientID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedPatient)


}




