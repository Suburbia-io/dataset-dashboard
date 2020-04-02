package application

import (
	"net/http"

	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminFingerprintList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Args struct {
			DatasetID string
		}
		Session shared.SharedSession
	}{Session: s}
	err = decoder.Decode(&ctx.Args, r.Form)
	if err != nil {
		return err
	}

	err = app.Views.ExecuteSST(w, "admin-fingerprint-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}
