package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/types"
)


func SendOTP(ToEmail, otpCode string , companyName string , userName string) error {

	apiKey :=config.GetEnv("SENDGRID_API_KEY" , "")
	url := "https://api.brevo.com/v3/smtp/email"

	emailRequest := types.EmailRequest{
		Sender: types.SenderInfo{
			Name: companyName,
			Email: "diabetes.care.center.syria@gmail.com",
		},
		To: []types.Recipient{
			{Email:ToEmail, Name: userName},
		},
		Subject:     "Hello from center!",
		HTMLContent: fmt.Sprintf("<h1>Hello, this is your otp code is : %s !</h1>",otpCode),
	}


	payload, err := json.Marshal(emailRequest)
	if err != nil {
		return err
	}


	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// تعيين الرؤوس المطلوبة
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	// إرسال الطلب
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in sending otp", err)
		return err
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == http.StatusCreated {
		return nil
	} 

	return nil

}