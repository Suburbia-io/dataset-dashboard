package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/mailer"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) CustomerUserLogin(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	if s.UserIsAuthenticated() {
		app.loginRedirectAuthenticatedUser(w, r, s)
		return nil
	}

	loginToken := r.Form.Get("loginToken")

	if loginToken != "" {
		user, err := app.DBAL.UserGetByValidLoginToken(loginToken)
		if err != nil {
			defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), "", "customerUserLoginTokenFail", loginToken)
			http.Redirect(w, r, "/login/?error=1", http.StatusSeeOther)
			return nil
		}

		// This shouldn't happen because we shouldn't send the loginToken for users that don't have role CustomerUser but
		// it's a nice safety check.
		if !user.IsRoleCustomerUser {
			defer app.DBAL.AuditTrailBySystemInsertAsync(
				tables.Users.Table(), s.User.UserID, "customerUserLoginRoleFail", user.Email)
			http.Redirect(w, r, "/login/?error=1", http.StatusSeeOther)
			return nil
		}

		s.Session, err = app.authenticateCustomerUser(w, user.UserID, loginToken)
		if err != nil {
			http.Redirect(w, r, "/login/?error=1", http.StatusSeeOther)
			return nil
		}

		go app.DBAL.UserClearLoginToken(user.CustomerID, user.UserID)

		http.Redirect(w, r, "/dashboard/", http.StatusSeeOther)
		return nil
	}

	ctx := struct {
		EmailSent     bool
		HasError      bool
		HasErrorEmail bool
		Session       shared.SharedSession
	}{
		EmailSent:     r.Form.Get("emailSent") == "1",
		HasError:      r.Form.Get("error") == "1",
		HasErrorEmail: r.Form.Get("errorEmail") == "1",
		Session:       s,
	}
	return app.Views.ExecuteSST(w, "customer-user-login.gohtml", ctx)
}

func (app *App) CustomerUserLoginSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	secureFail := func(userID, message, email string) {
		defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), userID, message, email)
		// Use the success flow we don't leak information from our system.
		_ = mailer.SendMailTemplate(app.Mailer, "customer-user-login", struct{}{}, "Suburbia Login", false,
			"dummy@example.com")
		http.Redirect(w, r, "/login/?emailSent=1", http.StatusSeeOther)
	}

	args := struct {
		Auth struct {
			Email    string `schema:"email"`
			Password string `schema:"password"`
		}
		Session shared.SharedSession
	}{
		Session: s,
	}

	if args.Session.UserIsAuthenticated() {
		app.loginRedirectAuthenticatedUser(w, r, s)
		return nil
	}

	err = decoder.Decode(&args.Auth, r.Form)
	if err != nil {
		defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), "", "customerUserLoginDecodeFail", "")
		http.Redirect(w, r, "/login/?errorEmail=1", http.StatusSeeOther)
		return nil
	}

	if err = validate.Email(args.Auth.Email); err != nil {
		defer app.DBAL.AuditTrailBySystemInsertAsync(
			tables.Users.Table(), "", "customerUserLoginEmailFail", args.Auth.Email)
		http.Redirect(w, r, "/login/?errorEmail=1", http.StatusSeeOther)
		return nil
	}

	// Honeypot
	if args.Auth.Password != "" {
		secureFail("", "customerUserLoginHoneypotFail", args.Auth.Email)
		fmt.Println("hp")
		return nil
	}

	user, err := app.DBAL.UserGetByEmail(args.Auth.Email)
	if err != nil {
		secureFail("", "customerUserLoginQueryFail", args.Auth.Email)
		return nil
	}

	if user.ArchivedAt != nil || user.Customer.ArchivedAt != nil {
		secureFail(user.UserID, "customerUserLoginArchivedFail", args.Auth.Email)
		return nil
	}

	if !user.IsRoleCustomerUser || user.CustomerID == shared.InternalCustomerUUID {
		secureFail(user.UserID, "customerUserLoginNotAllowedFail", args.Auth.Email)
		return nil
	}

	loginToken, err := app.DBAL.UserCreateLoginToken(user.CustomerID, user.UserID,
		shared.CustomerUserLoginTokenLifetimeSec*time.Second,
	)
	if err != nil {
		defer app.DBAL.AuditTrailBySystemInsertAsync(
			tables.Users.Table(), user.UserID, "customerUserEmailCreateTokenFail", args.Auth.Email)
		return err
	}

	ctx := struct {
		LoginLink      string
		ExpirationMins int
	}{
		LoginLink:      fmt.Sprintf("%s/login/?loginToken=%s", app.Config.ServerHostname, loginToken),
		ExpirationMins: shared.CustomerUserLoginTokenLifetimeSec / 60,
	}

	err = mailer.SendMailTemplate(app.Mailer, "customer-user-login", ctx, "Suburbia Login", true, user.Email)
	if err != nil {
		defer app.DBAL.AuditTrailBySystemInsertAsync(
			tables.Users.Table(), user.UserID, "customerUserEmailSendFail", args.Auth.Email)
		http.Redirect(w, r, "/login/?error=1", http.StatusSeeOther)
		return nil
	}

	defer app.DBAL.AuditTrailBySystemInsertAsync(
		tables.Users.Table(), user.UserID, "customerUserEmailSent", args.Auth.Email)
	http.Redirect(w, r, "/login/?emailSent=1", http.StatusSeeOther)
	return nil
}
