package mail

import (
	"fmt"
	"net/smtp"

	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/service/auth"
)

func Mailer(email string, name string) error {
    // Generate OTP
    otp, err := auth.GenerateOTP(email)
    if err != nil {
        return fmt.Errorf("failed to generate OTP: %v", err)
    }

    // SMTP إعدادات
    smtpHost := "smtp.mailersend.net"
    smtpPort := "587"
    smtpUser := config.Envs.MAILERSEND_SMTP_USER   
    smtpPassword :=config.Envs.MAILERSEND_SMTP_USER

    // محتوى الرسالة
    from := smtpUser
    to := email
    subject := "رمز التحقق الخاص بك"
    body := fmt.Sprintf("مرحبًا %s، رمز التحقق الخاص بك هو: %s", name, otp)

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

    // إرسال الرسالة
    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        return fmt.Errorf("فشل في إرسال البريد: %v", err)
    }

    fmt.Println("تم إرسال البريد الإلكتروني بنجاح إلى:", email)
    return nil
}
