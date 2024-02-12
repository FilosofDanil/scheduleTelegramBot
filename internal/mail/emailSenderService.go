package mail

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/smtp"
	"schedulerTelegramBot/configs"
)

type EmailSender struct {
	configurations configs.EmailConfigs
	id             string
}

func (e *EmailSender) SendMail() error {
	// Sender data.
	from := e.configurations.Sender
	password := e.configurations.Password
	id := uuid.New()
	e.id = id.String()
	fmt.Println(id.String())
	fmt.Println(e.id)
	// Receiver email address.
	to := []string{
		e.configurations.Receiver,
	}
	// smtp server configuration.
	smtpHost := e.configurations.SmtpHost
	smtpPort := e.configurations.SmtpPort

	// Message.

	message := []byte("Verification code: " + id.String())

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}

func (e *EmailSender) CheckUUID(idFromMessage string) error {
	fmt.Println(e.id)
	fmt.Println(idFromMessage)
	if e.id != idFromMessage {
		return errors.New("incorrect uuid")
	}
	return nil
}

func NewEmailSender(configurations configs.EmailConfigs) *EmailSender {
	fmt.Println("Configured Email Service!")
	return &EmailSender{configurations: configurations}
}
