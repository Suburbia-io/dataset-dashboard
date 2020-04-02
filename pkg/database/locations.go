package database

import (
	"database/sql"
	"encoding/csv"
	"io"
	"strings"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type Location = tables.Location

func (db *DBAL) LocationUpsertCSV(datasetID string, reader *csv.Reader) error {
	const (
		HASH   = 0
		STRING = 1
	)

	header, err := reader.Read()
	if err != nil {
		return errors.Unexpected.Wrap("Failed to read header.")
	}
	if header[HASH] != "location_hash" || header[STRING] != "location_string" {
		return errors.Unexpected.Wrap("Invalid header.")
	}

	return dbTX(db.DB, func(tx *sql.Tx) error {
		for {
			rec, err := reader.Read()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return errors.Unexpected.Wrap("Failed to read CSV: %w", err)
			}

			loc := tables.Location{
				DatasetID:           datasetID,
				LocationHash:        rec[HASH],
				LocationString:      rec[STRING],
				CreatedAt:           time.Now(),
				GeonamesHierarchy:   []byte("{}"),
				GeonamesPostalCodes: []byte("{}"),
			}
			if err := tables.Locations.Upsert(tx, &loc); err != nil {
				return err
			}
		}
	})
}

func (db *DBAL) LocationUpdate(location *Location) error {
	return tables.Locations.Update(db, location)
}

func (db *DBAL) LocationGetByLocationHash(locationHash string) (Location, error) {
	return tables.Locations.GetByLocationHash(db, locationHash)
}

func (db *DBAL) LocationGet(datasetID string, locationHash string) (Location, error) {
	return tables.Locations.Get(db, datasetID, locationHash)
}

type LocationListArgs struct {
	Search string  `json:"search" schema:"search"`
	Status *string `json:"status" schema:"status"`
	Limit  int     `json:"limit" schema:"limit"`
	Offset int     `json:"offset" schema:"offset"`
}

var EUCountryCodes = [37]string{"AD", "AT", "BE", "BG", "BV", "CH", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FR",
	"GB", "GR", "HR", "HU", "IE", "IS", "IT", "LI", "LT", "LU", "LV", "MC", "MT", "NL", "NO", "PL", "PT", "RO", "SE",
	"SI", "SK", "SM", "VA"}

func (db *DBAL) LocationList(datasetID string, args LocationListArgs) ([]Location, error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Locations.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Locations.View())
	builder.Write(`WHERE dataset_id = $1`, datasetID)

	if args.Status != nil {
		if *args.Status == "approved" {
			builder.Write(`AND approved IS TRUE`)
		} else if *args.Status == "rejected" {
			builder.Write(`AND approved IS FALSE`)
		} else if *args.Status == "pending" {
			builder.Write(`AND approved IS NULL`)
		}
	}

	search := strings.TrimSpace(args.Search)
	if search != "" {
		if search == "EU" {
			builder.Write(`AND parsed_country_code IN (`)
			for _, countryCode := range EUCountryCodes {
				builder.Write("'" + countryCode + "',")
			}
			builder.Write("'')")

		} else {

			search := "%" + search + "%"
			builder.Write(`AND (location_hash ILIKE $1 OR
                          location_string ILIKE $1 OR
                          parsed_country_code ILIKE $1 OR
                          name ILIKE $1 OR
                          geoname_id::text ILIKE $1 OR
                          parent_name ILIKE $1 OR
                          parent_geoname_id::text ILIKE $1 OR
                          country_code ILIKE $1)`, search)
		}
	}

	builder.Write(`ORDER BY location_hash ASC`)
	builder.Write(`LIMIT $1 OFFSET $2`, args.Limit, args.Offset)

	query, queryArgs := builder.MustBuild()
	return tables.Locations.List(db, query, queryArgs...)
}

const locationExportCreateTableQuery = `CREATE TABLE locations(
 location_hash TEXT NOT NULL PRIMARY KEY,
 name TEXT NOT NULL,
 population INT NOT NULL,
 geoname_id INT,
 country_code TEXT NOT NULL,
 approved BOOLEAN);`

const locationExportInsertQuery = `INSERT INTO locations(
 location_hash,
 name,
 population,
 geoname_id,
 country_code,
 approved
)VALUES(?,?,?,?,?,?)`

func (db *DBAL) locationExportToSQLite(
	datasetID string,
	sqlite *sql.DB,
) error {
	// Create table.
	if _, err := sqlite.Exec(locationExportCreateTableQuery); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create location table: %w", err).
			Alert()
	}

	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Locations.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Locations.View())
	builder.Write(`WHERE dataset_id = $1`, datasetID)
	query, queryArgs := builder.MustBuild()

	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to list locations: %w", err).
			Alert()
	}

	return dbTX(sqlite, func(tx *sql.Tx) error {
		insertStmt, err := tx.Prepare(locationExportInsertQuery)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to prepare location insert statement: %w", err).
				Alert()
		}

		for rows.Next() {
			row, err := tables.Locations.Scan(rows)
			if err != nil {
				return err
			}
			_, err = insertStmt.Exec(
				row.LocationHash, row.Name, row.Population, row.GeonameID, row.CountryCode, row.Approved)
			if err != nil {
				return errors.Unexpected.
					Wrap("Failed to insert location: %w", err).
					Alert()
			}
		}

		return nil
	})
}
