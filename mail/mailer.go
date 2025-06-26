package mail

import (
	
	"fmt"
	"net/http"
    "net/url"
    "strings"
	"github.com/AL-Hourani/care-center/service/auth"
)

func Mailer(email string, name string) error {
    // Generate OTP
    otp, err := auth.GenerateOTP(email)
    if err != nil {
        return fmt.Errorf("failed to generate OTP: %v", err)
    }
     
    apiKey := "7813682AEAD1740635FBAAB6BF82A8442C9AA5831B2137C661098F50ED8524D80B6DADD51B670609ECE28FE71D104E49"
    from := "diabetes.care.center.syria@gmail.com" // البريد الذي تحققت منه في Elastic Email

    data := url.Values{}
    data.Set("apikey", apiKey)
    data.Set("from", from)
    data.Set("to", email)
    data.Set("subject", "رمز التحقق OTP")
    data.Set("bodyText", fmt.Sprintf("مرحبًا %s، رمز التحقق الخاص بك هو: %s", name, otp))

    data.Set("isTransactional", "true")

    req, err := http.NewRequest("POST", "https://api.elasticemail.com/v2/email/send", strings.NewReader(data.Encode()))
    if err != nil {
        fmt.Println("Error creating request:", err)
        return err
    }

    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return err
    }
    defer resp.Body.Close()

 

    return nil

}
