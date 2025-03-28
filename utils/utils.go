package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AL-Hourani/care-center/mail"
	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/AL-Hourani/care-center/types"
	"github.com/go-playground/validator/v10"
)


var Validate = validator.New()

func ParseJSON(r *http.Request , payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}


func WriteJSON(w http.ResponseWriter , status int , v any) error {
	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter , status int , err error) {
	WriteJSON(w , status , map[string]string{"error":err.Error()})
}




//send opt to user 

func RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req types.OTPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	
	otpCode, err := auth.GenerateOTP(req.Email)
	if err != nil {
		http.Error(w, "Could not generate OTP", http.StatusInternalServerError)
		return
	}

	
	err = mail.SendOTP(req.Email, otpCode)
	if err != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OTP has been sent to your email"))
}

//check from opt 

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req types.VerifyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.OTPCode == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// التحقق من الرمز
	if auth.VerifyOTP(req.Email, req.OTPCode) {
		w.Write([]byte("OTP verified successfully"))
	} else {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
	}
}





