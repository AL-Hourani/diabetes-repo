package mail

import (
	"fmt"
	"log"

	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/resend/resend-go/v2"
)

func Mailer(email string , name string) error {

    apiKey := "re_Mqmq3x92_4y5ZSCrJh8bGX9xJV4MSMPV1"

    client := resend.NewClient(apiKey)

	//generate api ........
	otp , err  := auth.GenerateOTP(email)
	subject := fmt.Sprintf("مرحبًا %s، هذا رمز التحقق الخاص بك", name)
	if err != nil {
		log.Fatal("error generate otp")
	}

    params := &resend.SendEmailRequest{
        From:    "onboarding@resend.dev",
        To:      []string{email},
        Subject: subject,
        Html:    fmt.Sprintf("<h2>رمز التحقق الخاص بك هو: <strong>%s</strong></h2>", otp),
    }

    _ , err = client.Emails.Send(params)
	if err != nil {
        log.Fatalf("فشل الإرسال: %v", err)
    }

	return err
}