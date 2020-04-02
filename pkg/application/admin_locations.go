package application

import (
	"encoding/json"
	"net/http"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminLocationList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		ListArgs      database.LocationListArgs
		Dataset       database.Dataset
		Locations     []database.Location
		LocationsJSON string
		Session       shared.SharedSession
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

	// No pagination at the moment.
	ctx.ListArgs.Limit = 200
	ctx.ListArgs.Offset = 0

	ctx.Locations, err = app.DBAL.LocationList(datasetID, ctx.ListArgs)
	if err != nil {
		return err
	}

	locationsJSON, err := json.Marshal(ctx.Locations)
	if err != nil {
		return err
	}
	ctx.LocationsJSON = string(locationsJSON)

	err = app.Views.ExecuteSST(w, "admin-location-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminLocationSearchGeonamesFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	dataset, err := app.DBAL.DatasetGet(r.FormValue("DatasetID"))
	if err != nil {
		return err
	}

	go app.SearchGeonamesForLocations(dataset.DatasetID)

	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Locations.Table(), dataset.DatasetID, "searchGeonamesForLocations", "")
	http.Redirect(w, r, "/admin/datasets/"+dataset.DatasetID+"/locations/", http.StatusSeeOther)
	return nil
}
