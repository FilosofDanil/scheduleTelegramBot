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
	id             uuid.UUID
}

func (e EmailSender) SendMail() error {
	// Sender data.
	from := e.configurations.Sender
	password := e.configurations.Password
	id := uuid.New()
	e.id = id
	// Receiver email address.
	to := []string{
		e.configurations.Receiver,
	}

	// smtp server configuration.
	smtpHost := e.configurations.SmtpHost
	smtpPort := e.configurations.SmtpPort

	// Message.
	message := []byte("Test message")

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

func (e EmailSender) CheckUUID(idFromMessage string) error {
	if e.id.String() != idFromMessage {
		return errors.New("incorrect uuid")
	}
	return nil
}

func NewEmailSender(configurations configs.EmailConfigs) *EmailSender {
	fmt.Println("Configured Email Service!")
	return &EmailSender{configurations: configurations}
}
