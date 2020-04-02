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

type TagApp struct {
	TagAppID   string     `json:"tagAppID"`
	Name       string     `json:"name"`
	Weight     float64    `json:"weight"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

type TagAppTable struct{}

var TagApps = TagAppTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row TagApp) Equals(rhs TagApp) bool {
	if row.TagAppID != rhs.TagAppID {
		return false
	}
	if row.Name != rhs.Name {
		return false
	}
	if row.Weight != rhs.Weight {
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

// Scan a database row into a `TagApp` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t TagAppTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row TagApp,
	err error,
) {
	err = src.Scan(
		&row.TagAppID,
		&row.Name,
		&row.Weight,
		&row.ArchivedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan TagApp: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t TagAppTable) Table() string {
	return `tag_apps`
}

// View returns the table's view (for reading). May be the same as Table().
func (t TagAppTable) View() string {
	return `tag_apps`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t TagAppTable) SelectCols() string {
	return `tag_app_id,name,weight,archived_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_TagApp = `INSERT INTO tag_apps(
tag_app_id,
name,
weight,
archived_at
) VALUES (
 $1,$2,$3,$4
)`

// Insert will validate and insert a new `TagApp`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagAppTable) Insert(
	db DBi,
	row *TagApp,
) (
	err error,
) {

	if row.TagAppID == "" {
		row.TagAppID = crypto.NewUUID()
	}

	// Validate TagAppID.
	if err := validate.UUID(row.TagAppID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagAppID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_TagApp,
		row.TagAppID,
		row.Name,
		row.Weight,
		row.ArchivedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagApp.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_TagApp = `INSERT INTO tag_apps(
 tag_app_id,
 name,
 weight,
 archived_at
) VALUES (
 $1,$2,$3,$4
)
ON CONFLICT (tag_app_id)
DO UPDATE SET
 name=EXCLUDED.name,
 weight=EXCLUDED.weight,
 archived_at=EXCLUDED.archived_at`

func (t TagAppTable) Upsert(
	db DBi,
	row *TagApp,
) (
	err error,
) {

	if row.TagAppID == "" {
		row.TagAppID = crypto.NewUUID()
	}

	// Validate TagAppID.
	if err := validate.UUID(row.TagAppID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagAppID.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_TagApp,
		row.TagAppID,
		row.Name,
		row.Weight,
		row.ArchivedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("TagApp update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagApp update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagApp update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_TagApp = `UPDATE
 tag_apps
SET
 name=$1,
 weight=$2,
 archived_at=$3
WHERE
 tag_app_id=$4`

// Update updates the following column values:
//   - Name
//   - Weight
//   - ArchivedAt
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagAppTable) Update(
	db DBi,
	row *TagApp,
) (
	err error,
) {

	// Validate TagAppID.
	if err := validate.UUID(row.TagAppID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagAppID.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_TagApp,
		row.Name,
		row.Weight,
		row.ArchivedAt,
		row.TagAppID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("TagApp update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagApp update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagApp update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_TagApp = `DELETE FROM
 tag_apps
WHERE
 tag_app_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t TagAppTable) Delete(
	db DBi,
	TagAppID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_TagApp,
		TagAppID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("TagApp.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_TagApp = `SELECT
 tag_app_id,
 name,
 weight,
 archived_at
FROM
 tag_apps
WHERE
 tag_app_id=$1`

// Get returns the `TagApp` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t TagAppTable) Get(
	db DBi,
	TagAppID string,
) (
	row TagApp,
	err error,
) {
	src := db.QueryRow(getQuery_TagApp,
		TagAppID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `TagApp` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t TagAppTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []TagApp,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("TagApp.List failed: %w", err).
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
			Wrap("TagApp.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_TagApp = `CREATE TABLE tag_apps(
tag_app_id,
name,
weight,
archived_at
)`

func (t TagAppTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_TagApp)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagApp.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_TagApp)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagApp.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagApp.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.TagAppID,
			row.Name,
			row.Weight,
			row.ArchivedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("TagApp.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("TagApp.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
