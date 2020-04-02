package application

import (
	"encoding/json"
	"net/http"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminTagList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		ListArgs database.TagListArgs
		Tags     []database.Tag
		Dataset  database.Dataset
		TagType  database.TagType
		TagsJSON string
		Session  shared.SharedSession
	}{Session: s}
	err = decoder.Decode(&ctx.ListArgs, r.Form)
	if err != nil {
		return err
	}

	datasetID := r.FormValue("DatasetID")
	ctx.Dataset, err = app.DBAL.DatasetGet(datasetID)
	if err != nil {
		return err
	}

	tagTypeID := r.FormValue("TagTypeID")
	ctx.TagType, err = app.DBAL.TagTypeGet(datasetID, tagTypeID)
	if err != nil {
		return err
	}

	ctx.Tags, err = app.DBAL.TagList(datasetID, tagTypeID, ctx.ListArgs)
	if err != nil {
		return err
	}

	tagsJSON, err := json.Marshal(ctx.Tags)
	if err != nil {
		return err
	}
	ctx.TagsJSON = string(tagsJSON)

	err = app.Views.ExecuteSST(w, "admin-tag-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminTagForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		Dataset database.Dataset
		TagType database.TagType
		Tag     database.Tag
	}{Session: s}
	err = decoder.Decode(&ctx.Tag, r.Form)
	if err != nil {
		return err
	}

	ctx.Dataset, err = app.DBAL.DatasetGet(ctx.Tag.DatasetID)
	if err != nil {
		return err
	}

	ctx.TagType, err = app.DBAL.TagTypeGet(ctx.Tag.DatasetID, ctx.Tag.TagTypeID)
	if err != nil {
		return err
	}

	if ctx.Tag.TagID != "" {
		ctx.Tag, err = app.DBAL.TagGet(ctx.Tag.DatasetID, ctx.Tag.TagTypeID, ctx.Tag.TagID)
		if err != nil {
			return err
		}
	}

	err = app.Views.ExecuteSST(w, "admin-tag-form.gohtml", ctx)
	if err != nil {
		app.serverErr(w, r, err)
	}

	return nil
}

func (app *App) AdminTagSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		Tag     database.Tag
	}{Session: s}
	err = decoder.Decode(&ctx.Tag, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.TagUpsert(&ctx.Tag); err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Tags.Table(), ctx.Tag.TagID, "upsert", ctx.Tag)

	http.Redirect(w, r, "/admin/datasets/"+ctx.Tag.DatasetID+`/tag-types/`+ctx.Tag.TagTypeID+`/tags/`, http.StatusSeeOther)

	return nil
}

func (app *App) AdminTagDelete(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session shared.SharedSession
		Tag     database.Tag
	}{Session: s}
	err = decoder.Decode(&ctx.Tag, r.Form)
	if err != nil {
		return err
	}

	if err = app.DBAL.TagDelete(ctx.Tag.DatasetID, ctx.Tag.TagTypeID, ctx.Tag.TagID); err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Tags.Table(), ctx.Tag.TagID, "delete", ctx.Tag)

	http.Redirect(w, r, "/admin/datasets/"+ctx.Tag.DatasetID+`/tag-types/`+ctx.Tag.TagTypeID+`/tags/`, http.StatusSeeOther)
	return nil
}

func (app *App) AdminUpdateTagCountsFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	dataset, err := app.DBAL.DatasetGet(r.FormValue("DatasetID"))
	if err != nil {
		return err
	}

	tagType, err := app.DBAL.TagTypeGet(dataset.DatasetID, r.FormValue("TagTypeID"))
	if err != nil {
		return err
	}

	go app.DBAL.UpdateTagCounts(dataset.DatasetID, tagType.TagTypeID)

	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Tags.Table(), dataset.DatasetID, "updateTagCounts", "")
	http.Redirect(w, r, "/admin/datasets/"+dataset.DatasetID+"/tag-types/"+tagType.TagTypeID+"/tags/", http.StatusSeeOther)
	return nil
}
