package application

import (
	"errors"
	"net/http"
	"net/url"

	errors2 "github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/shared"
	"github.com/gorilla/mux"
)

type RoleFlag uint32

const (
	// PUBLIC is for publicly accessible pages. This role should not be combined with other roles.
	Public RoleFlag = 1 << iota
	RoleCustomerUser
	RoleAdmin
	RoleSuperAdmin
	RoleLabeler
)

func hasFlag(requiredRoles, roleFlags RoleFlag) bool {
	return requiredRoles&roleFlags != 0
}

func checkAuthentication(r *http.Request, app *App) shared.SharedSession {
	session := shared.GetSessionFromCookie(r, app.DBAL)
	if !session.UserIsAuthenticated() {
		session = shared.GetSessionFromBearer(r, app.DBAL)
	}
	return session
}

func checkAuthenticatedRoles(requiredRoles RoleFlag, session shared.SharedSession) (err error) {
	var userRoles RoleFlag = 0

	if session.User.IsRoleCustomerUser {
		userRoles = userRoles | RoleCustomerUser
	}

	if session.User.IsRoleAdmin {
		userRoles = userRoles | RoleAdmin
	}

	if session.User.IsRoleSuperAdmin {
		userRoles = userRoles | RoleSuperAdmin
	}

	if session.User.IsRoleLabeler {
		userRoles = userRoles | RoleLabeler
	}

	if !hasFlag(requiredRoles, userRoles) {
		return errors2.HttpForbidden
	}

	return
}

func (app *App) pageWrapper(
	router *mux.Router,
	pattern string,
	name string,
	method string,
	requiredRoles RoleFlag,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		s shared.SharedSession,
	) (err error),
) {
	router.NewRoute().
		Path(pattern).
		Methods(method).
		Name(name).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if hasFlag(requiredRoles, Public) && requiredRoles != Public {
				w.WriteHeader(500)
				_, _ = w.Write([]byte("public role cannot be combined with other roles"))
				return
			}

			session := checkAuthentication(r, app)

			if requiredRoles != Public && !session.UserIsAuthenticated() {
				if hasFlag(requiredRoles, RoleAdmin|RoleSuperAdmin|RoleLabeler) {
					http.Redirect(w, r, "/admin/auth/", http.StatusSeeOther)
					return
				} else if hasFlag(requiredRoles, RoleCustomerUser) {
					http.Redirect(w, r, "/login/", http.StatusSeeOther)
					return
				}
			}

			if requiredRoles != Public {
				err := checkAuthenticatedRoles(requiredRoles, session)
				if err != nil {
					http.NotFound(w, r)
					return
				}
			}

			// At this point the roles have been checked for the users that are authenticated.

			_ = r.ParseForm()
			for key, value := range mux.Vars(r) {
				r.Form.Set(key, value)
			}

			err := handler(w, r, session)
			if session.UserIsAuthenticated() {
				if err == nil {
					_ = app.DBAL.UserTouchLastActiveAt(session.User.CustomerID, session.User.UserID)
				}

				// Basic error handling is only for the admin dashboard pages, not the customer user pages.
				if hasFlag(requiredRoles, RoleAdmin|RoleSuperAdmin|RoleLabeler) && err != nil {
					if errors.Unwrap(err) != nil {
						err = errors.Unwrap(err)
					}
					ctx := struct {
						Session shared.SharedSession
						Error   error
					}{
						Session: session,
						Error:   err,
					}
					tmplErr := app.Views.ExecuteSST(w, "admin-error.gohtml", ctx)
					if tmplErr != nil {
						app.serverErr(w, r, tmplErr)
					}
					return
				}
			}

			if err != nil {
				app.serverErr(w, r, err)
			}
		})
}

func (app *App) adminPage(
	router *mux.Router,
	pattern string,
	name string,
	method string,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		s shared.SharedSession,
	) (err error),
) {
	app.pageWrapper(router, pattern, name, method, RoleAdmin, handler)
}

// TODO Use pageWrapper() to implement adminAuthPage() if possible.
func (app *App) adminAuthPage(
	router *mux.Router,
	pattern string,
	name string,
	method string,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		s shared.SharedSession,
	),
) {
	router.NewRoute().
		Path(pattern).
		Methods(method).
		Name(name).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			for key, value := range mux.Vars(r) {
				r.Form.Set(key, value)
			}

			// if user is already logged in, redirect
			session := shared.GetSessionFromCookie(r, app.DBAL)
			if session.UserIsAuthenticated() {
				app.loginRedirectAuthenticatedUser(w, r, session)
				return
			}

			handler(w, r, session)
		})
}

func getOrigin(r *http.Request) (origin string) {
	if origin = r.Header.Get("Origin"); len(origin) > 0 {
		return origin
	}

	if origin = r.Header.Get("Referer"); len(origin) == 0 {
		return ""
	}

	if u, err := url.Parse(origin); err != nil {
		return ""
	} else {
		return u.Scheme + "://" + u.Host
	}
}

func addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := getOrigin(r)
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Referer")
}

func addAPIHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("ViewCache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func (app *App) apiWrapper(
	router *mux.Router,
	pattern string,
	name string,
	requiredRoles RoleFlag,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		s shared.SharedSession,
	) (err error),
) {
	router.NewRoute().
		Path(pattern).
		Methods(http.MethodPost).
		Name(name).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if hasFlag(requiredRoles, Public) {
				app.respondApi(w, r, false, errors.New("invalid configuration"))
				return
			}

			session := checkAuthentication(r, app)

			if !session.UserIsAuthenticated() {
				app.respondApi(w, r, false, errors2.HttpNotAuthorized)
				return
			}

			addCORSHeaders(w, r)
			if r.Method == "OPTIONS" {
				return
			}
			addAPIHeaders(w)

			err := checkAuthenticatedRoles(requiredRoles, session)
			if err != nil {
				app.respondApi(w, r, false, err)
				return
			}

			err = handler(w, r, session)
			if err != nil {
				app.respondApi(w, r, false, err)
				return
			}
			_ = app.DBAL.UserTouchLastActiveAt(session.User.CustomerID, session.User.UserID)
		})
}

func (app *App) adminAPI(
	router *mux.Router,
	pattern string,
	name string,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		s shared.SharedSession,
	) (err error),
) {
	app.apiWrapper(router, pattern, name, RoleAdmin, handler)
}
