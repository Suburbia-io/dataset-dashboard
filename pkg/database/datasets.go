package database

import (
	"database/sql"
	"strings"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
	_ "github.com/mattn/go-sqlite3"
)

type Dataset = tables.Dataset

func (db *DBAL) DatasetUpsert(ds *Dataset) error {
	return WithTx(db.DB, func(tx *sql.Tx) error {
		if ds.DatasetID == "" {
			ds.CreatedAt = time.Now()
		}
		if err := tables.Datasets.Upsert(tx, ds); err != nil {
			return err
		}
		return db.createFPTagSelectView(tx, ds.DatasetID)
	})
}

func (db *DBAL) DatasetGet(datasetID string) (Dataset, error) {
	return tables.Datasets.Get(db, datasetID)
}

func (db *DBAL) DatasetGetBySlug(slug string) (dataset Dataset, err error) {
	return tables.Datasets.GetBySlug(db, slug)
}

type DatasetListArgs struct {
	Search     string `json:"search" schema:"search"`
	Archived   *bool  `json:"archived" schema:"archived"`
	Limit      int    `json:"limit" schema:"limit"`
	Offset     int    `json:"offset" schema:"offset"`
	Manageable *bool  `json:"manageable" schema:"manageable"`
}

func (db *DBAL) DatasetList(args DatasetListArgs) (datasets []Dataset, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Datasets.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Datasets.Table())
	builder.Write(`WHERE TRUE`)

	if args.Archived != nil {
		if *args.Archived {
			builder.Write(`AND archived_at IS NOT NULL`)
		} else {
			builder.Write(`AND archived_at IS NULL`)
		}
	}

	if args.Manageable != nil {
		if *args.Manageable {
			builder.Write(`AND manageable IS TRUE`)
		} else {
			builder.Write(`AND manageable IS FALSE`)
		}
	}

	if args.Search != "" {
		search := "%" + args.Search + "%"
		builder.Write(`AND (name ILIKE $1 OR slug ILIKE $1)`, search)
	}

	builder.Write(`ORDER BY name ASC`)
	builder.Write(`LIMIT $1 OFFSET $2`, args.Limit, args.Offset)

	query, queryArgs := builder.MustBuild()
	return tables.Datasets.List(db, query, queryArgs...)
}

func (db *DBAL) DatasetListActiveForCustomer(customerID string) (datasets []Dataset, err error) {
	// Prefix columns with table name to avoid ambiguous columns.
	var datasetsSelectCols []string
	for _, col := range strings.Split(tables.Datasets.SelectCols(), ",") {
		datasetsSelectCols = append(datasetsSelectCols, tables.Datasets.Table()+"."+col)
	}

	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(strings.Join(datasetsSelectCols, ","))
	builder.Write(`FROM`)
	builder.Write(tables.Datasets.Table())
	builder.Write(`JOIN`)
	builder.Write(tables.CustomerDatasets.Table())
	builder.Write(`ON datasets.dataset_id = customer_datasets.dataset_entity`)
	builder.Write(`WHERE archived_at IS NULL AND`)
	builder.Write(`customer_entity = $1`, customerID)
	builder.Write(`ORDER BY name ASC`)
	query, queryArgs := builder.MustBuild()
	return tables.Datasets.List(db, query, queryArgs...)
}

func (db *DBAL) DatasetExportSQLite(datasetID, outPath string, fullExport bool) error {
	sqlite, err := sql.Open("sqlite3", outPath)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to open output path: %w", err).
			Alert()
	}
	if err := db.fingerprintExportToSQLite(datasetID, sqlite); err != nil {
		return err
	}
	if err := db.locationExportToSQLite(datasetID, sqlite); err != nil {
		return err
	}
	if err := db.tagTypeExportSQLite(datasetID, sqlite); err != nil {
		return err
	}
	if err := db.tagExportSQLite(datasetID, sqlite); err != nil {
		return err
	}
	if fullExport {
		if err := db.tagAppExportToSQLite(datasetID, sqlite); err != nil {
			return err
		}
		if err := db.tagAppTagsExportSQLite(datasetID, sqlite); err != nil {
			return err
		}
	}
	if err := sqlite.Close(); err != nil {
		return errors.Unexpected.
			Wrap("Failed to close sqlite file: %w", err).
			Alert()
	}
	return nil
}
