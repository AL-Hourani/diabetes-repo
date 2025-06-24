package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AL-Hourani/care-center/service/auth"
)

func Mailer(email string, name string) error {
    // Generate OTP
    otp, err := auth.GenerateOTP(email)
    if err != nil {
        return fmt.Errorf("failed to generate OTP: %v", err)
    }

        body := map[string]string{
        "from":    "noreply@onresend.com",
        "to":      email,
        "subject": "رمز التحقق الخاص بك",
        "html":     "<p>مرحبًا " + name + " </p><p><strong>رمز التحقق هو: " + otp + "</strong></p>",
    }

    jsonData, _ := json.Marshal(body)

    req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer  re_hxapdzCX_LvfsU4bqcsBjT54J8d4C8Ftr")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("فشل الإرسال:", err)
        return err
    }
    defer resp.Body.Close()
    log.Println("تم إرسال OTP إلى:", email)

    return nil

}
