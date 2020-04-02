package application

import (
	"net/http"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"

	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) Logout(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	redirectUrl := "/login/"
	if s.UserIsAuthenticated() {
		if s.User.IsRoleAdmin || s.User.IsRoleSuperAdmin || s.User.IsRoleLabeler {
			redirectUrl = "/admin/auth/"
		}
		err = app.DBAL.SessionDelete(s.Token)
		if err != nil {
			return err
		}
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
	return nil
}

func (app *App) loginRedirectAuthenticatedUser(w http.ResponseWriter, r *http.Request, s shared.SharedSession) {
	if s.User.IsRoleAdmin {
		http.Redirect(w, r, "/admin/audit-trail/", http.StatusSeeOther)
	} else if s.User.IsRoleLabeler {
		http.Redirect(w, r, "/admin/datasets/", http.StatusSeeOther)
	} else if s.User.IsRoleCustomerUser {
		http.Redirect(w, r, "/dashboard/", http.StatusSeeOther)
	} else {
		_ = app.Logout(w, r, s)
	}
}

func (app *App) authenticateCustomerUser(w http.ResponseWriter, userID, loginToken string) (database.Session, error) {
	s, err := app.DBAL.SessionCreate(userID, app.Config.SessionLifetimeSec*time.Second)
	if err != nil {
		defer app.DBAL.AuditTrailBySystemInsertAsync(
			tables.Users.Table(), userID, "customerUserLoginSessionCreateFail", loginToken)
		return database.Session{}, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     shared.SessionCookieName,
		Domain:   app.Config.SessionCookieDomain,
		Secure:   app.Config.SessionCookieSecure,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   0,
		Value:    s.Token,
	})

	defer app.DBAL.AuditTrailBySystemInsertAsync(
		tables.Users.Table(), userID, "customerUserLoginSuccess", loginToken)
	return s, nil
}
