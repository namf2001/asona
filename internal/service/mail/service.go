package mail

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strconv"

	"asona/config"
)

// Service defines the mail sending contract.
type Service interface {
	SendMail(to []string, subject, body string) error
	SendWithTemplate(to []string, subject, templatePath string, data interface{}) error
}

type service struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	fromEmail    string
}

func New() Service {
	cfg := config.GetConfig()
	return &service{
		smtpHost:     cfg.MailSMTPHost,
		smtpPort:     cfg.MailSMTPPort,
		smtpUser:     cfg.MailSMTPUser,
		smtpPassword: cfg.MailSMTPPassword,
		fromEmail:    cfg.MailEmailFrom,
	}
}

func (s *service) SendMail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte("Subject: " + subject + "\n" + mime + "\n" + body)
	addr := s.smtpHost + ":" + strconv.Itoa(s.smtpPort)

	return smtp.SendMail(addr, auth, s.fromEmail, to, msg)
}

func (s *service) SendWithTemplate(to []string, subject, templatePath string, data interface{}) error {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	return s.SendMail(to, subject, buf.String())
}
