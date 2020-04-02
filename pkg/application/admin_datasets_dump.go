package application

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database"

	"github.com/Suburbia-io/dashboard/pkg/helpers/slice"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"

	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminDumpDatasetForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	ctx := struct {
		Session        shared.SharedSession
		Dataset        database.Dataset
		DumpableTables []string
	}{Session: s}
	err := decoder.Decode(&ctx.Dataset, r.Form)
	if err != nil {
		return err
	}

	ctx.Dataset, err = app.DBAL.DatasetGet(ctx.Dataset.DatasetID)
	if err != nil {
		return err
	}

	ctx.DumpableTables = []string{
		tables.Fingerprints.View(),
		tables.TagTypes.View(),
		tables.Tags.View(),
		tables.Datasets.View(),
		tables.Corporations.View(),
		tables.CorpMappings.View(),
		tables.CorpMappingRules.View(),
		tables.CorporationTypes.View(),
		tables.TagAppHistoricalTags.View(),
		tables.TagAppTags.View(),
		tables.TagApps.View(),
		tables.ConsensusTags.View(),
		tables.Locations.View(),
	}

	err = app.Views.ExecuteSST(w, "admin-dataset-dump-form.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminDumpDatasetSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	ctx := struct {
		Session shared.SharedSession
		Dataset database.Dataset
	}{Session: s}
	err := decoder.Decode(&ctx.Dataset, r.Form)
	if err != nil {
		return err
	}

	// check if dataset exists
	ctx.Dataset, err = app.DBAL.DatasetGet(ctx.Dataset.DatasetID)
	if err != nil {
		return err
	}

	// Get the table names, comma separated
	tableQueryParam := r.URL.Query().Get("tables")
	var tablesToExport []string

	// remove duplicates to prevent retrieving the same table multiple times
	for _, table := range strings.Split(tableQueryParam, ",") {
		tableName := strings.TrimSpace(table)
		if !slice.ContainsString(tablesToExport, tableName) {
			tablesToExport = append(tablesToExport, tableName)
		}
	}

	// create a temporary sqlite file
	dir, err := ioutil.TempDir("", "export-dataset-v2")
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, "dump.sqlite")
	if err != nil {
		return err
	}
	sqlite, err := sql.Open("sqlite3", dbPath+"?_journal=OFF&_sync=OFF")
	if err != nil {
		return err
	}

	// dump the desired tables into the sqlite file
	for _, table := range tablesToExport {
		switch table {
		case tables.Fingerprints.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.Fingerprints)
		case tables.TagTypes.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.TagTypes)
		case tables.Tags.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.Tags)
		case tables.Datasets.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.Datasets)
		case tables.Corporations.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.Corporations)
		case tables.CorpMappings.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.CorpMappings)
		case tables.CorpMappingRules.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.CorpMappingRules)
		case tables.CorporationTypes.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.CorporationTypes)
		case tables.TagAppHistoricalTags.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.TagAppHistoricalTags)
		case tables.TagAppTags.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.TagAppTags)
		case tables.TagApps.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.TagApps)
		case tables.ConsensusTags.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.ConsensusTags)
		case tables.Locations.Table():
			err = app.DBAL.DumpTableToSqlite(sqlite, ctx.Dataset.DatasetID, tables.Locations)
		}
	}

	if err != nil {
		return err
	}

	if err := sqlite.Close(); err != nil {
		return err
	}

	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s_%s.sqlite"`, "db-export", time.Now().Format("2006-01-02-15:04:05")))
	http.ServeFile(w, r, dbPath)
	return nil
}
