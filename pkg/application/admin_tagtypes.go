package application

import (
	"net/http"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminTagTypeForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		Dataset database.Dataset
		TagType database.TagType
	}{Session: s}
	err = decoder.Decode(&ctx.TagType, r.Form)
	if err != nil {
		return err
	}

	ctx.Dataset, err = app.DBAL.DatasetGet(ctx.TagType.DatasetID)
	if err != nil {
		return err
	}

	if ctx.TagType.TagTypeID != "" {
		ctx.TagType, err = app.DBAL.TagTypeGet(ctx.TagType.DatasetID, ctx.TagType.TagTypeID)
		if err != nil {
			app.serverErr(w, r, err)
		}
	}

	err = app.Views.ExecuteSST(w, "admin-tagtype-form.gohtml", ctx)
	if err != nil {
		app.serverErr(w, r, err)
	}

	return nil
}

func (app *App) AdminTagTypeSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		TagType database.TagType
	}{Session: s}
	err = decoder.Decode(&ctx.TagType, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.TagTypeUpsert(&ctx.TagType); err != nil {
		return err
	}
	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.TagTypes.Table(), ctx.TagType.TagTypeID, "upsert", ctx.TagType)

	http.Redirect(w, r, "/admin/datasets/"+ctx.TagType.DatasetID+`/`, http.StatusSeeOther)

	return nil
}

func (app *App) AdminTagTypeDelete(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		TagType database.TagType
	}{Session: s}
	err = decoder.Decode(&ctx.TagType, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.TagTypeDelete(ctx.TagType.DatasetID, ctx.TagType.TagTypeID); err != nil {
		return err
	}
	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.TagTypes.Table(), ctx.TagType.TagTypeID, "delete", ctx.TagType)

	http.Redirect(w, r, "/admin/datasets/"+ctx.TagType.DatasetID+`/`, http.StatusSeeOther)
	return nil
}
