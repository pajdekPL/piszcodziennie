package main

import (
	"net/smtp"
)

func (cfg *Config) sendEmail(to string, link string) error {
	auth := smtp.PlainAuth("", cfg.SMTPEnvs.Email, cfg.SMTPEnvs.SMTP_PWD, cfg.SMTPEnvs.SMTP_SERVER)
	serverAndPort := cfg.SMTPEnvs.SMTP_SERVER + ":" + cfg.SMTPEnvs.SMTP_PORT
	msg := []byte("Subject: Verify Your Email\n\nClick the link to verify: " + link)
	err := smtp.SendMail(serverAndPort, auth, cfg.SMTPEnvs.Email, []string{to}, msg)
	return err
}
