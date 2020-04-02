package application

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminCorpMappings(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		CorpMappings     []database.CorpMapping
		TagTypes         []database.TagType
		TagTypesJSON     template.JS
		CorpMappingsJSON template.JS
		Corporations     []database.Corporation
		CorporationsJSON template.JS
		Dataset          database.Dataset
		CorpType         database.CorporationType
		Session          shared.SharedSession
	}{Session: s}

	datasetID := r.FormValue("DatasetID")
	ctx.Dataset, err = app.DBAL.DatasetGet(datasetID)
	if err != nil {
		return err
	}

	corpTypeID := r.FormValue("CorpTypeID")
	ctx.CorpType, err = app.DBAL.CorporationTypeGet(corpTypeID)
	if err != nil {
		return err
	}

	corpMappings, err := app.DBAL.CorpMappingList()
	if err != nil {
		return err
	}
	ctx.CorpMappings = corpMappings

	corpMappingsJSON, err := json.Marshal(corpMappings)
	if err != nil {
		return err
	}
	ctx.CorpMappingsJSON = template.JS(string(corpMappingsJSON))

	tagTypes, err := app.DBAL.TagTypeList(database.TagTypeListArgs{DatasetID: datasetID})
	if err != nil {
		return err
	}
	ctx.TagTypes = tagTypes

	tagTypesJSON, err := json.Marshal(ctx.TagTypes)
	if err != nil {
		return err
	}
	ctx.TagTypesJSON = template.JS(string(tagTypesJSON))

	corporations, err := app.DBAL.CorporationList(database.CorporationListArgs{DatasetID: datasetID})
	if err != nil {
		return err
	}
	ctx.Corporations = corporations

	corporationsJSON, err := json.Marshal(corporations)
	if err != nil {
		return err
	}
	ctx.CorporationsJSON = template.JS(string(corporationsJSON))

	err = app.Views.ExecuteSST(w, "admin-corpmapping-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

type ListTagsForTagTypeArgs struct {
	DatasetID string `json:"datasetID"`
	TagTypeID string `json:"tagTypeID"`
}

func (app *App) AdminApiListTagsForTagType(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Tags    []database.Tag
		Session shared.SharedSession
		Args    ListTagsForTagTypeArgs
	}{Session: s}
	app.decodeRequest(r, &ctx.Args)

	ctx.Tags, err = app.DBAL.TagList(ctx.Args.DatasetID, ctx.Args.TagTypeID, database.TagListArgs{
		Search:     "",
		IsIncluded: nil,
	})
	if err != nil {
		return err
	}

	app.respondApi(w, r, ctx.Tags, nil)
	return nil
}

func (app *App) AdminApiInsertCorpMapping(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session     shared.SharedSession
		CorpMapping tables.CorpMapping
	}{Session: s}
	app.decodeRequest(r, &ctx.CorpMapping)

	ctx.CorpMapping.CorpMappingID = crypto.NewUUID()
	err = app.DBAL.CorpMappingUpsert(&ctx.CorpMapping)
	if err != nil {
		print(err.Error())
		return err
	}

	app.respondApi(w, r, ctx.CorpMapping, err)
	return nil
}

func (app *App) AdminApiDeleteCorpMappingRule(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		CorpMappingRuleID string `json:"corpMappingRuleID"`
	}{}

	app.decodeRequest(r, &ctx)
	print(ctx.CorpMappingRuleID)
	err = app.DBAL.CorpMappingRuleDelete(ctx.CorpMappingRuleID)
	if err != nil {
		print(err.Error())
		return err
	}

	app.respondApi(w, r, true, err)
	return nil
}

func (app *App) AdminApiInsertCorpMappingRule(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session   shared.SharedSession
		FromDate  string `json:"fromDate"`
		Country   string `json:"country"`
		CorpID    string `json:"corpID"`
		MappingID string `json:"corpMappingID"`
	}{Session: s}
	app.decodeRequest(r, &ctx)

	fromDate, err := time.Parse("2006-01-02", ctx.FromDate)
	if err != nil {
		return err
	}

	rule := tables.CorpMappingRule{
		CorpMappingRuleID: crypto.NewUUID(),
		CorpMappingID:     ctx.MappingID,
		CorpID:            ctx.CorpID,
		InternalNotes:     "",
		ExternalNotes:     "",
		FromDate:          fromDate,
		Country:           ctx.Country,
	}

	err = app.DBAL.CorpMappingRuleUpsert(&rule)
	if err != nil {
		print(err.Error())
		return err
	}

	app.respondApi(w, r, rule, err)
	return nil
}
