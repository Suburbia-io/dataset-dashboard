package database

import (
	"database/sql"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type TagType = tables.TagType

func (db *DBAL) TagTypeUpsert(tt *TagType) error {
	return WithTx(db.DB, func(tx *sql.Tx) error {
		if err := tables.TagTypes.Upsert(tx, tt); err != nil {
			return err
		}
		return db.createFPTagSelectView(tx, tt.DatasetID)
	})
}

func (db *DBAL) TagTypeGet(datasetID, tagTypeID string) (TagType, error) {
	return tables.TagTypes.Get(db, datasetID, tagTypeID)
}

type TagTypeListArgs struct {
	DatasetID string `json:"datasetID"`
}

func (db *DBAL) TagTypeList(args TagTypeListArgs) ([]TagType, error) {
	return db.tagTypeList(db, args)
}

func (db *DBAL) tagTypeList(
	tx DBi,
	args TagTypeListArgs,
) (
	tagTypes []TagType,
	err error,
) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.TagTypes.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.TagTypes.Table())
	builder.Write(`WHERE dataset_id=$1`, args.DatasetID)
	builder.Write(`ORDER BY tag_type ASC`)

	query, queryArgs := builder.MustBuild()
	return tables.TagTypes.List(tx, query, queryArgs...)
}

func (db *DBAL) TagTypeDelete(datasetID, tagTypeID string) error {
	return WithTx(db.DB, func(tx *sql.Tx) error {
		if err := tables.TagTypes.Delete(tx, datasetID, tagTypeID); err != nil {
			return err
		}
		return db.createFPTagSelectView(tx, datasetID)
	})
}

const tagTypeExportSQLiteCreateQuery = `
CREATE TABLE tag_types(
  tag_type_id TEXT NOT NULL PRIMARY KEY,
  tag_type    TEXT NOT NULL
);
`

const tagTypeExportSQLiteInsertQuery = `
INSERT INTO tag_types(tag_type_id, tag_type) VALUES (?,?)
`

func (db *DBAL) tagTypeExportSQLite(datasetID string, sqlite *sql.DB) error {
	if _, err := sqlite.Exec(tagTypeExportSQLiteCreateQuery); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create tag type table: %w", err).
			Alert()
	}

	l, err := db.TagTypeList(TagTypeListArgs{DatasetID: datasetID})
	if err != nil {
		return err
	}

	return dbTX(sqlite, func(tx *sql.Tx) error {
		insertStmt, err := tx.Prepare(tagTypeExportSQLiteInsertQuery)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to prepare tag type insert statement: %w", err).
				Alert()
		}

		for _, tt := range l {
			if _, err := insertStmt.Exec(tt.TagTypeID, tt.TagType); err != nil {
				return errors.Unexpected.
					Wrap("Failed to insert tag type: %w", err).
					Alert()
			}
		}

		return nil
	})
}
