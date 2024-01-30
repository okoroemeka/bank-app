package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

type EmailSender interface {
	SendEmail(subject, content string, to, cc, bcc, attachFiles []string) error
}

type GmailSender struct {
	name, FromEmailAddress, FromEmailPassword string
}

const (
	smtpAuthServer = "smtp.gmail.com"
	smtpServerAddr = "smtp.gmail.com:587"
)

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{name: name, FromEmailAddress: fromEmailAddress, FromEmailPassword: fromEmailPassword}
}

func (sender *GmailSender) SendEmail(subject, content string, to, cc, bcc, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.FromEmailAddress)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.Subject = subject
	e.HTML = []byte(content)

	for _, attachFile := range attachFiles {
		if _, err := e.AttachFile(attachFile); err != nil {
			return fmt.Errorf("cannot attach file: %w", err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.FromEmailAddress, sender.FromEmailPassword, smtpAuthServer)

	return e.Send(smtpServerAddr, smtpAuth)
}
