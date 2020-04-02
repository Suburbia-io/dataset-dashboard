package application

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/shared"
)

// This handler doesn't use the error handling from the page wrapper function because it's regularly used in scripts
// which use the response codes.
func (app *App) AdminAPIExportDataset_v1(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	// Get the dataset ID.
	datasetID := r.URL.Query().Get("datasetID")

	fullExport := r.URL.Query().Get("fullExport") == "1"

	dataset, err := app.DBAL.DatasetGet(datasetID)
	if err != nil {
		if errors.DBNotFound.Is(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return nil
	}

	dir, err := ioutil.TempDir("", "export-dataset-v1")
	if err != nil {
		err = errors.Unexpected.
			Wrap("Failed to create temp dir: %w", err).
			Alert()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, "db.sqlite")
	if err := app.DBAL.DatasetExportSQLite(datasetID, dbPath, fullExport); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s_%s.sqlite"`, dataset.Slug, time.Now().Format("2006-01-02-15:04:05")))
	http.ServeFile(w, r, dbPath)
	return nil
}
