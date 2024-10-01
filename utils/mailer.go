package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
)

func SendMail(email []string, payload bytes.Buffer) {
	mailHost := Settings.MailHost
	mailPort := Settings.MailPort
	address := fmt.Sprintf("%v:%v", mailHost, mailPort)
	// Set up SMTP client and send the email (replace with your actual SMTP settings)
	auth := smtp.PlainAuth("", Settings.MailUsername, Settings.MailPassword, Settings.MailHost)
	err := smtp.SendMail(address, auth, Settings.MailFrom, email, payload.Bytes())
	if err != nil {
		log.Fatal("mail failed to send")
	}
}
