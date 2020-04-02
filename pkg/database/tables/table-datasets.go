package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/sanitize"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type Dataset struct {
	DatasetID  string     `json:"datasetID"`
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	Manageable bool       `json:"manageable"`
	CreatedAt  time.Time  `json:"createdAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

type DatasetTable struct{}

var Datasets = DatasetTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row Dataset) Equals(rhs Dataset) bool {
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.Name != rhs.Name {
		return false
	}
	if row.Slug != rhs.Slug {
		return false
	}
	if row.Manageable != rhs.Manageable {
		return false
	}
	if !row.CreatedAt.Equal(rhs.CreatedAt) {
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

// Scan a database row into a `Dataset` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t DatasetTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row Dataset,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.Name,
		&row.Slug,
		&row.Manageable,
		&row.CreatedAt,
		&row.ArchivedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan Dataset: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t DatasetTable) Table() string {
	return `datasets`
}

// View returns the table's view (for reading). May be the same as Table().
func (t DatasetTable) View() string {
	return `datasets`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t DatasetTable) SelectCols() string {
	return `dataset_id,name,slug,manageable,created_at,archived_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_Dataset = `INSERT INTO datasets(
dataset_id,
name,
slug,
manageable,
created_at,
archived_at
) VALUES (
 $1,$2,$3,$4,$5,$6
)`

// Insert will validate and insert a new `Dataset`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t DatasetTable) Insert(
	db DBi,
	row *Dataset,
) (
	err error,
) {

	if row.DatasetID == "" {
		row.DatasetID = crypto.NewUUID()
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.NonEmptyString(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Sanitize Slug.
	row.Slug = sanitize.SingleLineString(row.Slug)

	// Validate Slug.
	if err := validate.Slug(row.Slug); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Slug.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_Dataset,
		row.DatasetID,
		row.Name,
		row.Slug,
		row.Manageable,
		row.CreatedAt,
		row.ArchivedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Dataset.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_Dataset = `INSERT INTO datasets(
 dataset_id,
 name,
 slug,
 manageable,
 created_at,
 archived_at
) VALUES (
 $1,$2,$3,$4,$5,$6
)
ON CONFLICT (dataset_id)
DO UPDATE SET
 name=EXCLUDED.name,
 slug=EXCLUDED.slug,
 manageable=EXCLUDED.manageable,
 archived_at=EXCLUDED.archived_at`

func (t DatasetTable) Upsert(
	db DBi,
	row *Dataset,
) (
	err error,
) {

	if row.DatasetID == "" {
		row.DatasetID = crypto.NewUUID()
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.NonEmptyString(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Sanitize Slug.
	row.Slug = sanitize.SingleLineString(row.Slug)

	// Validate Slug.
	if err := validate.Slug(row.Slug); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Slug.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_Dataset,
		row.DatasetID,
		row.Name,
		row.Slug,
		row.Manageable,
		row.CreatedAt,
		row.ArchivedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Dataset update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Dataset update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Dataset update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_Dataset = `UPDATE
 datasets
SET
 name=$1,
 slug=$2,
 manageable=$3,
 archived_at=$4
WHERE
 dataset_id=$5`

// Update updates the following column values:
//   - Name
//   - Slug
//   - Manageable
//   - ArchivedAt
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t DatasetTable) Update(
	db DBi,
	row *Dataset,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.NonEmptyString(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Sanitize Slug.
	row.Slug = sanitize.SingleLineString(row.Slug)

	// Validate Slug.
	if err := validate.Slug(row.Slug); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Slug.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_Dataset,
		row.Name,
		row.Slug,
		row.Manageable,
		row.ArchivedAt,
		row.DatasetID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Dataset update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Dataset update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Dataset update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_Dataset = `DELETE FROM
 datasets
WHERE
 dataset_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t DatasetTable) Delete(
	db DBi,
	DatasetID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_Dataset,
		DatasetID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("Dataset.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_Dataset = `SELECT
 dataset_id,
 name,
 slug,
 manageable,
 created_at,
 archived_at
FROM
 datasets
WHERE
 dataset_id=$1`

// Get returns the `Dataset` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t DatasetTable) Get(
	db DBi,
	DatasetID string,
) (
	row Dataset,
	err error,
) {
	src := db.QueryRow(getQuery_Dataset,
		DatasetID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetBySlug
// ----------------------------------------------------------------------------

const getQuery_Dataset_bySlug = `SELECT
 dataset_id,
 name,
 slug,
 manageable,
 created_at,
 archived_at
FROM
  datasets
WHERE
 slug=$1`

// GetBySlug return the Dataset object by a natural key.
func (t DatasetTable) GetBySlug(
	db DBi,
	Slug string,
) (
	row Dataset,
	err error,
) {

	src := db.QueryRow(getQuery_Dataset_bySlug,
		Slug)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `Dataset` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t DatasetTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []Dataset,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("Dataset.List failed: %w", err).
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
			Wrap("Dataset.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_Dataset = `CREATE TABLE datasets(
dataset_id,
name,
slug,
manageable,
created_at,
archived_at
)`

func (t DatasetTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_Dataset)
	if err != nil {
		return errors.Unexpected.
			Wrap("Dataset.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_Dataset)
	if err != nil {
		return errors.Unexpected.
			Wrap("Dataset.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Dataset.List failed: %w", err).
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
			row.Name,
			row.Slug,
			row.Manageable,
			row.CreatedAt,
			row.ArchivedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("Dataset.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("Dataset.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
