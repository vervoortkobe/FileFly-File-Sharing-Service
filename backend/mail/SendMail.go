package mail

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type MailConfig struct {
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

func SendEmail(config MailConfig, recipientName, recipientAddress, subject, htmlBody string) error {
	if config.SMTPServer == "" || config.SMTPPort == 0 || config.FromEmail == "" {
		log.Println("SMTP server, port, or fromEmail not configured. Email not sent.")
		return fmt.Errorf("mail configuration incomplete")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(config.FromEmail, config.FromName))
	m.SetAddressHeader("To", recipientAddress, recipientName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	var d *gomail.Dialer
	if config.SMTPUser != "" && config.SMTPPassword != "" {
		d = gomail.NewDialer(config.SMTPServer, config.SMTPPort, config.SMTPUser, config.SMTPPassword)
	} else {
		d = &gomail.Dialer{Host: config.SMTPServer, Port: config.SMTPPort}
		log.Println("Attempting to send email without explicit SMTP username/password authentication.")
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email to %q: %v", recipientAddress, err)
		return err
	}

	log.Printf("Email sent to %q with subject %q", recipientAddress, subject)
	return nil
}
