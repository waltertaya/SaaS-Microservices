package utils

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/waltertaya/saas-microservices/auth-service/initializers"
)

func SendEmail(to, otp string) error {

	initializers.LoadEnvVariables()

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	subject := "Your OTP Code for Verification"

	msg := "From: " + from + "\n" +
	"To: " + to + "\n" +
	"Subject: " + subject + "\n" +
	"MIME-Version: 1.0\n" +
	"Content-Type: text/html; charset=\"UTF-8\"\n" +
	"Content-Transfer-Encoding: 8bit\n\n" +
	"<html><body>" +
	"<h1>Your OTP Code</h1>" +
	"<p>Your OTP code is: <strong>" + otp + "</strong></p>" +
	"<p>This code is valid for 10 minutes.</p>" +
	"<p>If you did not request this code, please ignore this email.</p>" +
	"<p>Thank you for using our service!</p>" +
	"</body></html>"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	fmt.Println("Email sent successfully!")
	return nil
}
