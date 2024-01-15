package service

import (
	"fmt"
	"log/slog"
	"net/smtp"
)

type mailService struct {
	Username   string
	Password   string
	Origin     string
	Port       string
	Encryption string
	Logger     *slog.Logger
}

type MailConfig struct {
	Username   string
	Password   string
	Origin     string
	Port       string
	Encryption string
	Logger     *slog.Logger
}

type MailService interface {
	SendResetEmail(email string, token string) error
}

func NewMailService(c *MailConfig) MailService {
	return &mailService{
		Username:   c.Username,
		Password:   c.Password,
		Origin:     c.Origin,
		Port:       c.Port,
		Encryption: c.Encryption,
		Logger:     c.Logger,
	}
}

// SendResetMail sends a password reset email with the given reset token
func (s *mailService) SendResetEmail(email string, token string) error {
	msg := "From: " + s.Username + "\n" +
		"To: " + email + "\n" +
		"Subject: Reset Email\n\n" +
		fmt.Sprintf("<a href=\"%s/reset-password/%s\">Reset Password</a>", "http://localhost:8080", token)

	err := smtp.SendMail(s.Origin+":"+s.Port,
		smtp.CRAMMD5Auth(s.Username, s.Password),
		s.Username, []string{email}, []byte(msg))

	return err
}
