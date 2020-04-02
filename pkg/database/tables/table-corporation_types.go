package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"

	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
)

type CorporationType struct {
	CorporationTypeID string    `json:"corporationTypeID"`
	DatasetID         string    `json:"datasetID"`
	CorporationType   string    `json:"corporationType"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CorporationTypeTable struct{}

var CorporationTypes = CorporationTypeTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row CorporationType) Equals(rhs CorporationType) bool {
	if row.CorporationTypeID != rhs.CorporationTypeID {
		return false
	}
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.CorporationType != rhs.CorporationType {
		return false
	}
	if row.Description != rhs.Description {
		return false
	}
	if !row.CreatedAt.Equal(rhs.CreatedAt) {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `CorporationType` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t CorporationTypeTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row CorporationType,
	err error,
) {
	err = src.Scan(
		&row.CorporationTypeID,
		&row.DatasetID,
		&row.CorporationType,
		&row.Description,
		&row.CreatedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan CorporationType: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t CorporationTypeTable) Table() string {
	return `corporation_types`
}

// View returns the table's view (for reading). May be the same as Table().
func (t CorporationTypeTable) View() string {
	return `corporation_types`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t CorporationTypeTable) SelectCols() string {
	return `corporation_type_id,dataset_id,corporation_type,description,created_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_CorporationType = `INSERT INTO corporation_types(
corporation_type_id,
dataset_id,
corporation_type,
description,
created_at
) VALUES (
 $1,$2,$3,$4,$5
)`

// Insert will validate and insert a new `CorporationType`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorporationTypeTable) Insert(
	db DBi,
	row *CorporationType,
) (
	err error,
) {

	if row.CorporationTypeID == "" {
		row.CorporationTypeID = crypto.NewUUID()
	}

	// Validate CorporationTypeID.
	if err := validate.UUID(row.CorporationTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorporationTypeID.")
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_CorporationType,
		row.CorporationTypeID,
		row.DatasetID,
		row.CorporationType,
		row.Description,
		row.CreatedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorporationType.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_CorporationType = `INSERT INTO corporation_types(
 corporation_type_id,
 dataset_id,
 corporation_type,
 description,
 created_at
) VALUES (
 $1,$2,$3,$4,$5
)
ON CONFLICT (corporation_type_id)
DO UPDATE SET
 corporation_type=EXCLUDED.corporation_type,
 description=EXCLUDED.description`

func (t CorporationTypeTable) Upsert(
	db DBi,
	row *CorporationType,
) (
	err error,
) {

	if row.CorporationTypeID == "" {
		row.CorporationTypeID = crypto.NewUUID()
	}

	// Validate CorporationTypeID.
	if err := validate.UUID(row.CorporationTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorporationTypeID.")
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_CorporationType,
		row.CorporationTypeID,
		row.DatasetID,
		row.CorporationType,
		row.Description,
		row.CreatedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("CorporationType update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("CorporationType update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorporationType update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_CorporationType = `UPDATE
 corporation_types
SET
 corporation_type=$1,
 description=$2
WHERE
 corporation_type_id=$3`

// Update updates the following column values:
//   - CorporationType
//   - Description
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorporationTypeTable) Update(
	db DBi,
	row *CorporationType,
) (
	err error,
) {

	// Validate CorporationTypeID.
	if err := validate.UUID(row.CorporationTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorporationTypeID.")
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_CorporationType,
		row.CorporationType,
		row.Description,
		row.CorporationTypeID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("CorporationType update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("CorporationType update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorporationType update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_CorporationType = `DELETE FROM
 corporation_types
WHERE
 corporation_type_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t CorporationTypeTable) Delete(
	db DBi,
	CorporationTypeID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_CorporationType,
		CorporationTypeID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("CorporationType.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_CorporationType = `SELECT
 corporation_type_id,
 dataset_id,
 corporation_type,
 description,
 created_at
FROM
 corporation_types
WHERE
 corporation_type_id=$1`

// Get returns the `CorporationType` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t CorporationTypeTable) Get(
	db DBi,
	CorporationTypeID string,
) (
	row CorporationType,
	err error,
) {
	src := db.QueryRow(getQuery_CorporationType,
		CorporationTypeID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `CorporationType` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t CorporationTypeTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []CorporationType,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("CorporationType.List failed: %w", err).
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
			Wrap("CorporationType.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_CorporationType = `CREATE TABLE corporation_types(
corporation_type_id,
dataset_id,
corporation_type,
description,
created_at
)`

func (t CorporationTypeTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_CorporationType)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorporationType.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_CorporationType)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorporationType.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorporationType.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.CorporationTypeID,
			row.DatasetID,
			row.CorporationType,
			row.Description,
			row.CreatedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("CorporationType.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("CorporationType.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
