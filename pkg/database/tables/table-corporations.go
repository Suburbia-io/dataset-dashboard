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

type Corporation struct {
	CorporationID string     `json:"corporationID"`
	Exchange      string     `json:"exchange"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Slug          string     `json:"slug"`
	Isin          string     `json:"isin"`
	Cusip         string     `json:"cusip"`
	DatasetID     string     `json:"datasetID"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	ArchivedAt    *time.Time `json:"archivedAt"`
}

type CorporationTable struct{}

var Corporations = CorporationTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row Corporation) Equals(rhs Corporation) bool {
	if row.CorporationID != rhs.CorporationID {
		return false
	}
	if row.Exchange != rhs.Exchange {
		return false
	}
	if row.Code != rhs.Code {
		return false
	}
	if row.Name != rhs.Name {
		return false
	}
	if row.Slug != rhs.Slug {
		return false
	}
	if row.Isin != rhs.Isin {
		return false
	}
	if row.Cusip != rhs.Cusip {
		return false
	}
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if !row.CreatedAt.Equal(rhs.CreatedAt) {
		return false
	}

	if !row.UpdatedAt.Equal(rhs.UpdatedAt) {
		return false
	}

	if row.ArchivedAt != nil || rhs.ArchivedAt != nil {
		if row.ArchivedAt == nil || rhs.ArchivedAt == nil {
			return false
		}
		if !row.ArchivedAt.Equal(*rhs.ArchivedAt) {
			return false
		}
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `Corporation` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t CorporationTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row Corporation,
	err error,
) {
	err = src.Scan(
		&row.CorporationID,
		&row.Exchange,
		&row.Code,
		&row.Name,
		&row.Slug,
		&row.Isin,
		&row.Cusip,
		&row.DatasetID,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.ArchivedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan Corporation: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t CorporationTable) Table() string {
	return `corporations`
}

// View returns the table's view (for reading). May be the same as Table().
func (t CorporationTable) View() string {
	return `corporations`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t CorporationTable) SelectCols() string {
	return `corporation_id,exchange,code,name,slug,isin,cusip,dataset_id,created_at,updated_at,archived_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_Corporation = `INSERT INTO corporations(
corporation_id,
exchange,
code,
name,
slug,
isin,
cusip,
dataset_id,
created_at,
updated_at,
archived_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
)`

// Insert will validate and insert a new `Corporation`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorporationTable) Insert(
	db DBi,
	row *Corporation,
) (
	err error,
) {

	if row.CorporationID == "" {
		row.CorporationID = crypto.NewUUID()
	}

	// Validate CorporationID.
	if err := validate.UUID(row.CorporationID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorporationID.")
	}

	// Validate Slug.
	if err := validate.Slug(row.Slug); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Slug.")
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_Corporation,
		row.CorporationID,
		row.Exchange,
		row.Code,
		row.Name,
		row.Slug,
		row.Isin,
		row.Cusip,
		row.DatasetID,
		row.CreatedAt,
		row.UpdatedAt,
		row.ArchivedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Corporation.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_Corporation = `INSERT INTO corporations(
 corporation_id,
 exchange,
 code,
 name,
 slug,
 isin,
 cusip,
 dataset_id,
 created_at,
 updated_at,
 archived_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
)
ON CONFLICT (corporation_id)
DO UPDATE SET
 exchange=EXCLUDED.exchange,
 code=EXCLUDED.code,
 name=EXCLUDED.name,
 slug=EXCLUDED.slug,
 isin=EXCLUDED.isin,
 cusip=EXCLUDED.cusip`

func (t CorporationTable) Upsert(
	db DBi,
	row *Corporation,
) (
	err error,
) {

	if row.CorporationID == "" {
		row.CorporationID = crypto.NewUUID()
	}

	// Validate CorporationID.
	if err := validate.UUID(row.CorporationID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorporationID.")
	}

	// Validate Slug.
	if err := validate.Slug(row.Slug); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Slug.")
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_Corporation,
		row.CorporationID,
		row.Exchange,
		row.Code,
		row.Name,
		row.Slug,
		row.Isin,
		row.Cusip,
		row.DatasetID,
		row.CreatedAt,
		row.UpdatedAt,
		row.ArchivedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Corporation update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Corporation update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Corporation update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_Corporation = `UPDATE
 corporations
SET
 exchange=$1,
 code=$2,
 name=$3,
 slug=$4,
 isin=$5,
 cusip=$6
WHERE
 corporation_id=$7`

// Update updates the following column values:
//   - Exchange
//   - Code
//   - Name
//   - Slug
//   - Isin
//   - Cusip
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorporationTable) Update(
	db DBi,
	row *Corporation,
) (
	err error,
) {

	// Validate CorporationID.
	if err := validate.UUID(row.CorporationID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorporationID.")
	}

	// Validate Slug.
	if err := validate.Slug(row.Slug); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Slug.")
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_Corporation,
		row.Exchange,
		row.Code,
		row.Name,
		row.Slug,
		row.Isin,
		row.Cusip,
		row.CorporationID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Corporation update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Corporation update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Corporation update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_Corporation = `DELETE FROM
 corporations
WHERE
 corporation_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t CorporationTable) Delete(
	db DBi,
	CorporationID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_Corporation,
		CorporationID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("Corporation.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_Corporation = `SELECT
 corporation_id,
 exchange,
 code,
 name,
 slug,
 isin,
 cusip,
 dataset_id,
 created_at,
 updated_at,
 archived_at
FROM
 corporations
WHERE
 corporation_id=$1`

// Get returns the `Corporation` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t CorporationTable) Get(
	db DBi,
	CorporationID string,
) (
	row Corporation,
	err error,
) {
	src := db.QueryRow(getQuery_Corporation,
		CorporationID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `Corporation` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t CorporationTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []Corporation,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("Corporation.List failed: %w", err).
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
			Wrap("Corporation.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_Corporation = `CREATE TABLE corporations(
corporation_id,
exchange,
code,
name,
slug,
isin,
cusip,
dataset_id,
created_at,
updated_at,
archived_at
)`

func (t CorporationTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_Corporation)
	if err != nil {
		return errors.Unexpected.
			Wrap("Corporation.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_Corporation)
	if err != nil {
		return errors.Unexpected.
			Wrap("Corporation.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Corporation.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.CorporationID,
			row.Exchange,
			row.Code,
			row.Name,
			row.Slug,
			row.Isin,
			row.Cusip,
			row.DatasetID,
			row.CreatedAt,
			row.UpdatedAt,
			row.ArchivedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("Corporation.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("Corporation.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
