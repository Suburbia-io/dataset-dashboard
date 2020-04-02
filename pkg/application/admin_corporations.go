package application

import (
	"net/http"

	"github.com/Suburbia-io/dashboard/shared"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
)

func (app *App) AdminCorporationForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session     shared.SharedSession
		Dataset     database.Dataset
		Corporation database.Corporation
	}{Session: s}
	err = decoder.Decode(&ctx.Corporation, r.Form)
	if err != nil {
		return err
	}

	ctx.Dataset, err = app.DBAL.DatasetGet(ctx.Corporation.DatasetID)
	if err != nil {
		return err
	}

	if ctx.Corporation.CorporationID != "" {
		ctx.Corporation, err = app.DBAL.CorporationGet(ctx.Corporation.CorporationID)
		if err != nil {
			app.serverErr(w, r, err)
		}
	}

	err = app.Views.ExecuteSST(w, "admin-corporation-form.gohtml", ctx)
	if err != nil {
		app.serverErr(w, r, err)
	}

	return nil
}

func (app *App) AdminCorporationSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session     shared.SharedSession
		Corporation database.Corporation
	}{Session: s}
	err = decoder.Decode(&ctx.Corporation, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.CorporationUpsert(&ctx.Corporation); err != nil {
		return err
	}
	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Corporations.Table(), ctx.Corporation.CorporationID, "upsert", ctx.Corporation)

	http.Redirect(w, r, "/admin/datasets/"+ctx.Corporation.DatasetID+`/`, http.StatusSeeOther)

	return nil
}

func (app *App) AdminCorporationDelete(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session     shared.SharedSession
		Corporation database.Corporation
	}{Session: s}
	err = decoder.Decode(&ctx.Corporation, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.CorporationDelete(ctx.Corporation.CorporationID); err != nil {
		return err
	}
	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Corporations.Table(), ctx.Corporation.CorporationID, "delete", ctx.Corporation)

	http.Redirect(w, r, "/admin/datasets/"+ctx.Corporation.DatasetID+`/`, http.StatusSeeOther)
	return nil
}
