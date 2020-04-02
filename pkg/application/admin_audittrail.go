package application

import (
	"net/http"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminAuditTrailList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		ListArgs    database.AuditTrailListArgs
		AuditTrails []database.AuditTrail
		Session     shared.SharedSession
	}{Session: s}
	err = decoder.Decode(&ctx.ListArgs, r.Form)
	if err != nil {
		return err
	}

	// No pagination at the moment.
	ctx.ListArgs.Limit = 1000
	ctx.ListArgs.Offset = 0

	ctx.AuditTrails, err = app.DBAL.AuditTrailList(ctx.ListArgs)
	if err != nil {
		return err
	}

	err = app.Views.ExecuteSST(w, "admin-audittrail-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}
