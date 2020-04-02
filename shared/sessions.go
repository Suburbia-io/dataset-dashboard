package shared

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/gorilla/csrf"
)

type SharedSession struct {
	database.Session
	AuthMethod *string
	r          *http.Request
}

func (s SharedSession) CSRFField() template.HTML {
	return csrf.TemplateField(s.r)
}

func (s SharedSession) CSRFToken() string {
	return csrf.Token(s.r)
}

func (s SharedSession) UserIsAuthenticated() bool {
	return s.UserID != "" && s.User.UserID != "" && s.AuthMethod != nil
}

func GetSessionFromCookie(r *http.Request, db *database.DBAL) (s SharedSession) {
	s.r = r
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return
	}
	s.Session, _ = db.SessionGet(cookie.Value)
	authMethodCookie := "C"
	s.AuthMethod = &authMethodCookie
	return
}

func GetSessionFromBearer(r *http.Request, db *database.DBAL) (s SharedSession) {
	s.r = r
	authHeaderValue := r.Header.Get("Authorization")
	if authHeaderValue == "" {
		return
	}

	if !strings.HasPrefix(authHeaderValue, "Bearer ") {
		return
	}

	apiKey := authHeaderValue[7:]

	user, err := db.UserGetByAPIKey(apiKey)
	if err != nil {
		return
	}

	// API access is only for internal use.
	if user.CustomerID != InternalCustomerUUID {
		return
	}

	s.UserID = user.UserID
	s.User = user
	s.Token = apiKey
	authMethodCookie := "B"
	s.AuthMethod = &authMethodCookie
	return
}
