package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"

	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
)

type CorpMapping struct {
	CorpMappingID string `json:"corpMappingID"`
	CorpTypeID    string `json:"corpTypeID"`
	TagTypeID     string `json:"tagTypeID"`
	TagID         string `json:"tagID"`
}

type CorpMappingTable struct{}

var CorpMappings = CorpMappingTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row CorpMapping) Equals(rhs CorpMapping) bool {
	if row.CorpMappingID != rhs.CorpMappingID {
		return false
	}
	if row.CorpTypeID != rhs.CorpTypeID {
		return false
	}
	if row.TagTypeID != rhs.TagTypeID {
		return false
	}
	if row.TagID != rhs.TagID {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `CorpMapping` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t CorpMappingTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row CorpMapping,
	err error,
) {
	err = src.Scan(
		&row.CorpMappingID,
		&row.CorpTypeID,
		&row.TagTypeID,
		&row.TagID)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan CorpMapping: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t CorpMappingTable) Table() string {
	return `corp_mappings`
}

// View returns the table's view (for reading). May be the same as Table().
func (t CorpMappingTable) View() string {
	return `corp_mappings`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t CorpMappingTable) SelectCols() string {
	return `corp_mapping_id,corp_type_id,tag_type_id,tag_id`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_CorpMapping = `INSERT INTO corp_mappings(
corp_mapping_id,
corp_type_id,
tag_type_id,
tag_id
) VALUES (
 $1,$2,$3,$4
)`

// Insert will validate and insert a new `CorpMapping`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorpMappingTable) Insert(
	db DBi,
	row *CorpMapping,
) (
	err error,
) {

	if row.CorpMappingID == "" {
		row.CorpMappingID = crypto.NewUUID()
	}

	// Validate CorpMappingID.
	if err := validate.UUID(row.CorpMappingID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingID.")
	}

	// Validate CorpTypeID.
	if err := validate.UUID(row.CorpTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpTypeID.")
	}

	// Validate TagTypeID.
	if err := validate.UUID(row.TagTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagTypeID.")
	}

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_CorpMapping,
		row.CorpMappingID,
		row.CorpTypeID,
		row.TagTypeID,
		row.TagID)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorpMapping.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_CorpMapping = `INSERT INTO corp_mappings(
 corp_mapping_id,
 corp_type_id,
 tag_type_id,
 tag_id
) VALUES (
 $1,$2,$3,$4
)
ON CONFLICT (corp_mapping_id)
DO UPDATE SET
 corp_type_id=EXCLUDED.corp_type_id,
 tag_type_id=EXCLUDED.tag_type_id,
 tag_id=EXCLUDED.tag_id`

func (t CorpMappingTable) Upsert(
	db DBi,
	row *CorpMapping,
) (
	err error,
) {

	if row.CorpMappingID == "" {
		row.CorpMappingID = crypto.NewUUID()
	}

	// Validate CorpMappingID.
	if err := validate.UUID(row.CorpMappingID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingID.")
	}

	// Validate CorpTypeID.
	if err := validate.UUID(row.CorpTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpTypeID.")
	}

	// Validate TagTypeID.
	if err := validate.UUID(row.TagTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagTypeID.")
	}

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_CorpMapping,
		row.CorpMappingID,
		row.CorpTypeID,
		row.TagTypeID,
		row.TagID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("CorpMapping update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("CorpMapping update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorpMapping update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_CorpMapping = `UPDATE
 corp_mappings
SET
 corp_type_id=$1,
 tag_type_id=$2,
 tag_id=$3
WHERE
 corp_mapping_id=$4`

// Update updates the following column values:
//   - CorpTypeID
//   - TagTypeID
//   - TagID
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorpMappingTable) Update(
	db DBi,
	row *CorpMapping,
) (
	err error,
) {

	// Validate CorpMappingID.
	if err := validate.UUID(row.CorpMappingID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingID.")
	}

	// Validate CorpTypeID.
	if err := validate.UUID(row.CorpTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpTypeID.")
	}

	// Validate TagTypeID.
	if err := validate.UUID(row.TagTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagTypeID.")
	}

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_CorpMapping,
		row.CorpTypeID,
		row.TagTypeID,
		row.TagID,
		row.CorpMappingID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("CorpMapping update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("CorpMapping update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorpMapping update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_CorpMapping = `DELETE FROM
 corp_mappings
WHERE
 corp_mapping_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t CorpMappingTable) Delete(
	db DBi,
	CorpMappingID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_CorpMapping,
		CorpMappingID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("CorpMapping.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_CorpMapping = `SELECT
 corp_mapping_id,
 corp_type_id,
 tag_type_id,
 tag_id
FROM
 corp_mappings
WHERE
 corp_mapping_id=$1`

// Get returns the `CorpMapping` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t CorpMappingTable) Get(
	db DBi,
	CorpMappingID string,
) (
	row CorpMapping,
	err error,
) {
	src := db.QueryRow(getQuery_CorpMapping,
		CorpMappingID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `CorpMapping` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t CorpMappingTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []CorpMapping,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("CorpMapping.List failed: %w", err).
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
			Wrap("CorpMapping.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_CorpMapping = `CREATE TABLE corp_mappings(
corp_mapping_id,
corp_type_id,
tag_type_id,
tag_id
)`

func (t CorpMappingTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_CorpMapping)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorpMapping.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_CorpMapping)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorpMapping.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorpMapping.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.CorpMappingID,
			row.CorpTypeID,
			row.TagTypeID,
			row.TagID,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("CorpMapping.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("CorpMapping.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
