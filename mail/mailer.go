package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
)


func SendOTP(email, otpCode string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "your-email@example.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Your OTP Code")
	mailer.SetBody("text/html", fmt.Sprintf(`
	    <h3>مرحبا بك في مركز عناية</h3>
        <h3>Your OTP code is:</h3>
        <h2>%s</h2>
    `, otpCode))

	dialer := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-email-password")

	return dialer.DialAndSend(mailer)

}