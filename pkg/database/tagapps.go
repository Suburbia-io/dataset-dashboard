package database

import (
	"database/sql"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

const TAG_APP_HUMAN_POOL = "00000000-0000-0000-0000-100000000000"

type TagApp = tables.TagApp

func (db *DBAL) TagAppUpsert(app *TagApp) error {
	return tables.TagApps.Upsert(db, app)
}

func (db *DBAL) TagAppGet(tagAppID string) (TagApp, error) {
	return tables.TagApps.Get(db, tagAppID)
}

type TagAppListArgs struct {
	Archived *bool
}

func (db *DBAL) TagAppList(args TagAppListArgs) ([]TagApp, error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.TagApps.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.TagApps.View())
	builder.Write(`WHERE TRUE`)

	if args.Archived != nil {
		if *args.Archived {
			builder.Write(`AND archived_at IS NOT NULL`)
		} else {
			builder.Write(`AND archived_at IS NULL`)
		}
	}

	builder.Write(`ORDER BY name`)
	query, queryArgs := builder.MustBuild()
	return tables.TagApps.List(db, query, queryArgs...)
}

const tagAppsExportCreateTableQuery = `CREATE TABLE tag_apps(
 tag_app_id TEXT NOT NULL PRIMARY KEY,
 name TEXT NOT NULL,
 weight REAL NOT NULL,
 archived_at TIMESTAMP)`

const tagAppsExportInsertQuery = `INSERT INTO tag_apps(
 tag_app_id,
 name,
 weight,
 archived_at
)VALUES(?,?,?,?)`

func (db *DBAL) tagAppExportToSQLite(
	datasetID string,
	sqlite *sql.DB,
) error {
	// Create table.
	if _, err := sqlite.Exec(tagAppsExportCreateTableQuery); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create tag_apps table: %w", err).
			Alert()
	}

	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.TagApps.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.TagApps.View())
	query, queryArgs := builder.MustBuild()

	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to list tag_apps: %w", err).
			Alert()
	}

	return dbTX(sqlite, func(tx *sql.Tx) error {
		insertStmt, err := tx.Prepare(tagAppsExportInsertQuery)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to prepare tag_apps insert statement: %w", err).
				Alert()
		}

		for rows.Next() {
			row, err := tables.TagApps.Scan(rows)
			if err != nil {
				return err
			}
			_, err = insertStmt.Exec(row.TagAppID, row.Name, row.Weight, row.ArchivedAt)
			if err != nil {
				return errors.Unexpected.
					Wrap("Failed to insert tag_app: %w", err).
					Alert()
			}
		}

		return nil
	})
}
