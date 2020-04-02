package application

import (
	"net/http"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminDatasetList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		ListArgs database.DatasetListArgs
		Datasets []database.Dataset
		Session  shared.SharedSession
	}{Session: s}
	err = decoder.Decode(&ctx.ListArgs, r.Form)
	if err != nil {
		return err
	}

	// Don't allow labelers to see archived and non-manageable datasets.
	if s.User.IsRoleLabeler && !s.User.IsRoleAdmin {
		archived := false
		ctx.ListArgs.Archived = &archived
		manageable := true
		ctx.ListArgs.Manageable = &manageable
	}

	// No pagination at the moment.
	ctx.ListArgs.Limit = 100000
	ctx.ListArgs.Offset = 0

	ctx.Datasets, err = app.DBAL.DatasetList(ctx.ListArgs)
	if err != nil {
		return err
	}

	err = app.Views.ExecuteSST(w, "admin-dataset-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminDatasetForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		Dataset struct {
			database.Dataset
			Archived bool
		}
		TagTypes         []database.TagType
		CorporationTypes []database.CorporationType
		Corporations     []database.Corporation
	}{Session: s}
	err = decoder.Decode(&ctx.Dataset, r.Form)
	if err != nil {
		return err
	}

	if ctx.Dataset.DatasetID != "" {
		ctx.Dataset.Dataset, err = app.DBAL.DatasetGet(ctx.Dataset.DatasetID)
		if err != nil {
			return err
		}
		ctx.TagTypes, err = app.DBAL.TagTypeList(database.TagTypeListArgs{DatasetID: ctx.Dataset.DatasetID})
		if err != nil {
			return err
		}
		ctx.CorporationTypes, err = app.DBAL.CorporationTypeList(database.CorporationTypeListArgs{DatasetID: ctx.Dataset.DatasetID})
		if err != nil {
			return err
		}
		ctx.Corporations, err = app.DBAL.CorporationList(database.CorporationListArgs{DatasetID: ctx.Dataset.DatasetID})
		if err != nil {
			return err
		}
	}

	err = app.Views.ExecuteSST(w, "admin-dataset-form.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminDatasetFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		Dataset struct {
			database.Dataset
			Archived bool
		}
	}{Session: s}
	err = decoder.Decode(&ctx.Dataset, r.Form)
	if err != nil {
		return err
	}

	if (ctx.Dataset.ArchivedAt != nil) != ctx.Dataset.Archived {
		if ctx.Dataset.Archived {
			t := time.Now()
			ctx.Dataset.ArchivedAt = &t
		} else {
			ctx.Dataset.ArchivedAt = nil
		}
	}

	if err = app.DBAL.DatasetUpsert(&ctx.Dataset.Dataset); err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Datasets.Table(), ctx.Dataset.DatasetID, "upsert", ctx.Dataset)

	http.Redirect(w, r, "/admin/datasets/", http.StatusSeeOther)

	return nil
}
