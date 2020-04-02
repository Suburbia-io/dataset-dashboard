package application

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

type middleware func(next http.Handler) http.Handler

func middlewareGroup(mdlList ...middleware) middleware {
	return func(handler http.Handler) http.Handler {
		for i := len(mdlList); i > 0; i -= 1 {
			handler = mdlList[i-1](handler)
		}
		return handler
	}
}

func (app *App) secureHeadersMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *App) logMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO: %s - %s %s %s %s", r.RemoteAddr, r.Proto, r.Host, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *App) catchPanicMdl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErr(w, r, errors.Unexpected.WithErr(fmt.Errorf("%s", err)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
