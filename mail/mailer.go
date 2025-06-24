package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AL-Hourani/care-center/service/auth"
)

func Mailer(email string, name string) error {
    // Generate OTP
    otp, err := auth.GenerateOTP(email)
    if err != nil {
        return fmt.Errorf("failed to generate OTP: %v", err)
    }
url := "https://api.emailjs.com/api/v1.0/email/send"

    data := map[string]interface{}{
        "service_id":  "service_cruvhkn",
        "template_id": "template_cugk7wr",
        "user_id":     "ytF93swbS8Hw5negh", // هذا هو api key
        "template_params": map[string]string{
            "name":     name,
            "otp":      otp,
            "to_email": email,
        },
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
    }

    return nil

}
