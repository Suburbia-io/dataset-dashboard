package database

import (
	"database/sql"
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type Fingerprint = tables.Fingerprint

func (db *DBAL) FingerprintUpsert(fp *Fingerprint) (err error) {
	return tables.Fingerprints.Upsert(db, fp)
}

func (db *DBAL) FingerprintUpsertCSV(datasetID string, reader *csv.Reader) error {
	const (
		FP    = 0
		RAW   = 1
		COUNT = 2
	)

	header, err := reader.Read()
	if err != nil {
		return errors.Unexpected.Wrap("Failed to read header.")
	}
	if header[FP] != "fingerprint" || header[RAW] != "raw_text" || header[COUNT] != "count" {
		return errors.Unexpected.Wrap("Invalid header.")
	}

	return WithTx(db.DB, func(tx *sql.Tx) error {
		for {
			rec, err := reader.Read()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return errors.Unexpected.Wrap("Failed to read CSV: %w", err)
			}
			count, err := strconv.ParseInt(rec[COUNT], 10, 64)
			if err != nil {
				return errors.Unexpected.Wrap("Failed to parse count: %w", err)
			}

			fp := tables.Fingerprint{
				DatasetID:   datasetID,
				Fingerprint: rec[FP],
				RawText:     rec[RAW],
				Count:       int(count),
				UpdatedAt:   time.Now(),
			}
			if err := tables.Fingerprints.Upsert(tx, &fp); err != nil {
				return err
			}
		}
	})
}

func (db *DBAL) FingerprintAnnotationsUpsertCSV(
	datasetID string,
	reader *csv.Reader,
) error {
	const (
		FP    = 0
		NOTES = 1
	)

	header, err := reader.Read()
	if err != nil {
		return errors.Unexpected.Wrap("Failed to read header.")
	}
	if header[FP] != "fingerprint" || header[NOTES] != "annotations" {
		return errors.Unexpected.Wrap("Invalid header.")
	}

	return dbTX(db.DB, func(tx *sql.Tx) error {
		for {
			rec, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return errors.Unexpected.Wrap("Failed to read CSV: %w", err)
			}

			err = tables.Fingerprints.UpdateAnnotations(
				tx,
				datasetID,
				rec[FP],
				rec[NOTES])

			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DBAL) fingerprintExportToSQLite(
	datasetID string,
	sqlite *sql.DB,
) error {

	q := FPTagViewQuery{
		DatasetID: datasetID,
		TagAppID:  TAG_APP_HUMAN_POOL,
		Limit:     -1,
	}

	tagTypes, err := q.TagTypes(db)
	if err != nil {
		return err
	}

	// Build query for create table.
	builder := qb.Builder{}
	builder.Write(`CREATE TABLE fingerprints(`)
	builder.Write(`fingerprint TEXT NOT NULL PRIMARY KEY,`)
	builder.Write(`raw_text    TEXT NOT NULL,`)
	builder.Write(`annotations TEXT NOT NULL,`)
	builder.Write(`count       INT  NOT NULL`)

	for _, tt := range tagTypes {
		builder.Write(`,` + tt.TagType + ` TEXT NOT NULL`)
		builder.Write(`,` + tt.TagType + `_confidence REAL NOT NULL`)
	}
	builder.Write(`);`)

	query, _ := builder.MustBuild()
	if _, err := sqlite.Exec(query); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create fingerprints table: %w", err).
			Alert()
	}

	insertQuery := `INSERT INTO fingerprints(fingerprint,raw_text,annotations,count`
	placeholders := `(?,?,?,?`
	for _, tt := range tagTypes {
		insertQuery += `,` + tt.TagType + `,` + tt.TagType + `_confidence`
		placeholders += `,?,?`
	}
	insertQuery = insertQuery + `)VALUES` + placeholders + `)`

	return dbTX(sqlite, func(tx *sql.Tx) error {
		insertStmt, err := tx.Prepare(insertQuery)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to prepare fingerprint insert statement: %w", err).
				Alert()
		}

		return q.iterate(db, func(fpRow FPTagRow) error {
			row := make([]interface{}, 4+2*len(tagTypes))
			row[0] = fpRow.Fingerprint
			row[1] = fpRow.RawText
			row[2] = fpRow.Annotations
			row[3] = fpRow.Count

			for i, strPtr := range fpRow.ConsTags {
				if strPtr == nil {
					row[4+2*i+0] = ""
					row[4+2*i+1] = float64(0)
				} else {
					row[4+2*i] = *strPtr
					row[4+2*i+1] = *fpRow.ConsConfidences[i]
				}
			}

			if _, err := insertStmt.Exec(row...); err != nil {
				return errors.Unexpected.
					Wrap("Failed to insert fingerprint: %w", err).
					Alert()
			}

			return nil
		})
	})
}
