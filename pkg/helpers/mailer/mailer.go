package mailer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

type MailCatcher struct {
	SentMail []CaughtMail
}

type CaughtMail struct {
	Mail Mail
	To   []string
}

func NewMailCatcher() *MailCatcher {
	return &MailCatcher{
		SentMail: make([]CaughtMail, 0),
	}
}

func (mailer *MailCatcher) Send(mail Mail, to ...string) error {
	mailer.SentMail = append(mailer.SentMail, CaughtMail{mail, to})
	return nil
}

func (mailer *MailCatcher) Empty() {
	mailer.SentMail = []CaughtMail{}
}

type LogMailer struct{}

func NewLogMailer() LogMailer {
	return LogMailer{}
}

func (mailer LogMailer) Send(mail Mail, to ...string) error {
	fmt.Printf("---------------------------------------\n")
	log.Printf("Email:\n")
	fmt.Printf("to:      %s\n", to)
	fmt.Printf("subject: %s\n", mail.Subject)
	fmt.Printf("body:    %s\n", mail.Body)
	fmt.Printf("---------------------------------------\n")

	return nil
}

type Config struct {
	Driver        string
	SenderName    string
	SenderEmail   string
	RedirectEmail string

	MailgunConfig
}

type Mailer interface {
	Send(mail Mail, to ...string) error
}

type Mail struct {
	Subject string
	Body    string
}

func NewMailer(conf Config) Mailer {
	switch conf.Driver {
	case "log":
		return NewLogMailer()
	case "mailgun":
		return NewMailgunMailer(conf)
	default:
		panic("unknown mailer driver")
	}
}

type MailgunConfig struct {
	MailgunDomain string
	MailgunAPIKey string
}

type MailgunMailer struct {
	redirectEmail string
	senderEmail   string
	senderName    string
	mg            mailgun.Mailgun
}

func NewMailgunMailer(conf Config) MailgunMailer {
	return MailgunMailer{
		mg:            mailgun.NewMailgun(conf.MailgunDomain, conf.MailgunAPIKey),
		redirectEmail: conf.RedirectEmail,
		senderEmail:   conf.SenderEmail,
		senderName:    conf.SenderName,
	}
}

func (mailer MailgunMailer) Send(mail Mail, to ...string) error {

	//if redirectEmail, then don't really send to recipients
	if mailer.redirectEmail != "" {
		to = []string{mailer.redirectEmail}
	}

	message := mailer.mg.NewMessage(
		fmt.Sprintf("%s <%s>", mailer.senderName, mailer.senderEmail),
		mail.Subject,
		mail.Body,
		to...,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mailer.mg.Send(ctx, message)
	return err
}
