package application

import (
	"net/http"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/slice"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminCustomerList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		ListArgs  database.CustomerListArgs
		Customers []database.Customer
		Session   shared.SharedSession
	}{Session: s}
	err = decoder.Decode(&ctx.ListArgs, r.Form)
	if err != nil {
		return err
	}

	// No pagination at the moment.
	ctx.ListArgs.Limit = 100000
	ctx.ListArgs.Offset = 0

	ctx.Customers, err = app.DBAL.CustomerList(ctx.ListArgs)
	if err != nil {
		return err
	}

	err = app.Views.ExecuteSST(w, "admin-customer-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminCustomerForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	type Dataset = struct {
		database.Dataset
		Mapped bool
	}

	ctx := struct {
		Session  shared.SharedSession
		Customer struct {
			database.Customer
			Archived bool
		}
		Datasets []Dataset
	}{Session: s}
	err = decoder.Decode(&ctx.Customer, r.Form)
	if err != nil {
		return err
	}

	if ctx.Customer.Customer.CustomerID != "" {
		ctx.Customer.Customer, err = app.DBAL.CustomerGet(ctx.Customer.CustomerID)

		mappedDatasets := []string{}
		customerDatasets, err := app.DBAL.CustomerDatasetMappingList(ctx.Customer.CustomerID)
		for _, cd := range customerDatasets {
			mappedDatasets = append(mappedDatasets, cd.DatasetEntity)
		}

		datasetsArchived := false
		datasets, err := app.DBAL.DatasetList(database.DatasetListArgs{
			Limit:    1000,
			Offset:   0,
			Archived: &datasetsArchived,
		})
		for _, ds := range datasets {
			ctx.Datasets = append(ctx.Datasets, Dataset{
				ds,
				slice.ContainsString(mappedDatasets, ds.DatasetID),
			})
		}

		if err != nil {
			return err
		}
	}

	err = app.Views.ExecuteSST(w, "admin-customer-form.gohtml", ctx)
	if err != nil {
		app.serverErr(w, r, err)
	}

	return err
}

func (app *App) AdminCustomerFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session  shared.SharedSession
		Customer struct {
			database.Customer
			Archived        bool
			DatasetMappings []struct {
				DatasetID string
			}
		}
	}{Session: s}
	err = decoder.Decode(&ctx.Customer, r.Form)
	if err != nil {
		return err
	}

	if (ctx.Customer.ArchivedAt != nil) != ctx.Customer.Archived {
		if ctx.Customer.Archived {
			t := time.Now()
			ctx.Customer.ArchivedAt = &t
		} else {
			ctx.Customer.ArchivedAt = nil
		}
	}

	err = app.DBAL.CustomerUpsert(&ctx.Customer.Customer)
	if err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Customers.Table(), ctx.Customer.CustomerID, "upsert", ctx.Customer)

	if ctx.Customer.CustomerID != "" {
		datasetIDs := []string{}
		for _, dsMapping := range ctx.Customer.DatasetMappings {
			if dsMapping.DatasetID != "" {
				datasetIDs = append(datasetIDs, dsMapping.DatasetID)
			}
		}
		err = app.DBAL.CustomerDatasetMappingSet(ctx.Customer.CustomerID, datasetIDs)
		defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.CustomerDatasets.Table(), ctx.Customer.CustomerID, "insert", datasetIDs)
	}

	if err != nil {
		return err
	}

	http.Redirect(w, r, "/admin/customers", http.StatusSeeOther)

	return nil
}
