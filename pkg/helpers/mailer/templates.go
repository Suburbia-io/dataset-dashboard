package mailer

import (
	"html/template"
	"strings"
)

// TODO Replace with HTML emails.
var mailTemplates = map[string]*template.Template{
	"customer-user-login": template.Must(
		template.New("customer-user-login").Parse(`Login to Suburbia

Please click the link to login to your account on Suburbia.io.

{{.LoginLink}}

The login link will expire in {{.ExpirationMins}} minutes.

--

Suburbia
Alternative Data
Better Decisions`)),
}

func SendMailTemplate(mailer Mailer, tmplName string, ctx interface{}, subject string, sendEmail bool, to ...string) (err error) {
	var sb strings.Builder
	mailTmpl := mailTemplates[tmplName]
	err = mailTmpl.Execute(&sb, ctx)
	if err != nil {
		return err
	}
	emailBody := sb.String()

	mail := Mail{
		Subject: subject,
		Body:    emailBody,
	}

	if sendEmail {
		go mailer.Send(mail, to...)
	}
	return nil
}
