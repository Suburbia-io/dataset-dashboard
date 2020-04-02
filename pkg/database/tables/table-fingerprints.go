package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/sanitize"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type Fingerprint struct {
	DatasetID   string    `json:"datasetID"`
	Fingerprint string    `json:"fingerprint"`
	RawText     string    `json:"rawText"`
	Annotations string    `json:"annotations"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Count       int       `json:"count"`
}

type FingerprintTable struct{}

var Fingerprints = FingerprintTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row Fingerprint) Equals(rhs Fingerprint) bool {
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.Fingerprint != rhs.Fingerprint {
		return false
	}
	if row.RawText != rhs.RawText {
		return false
	}
	if row.Annotations != rhs.Annotations {
		return false
	}
	if !row.UpdatedAt.Equal(rhs.UpdatedAt) {
		return false
	}

	if row.Count != rhs.Count {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `Fingerprint` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t FingerprintTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row Fingerprint,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.Fingerprint,
		&row.RawText,
		&row.Annotations,
		&row.UpdatedAt,
		&row.Count)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan Fingerprint: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t FingerprintTable) Table() string {
	return `fingerprints`
}

// View returns the table's view (for reading). May be the same as Table().
func (t FingerprintTable) View() string {
	return `fingerprints`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t FingerprintTable) SelectCols() string {
	return `dataset_id,fingerprint,raw_text,annotations,updated_at,count`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_Fingerprint = `INSERT INTO fingerprints(
dataset_id,
fingerprint,
raw_text,
annotations,
updated_at,
count
) VALUES (
 $1,$2,$3,$4,$5,$6
)`

// Insert will validate and insert a new `Fingerprint`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t FingerprintTable) Insert(
	db DBi,
	row *Fingerprint,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize Fingerprint.
	row.Fingerprint = sanitize.SingleLineString(row.Fingerprint)

	// Validate Fingerprint.
	if err := validate.NonEmptyString(row.Fingerprint); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Fingerprint.")
	}

	// Sanitize RawText.
	row.RawText = sanitize.SingleLineString(row.RawText)

	// Sanitize Annotations.
	row.Annotations = sanitize.SingleLineString(row.Annotations)

	// Execute query.
	_, err = db.Exec(insertQuery_Fingerprint,
		row.DatasetID,
		row.Fingerprint,
		row.RawText,
		row.Annotations,
		row.UpdatedAt,
		row.Count)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Fingerprint.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_Fingerprint = `INSERT INTO fingerprints(
 dataset_id,
 fingerprint,
 raw_text,
 annotations,
 updated_at,
 count
) VALUES (
 $1,$2,$3,$4,$5,$6
)
ON CONFLICT (dataset_id,fingerprint)
DO UPDATE SET
 raw_text=EXCLUDED.raw_text,
 updated_at=EXCLUDED.updated_at,
 count=EXCLUDED.count`

func (t FingerprintTable) Upsert(
	db DBi,
	row *Fingerprint,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize Fingerprint.
	row.Fingerprint = sanitize.SingleLineString(row.Fingerprint)

	// Validate Fingerprint.
	if err := validate.NonEmptyString(row.Fingerprint); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Fingerprint.")
	}

	// Sanitize RawText.
	row.RawText = sanitize.SingleLineString(row.RawText)

	// Sanitize Annotations.
	row.Annotations = sanitize.SingleLineString(row.Annotations)

	// Execute query.
	result, err := db.Exec(upsertQuery_Fingerprint,
		row.DatasetID,
		row.Fingerprint,
		row.RawText,
		row.Annotations,
		row.UpdatedAt,
		row.Count)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Fingerprint update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Fingerprint update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Fingerprint update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_Fingerprint = `UPDATE
 fingerprints
SET
 raw_text=$1,
 updated_at=$2,
 count=$3
WHERE
 dataset_id=$4 AND 
 fingerprint=$5`

// Update updates the following column values:
//   - RawText
//   - UpdatedAt
//   - Count
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t FingerprintTable) Update(
	db DBi,
	row *Fingerprint,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize Fingerprint.
	row.Fingerprint = sanitize.SingleLineString(row.Fingerprint)

	// Validate Fingerprint.
	if err := validate.NonEmptyString(row.Fingerprint); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Fingerprint.")
	}

	// Sanitize RawText.
	row.RawText = sanitize.SingleLineString(row.RawText)

	// Sanitize Annotations.
	row.Annotations = sanitize.SingleLineString(row.Annotations)

	// Execute query.
	result, err := db.Exec(updateQuery_Fingerprint,
		row.RawText,
		row.UpdatedAt,
		row.Count,
		row.DatasetID,
		row.Fingerprint)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Fingerprint update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Fingerprint update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Fingerprint update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateAnnotations
// ----------------------------------------------------------------------------

const updateQuery_Fingerprint_Annotations = `UPDATE
 fingerprints
SET
 annotations=$1
WHERE
 dataset_id=$2 AND 
 fingerprint=$3`

// UpdateAnnotations will attempt to update the Annotations column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t FingerprintTable) UpdateAnnotations(
	db DBi,
	DatasetID string,
	Fingerprint string,
	Annotations string,
) (
	err error,
) {
	Annotations = sanitize.SingleLineString(Annotations)

	result, err := db.Exec(updateQuery_Fingerprint_Annotations,
		Annotations,
		DatasetID,
		Fingerprint)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Fingerprint update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Fingerprint update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Fingerprint update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_Fingerprint = `DELETE FROM
 fingerprints
WHERE
 dataset_id=$1 AND 
 fingerprint=$2`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t FingerprintTable) Delete(
	db DBi,
	DatasetID string,
	Fingerprint string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_Fingerprint,
		DatasetID,
		Fingerprint)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("Fingerprint.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_Fingerprint = `SELECT
 dataset_id,
 fingerprint,
 raw_text,
 annotations,
 updated_at,
 count
FROM
 fingerprints
WHERE
 dataset_id=$1 AND 
 fingerprint=$2`

// Get returns the `Fingerprint` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t FingerprintTable) Get(
	db DBi,
	DatasetID string,
	Fingerprint string,
) (
	row Fingerprint,
	err error,
) {
	src := db.QueryRow(getQuery_Fingerprint,
		DatasetID,
		Fingerprint)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `Fingerprint` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t FingerprintTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []Fingerprint,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("Fingerprint.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return l, err
		}
		l = append(l, row)
	}

	if err := rows.Err(); err != nil {
		return l, errors.Unexpected.
			Wrap("Fingerprint.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_Fingerprint = `CREATE TABLE fingerprints(
dataset_id,
fingerprint,
raw_text,
annotations,
updated_at,
count
)`

func (t FingerprintTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_Fingerprint)
	if err != nil {
		return errors.Unexpected.
			Wrap("Fingerprint.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_Fingerprint)
	if err != nil {
		return errors.Unexpected.
			Wrap("Fingerprint.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Fingerprint.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.DatasetID,
			row.Fingerprint,
			row.RawText,
			row.Annotations,
			row.UpdatedAt,
			row.Count,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("Fingerprint.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("Fingerprint.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
