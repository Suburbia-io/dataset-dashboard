package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/mailer"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminAuth(w http.ResponseWriter, r *http.Request, s shared.SharedSession) {
	ctx := struct {
		HasError bool
		Session  shared.SharedSession
	}{
		Session: s,
	}

	ctx.HasError = r.Form.Get("error") == "1"

	err := app.Views.ExecuteSST(w, "admin-auth.gohtml", ctx)
	if err != nil {
		app.serverErr(w, r, err)
	}
}

func (app *App) secureAuthFail(w http.ResponseWriter, r *http.Request, password string) {
	// Just check against a random hash anyway, even though
	// we know this admin does not exist, to make the request slower
	// and therefore don't give away if this admin actually exists.
	crypto.CheckPasswordHash("$2y$14$Ga/NNTftkPZPL38f10816eoZLw6nloBAYbjT3K0Z86yJRQMXqLz96", password)
	http.Redirect(w, r, ".?error=1", http.StatusSeeOther)
}

func (app *App) AdminAuthSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) {
	args := struct {
		HasError bool
		Auth     struct {
			Email    string `schema:"email"`
			Password string `schema:"password"`
		}
		Session shared.SharedSession
	}{
		Session: s,
	}

	args.HasError = r.Form.Get("error") == "1"

	err := decoder.Decode(&args.Auth, r.Form)
	if err != nil {
		// Use a fixed random password because we haven't decoded anything.
		app.secureAuthFail(w, r, "ec636c2863a4cb79bff586b33f09a3b8")
		defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), "", "adminLoginFail", "")
		return
	}

	user, err := app.DBAL.UserGetByEmail(args.Auth.Email)
	if err != nil {
		app.secureAuthFail(w, r, args.Auth.Password)
		defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), "", "adminLoginFail", args.Auth.Email)
		return
	}

	if !crypto.CheckPasswordHash(user.Hash, args.Auth.Password) {
		defer app.DBAL.AuditTrailBySystemInsertAsync(
			tables.Users.Table(), user.UserID, "adminLoginPasswordFail", args.Auth.Email)
		http.Redirect(w, r, ".?error=1", http.StatusSeeOther)
		return
	}

	if user.CustomerID != shared.SuburbiaInternalCustomerUUID {
		defer app.DBAL.AuditTrailBySystemInsertAsync(
			tables.Users.Table(), user.UserID, "adminLoginNonAdminFail", args.Auth.Email)
		http.Redirect(w, r, ".?error=1", http.StatusSeeOther)
		return
	}

	emailToken, browserToken, err := app.DBAL.AuthCreateAdminToken(
		user.UserID,
		app.Config.AdminAuthTokenLifetimeSec*time.Second,
	)
	if err != nil {
		app.serverErr(w, r, err)
		return
	}

	mail := mailer.Mail{
		Subject: "Suburbia Verify",
		Body:    emailToken,
	}
	if err := app.Mailer.Send(mail, user.Email); err != nil {
		app.serverErr(w, r, errors.UnexpectedError(err, "Failed sending admin auth email"))
		return
	}
	defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), user.UserID, "adminLoginSuccess", args.Auth.Email)
	http.Redirect(w, r, fmt.Sprintf("/admin/auth/two-factor/?browserToken=%s", browserToken), http.StatusSeeOther)
}

func (app *App) AdminTwoFactor(w http.ResponseWriter, r *http.Request, s shared.SharedSession) {
	args := struct {
		Auth struct {
			BrowserToken string `schema:"browserToken"`
			EmailToken   string `schema:"emailToken"`
		}
		Session shared.SharedSession
	}{
		Session: s,
	}
	err := decoder.Decode(&args.Auth, r.Form)
	if err != nil {
		http.Redirect(w, r, "/admin/auth/?error=1", http.StatusSeeOther)
		return
	}

	if args.Auth.BrowserToken == "" {
		http.Redirect(w, r, "/admin/auth/?error=1", http.StatusSeeOther)
		return
	}

	err = app.Views.ExecuteSST(w, "admin-auth-twoFactor.gohtml", args)
	if err != nil {
		app.serverErr(w, r, err)
	}
}

func (app *App) AdminTwoFactorSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) {
	args := struct {
		Auth struct {
			BrowserToken string `schema:"browserToken"`
			EmailToken   string `schema:"emailToken"`
		}
		Session shared.SharedSession
	}{
		Session: s,
	}
	err := decoder.Decode(&args.Auth, r.Form)
	if err != nil {
		http.Redirect(w, r, "/admin/auth/?error=1", http.StatusSeeOther)
		return
	}

	if args.Auth.BrowserToken == "" {
		http.Redirect(w, r, "/admin/auth/?error=1", http.StatusSeeOther)
		return
	}

	userID, err := app.DBAL.AuthAuthenticateAdminToken(args.Auth.EmailToken, args.Auth.BrowserToken, shared.SuburbiaInternalCustomerUUID)
	if err != nil {
		defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), "", "adminTwoFactorFail", args.Auth.EmailToken)
		http.Redirect(w, r, "/admin/auth/?error=1", http.StatusSeeOther)
		return
	}

	session, err := app.DBAL.SessionCreate(userID, app.Config.SessionLifetimeSec*time.Second)
	if err != nil {
		app.respondApi(w, r, nil, err)
		return
	}

	if session.User.CustomerID != shared.SuburbiaInternalCustomerUUID {
		defer app.DBAL.AuditTrailBySystemInsertAsync(
			tables.Users.Table(), session.User.UserID, "adminTwoFactorNonAdminFail", session.User.Email)
		http.Redirect(w, r, "/admin/auth/?error=1", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     shared.SessionCookieName,
		Domain:   app.Config.SessionCookieDomain,
		Secure:   app.Config.SessionCookieSecure,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   0,
		Value:    session.Token,
	})

	defer app.DBAL.AuditTrailBySystemInsertAsync(
		tables.Users.Table(), userID, "adminTwoFactorSuccess", args.Auth.EmailToken)

	if session.User.IsRoleAdmin {
		http.Redirect(w, r, "/admin/audit-trail/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/admin/datasets/", http.StatusSeeOther)
	}
}
