package application

import (
	"net/http"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminCorporationTypeForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session         shared.SharedSession
		Dataset         database.Dataset
		CorporationType database.CorporationType
	}{Session: s}
	err = decoder.Decode(&ctx.CorporationType, r.Form)
	if err != nil {
		return err
	}

	ctx.Dataset, err = app.DBAL.DatasetGet(ctx.CorporationType.DatasetID)
	if err != nil {
		return err
	}

	if ctx.CorporationType.CorporationTypeID != "" {
		ctx.CorporationType, err = app.DBAL.CorporationTypeGet(ctx.CorporationType.CorporationTypeID)
		if err != nil {
			app.serverErr(w, r, err)
		}
	}

	err = app.Views.ExecuteSST(w, "admin-corporationtype-form.gohtml", ctx)
	if err != nil {
		app.serverErr(w, r, err)
	}

	return nil
}

func (app *App) AdminCorporationTypeSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session         shared.SharedSession
		CorporationType database.CorporationType
	}{Session: s}
	err = decoder.Decode(&ctx.CorporationType, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.CorporationTypeUpsert(&ctx.CorporationType); err != nil {
		return err
	}
	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.CorporationTypes.Table(), ctx.CorporationType.CorporationTypeID, "upsert", ctx.CorporationType)

	http.Redirect(w, r, "/admin/datasets/"+ctx.CorporationType.DatasetID+`/`, http.StatusSeeOther)

	return nil
}

func (app *App) AdminCorporationTypeDelete(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session         shared.SharedSession
		CorporationType database.CorporationType
	}{Session: s}
	err = decoder.Decode(&ctx.CorporationType, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.CorporationTypeDelete(ctx.CorporationType.CorporationTypeID); err != nil {
		return err
	}
	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.CorporationTypes.Table(), ctx.CorporationType.CorporationTypeID, "delete", ctx.CorporationType)

	http.Redirect(w, r, "/admin/datasets/"+ctx.CorporationType.DatasetID+`/`, http.StatusSeeOther)
	return nil
}
