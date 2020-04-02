package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/sanitize"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type TagType struct {
	DatasetID   string `json:"datasetID"`
	TagTypeID   string `json:"tagTypeID"`
	TagType     string `json:"tagType"`
	Description string `json:"description"`
}

type TagTypeTable struct{}

var TagTypes = TagTypeTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row TagType) Equals(rhs TagType) bool {
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.TagTypeID != rhs.TagTypeID {
		return false
	}
	if row.TagType != rhs.TagType {
		return false
	}
	if row.Description != rhs.Description {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `TagType` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t TagTypeTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row TagType,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.TagTypeID,
		&row.TagType,
		&row.Description)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan TagType: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t TagTypeTable) Table() string {
	return `tag_types`
}

// View returns the table's view (for reading). May be the same as Table().
func (t TagTypeTable) View() string {
	return `tag_types`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t TagTypeTable) SelectCols() string {
	return `dataset_id,tag_type_id,tag_type,description`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_TagType = `INSERT INTO tag_types(
dataset_id,
tag_type_id,
tag_type,
description
) VALUES (
 $1,$2,$3,$4
)`

// Insert will validate and insert a new `TagType`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagTypeTable) Insert(
	db DBi,
	row *TagType,
) (
	err error,
) {

	if row.TagTypeID == "" {
		row.TagTypeID = crypto.NewUUID()
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Validate TagTypeID.
	if err := validate.UUID(row.TagTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagTypeID.")
	}

	// Sanitize TagType.
	row.TagType = sanitize.SingleLineString(row.TagType)

	// Validate TagType.
	if err := validate.TagType(row.TagType); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagType.")
	}

	// Sanitize Description.
	row.Description = sanitize.TrimSpace(row.Description)

	// Execute query.
	_, err = db.Exec(insertQuery_TagType,
		row.DatasetID,
		row.TagTypeID,
		row.TagType,
		row.Description)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagType.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_TagType = `INSERT INTO tag_types(
 dataset_id,
 tag_type_id,
 tag_type,
 description
) VALUES (
 $1,$2,$3,$4
)
ON CONFLICT (dataset_id,tag_type_id)
DO UPDATE SET
 tag_type=EXCLUDED.tag_type,
 description=EXCLUDED.description`

func (t TagTypeTable) Upsert(
	db DBi,
	row *TagType,
) (
	err error,
) {

	if row.TagTypeID == "" {
		row.TagTypeID = crypto.NewUUID()
	}

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Validate TagTypeID.
	if err := validate.UUID(row.TagTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagTypeID.")
	}

	// Sanitize TagType.
	row.TagType = sanitize.SingleLineString(row.TagType)

	// Validate TagType.
	if err := validate.TagType(row.TagType); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagType.")
	}

	// Sanitize Description.
	row.Description = sanitize.TrimSpace(row.Description)

	// Execute query.
	result, err := db.Exec(upsertQuery_TagType,
		row.DatasetID,
		row.TagTypeID,
		row.TagType,
		row.Description)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("TagType update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagType update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagType update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_TagType = `UPDATE
 tag_types
SET
 tag_type=$1,
 description=$2
WHERE
 dataset_id=$3 AND 
 tag_type_id=$4`

// Update updates the following column values:
//   - TagType
//   - Description
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagTypeTable) Update(
	db DBi,
	row *TagType,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Validate TagTypeID.
	if err := validate.UUID(row.TagTypeID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagTypeID.")
	}

	// Sanitize TagType.
	row.TagType = sanitize.SingleLineString(row.TagType)

	// Validate TagType.
	if err := validate.TagType(row.TagType); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagType.")
	}

	// Sanitize Description.
	row.Description = sanitize.TrimSpace(row.Description)

	// Execute query.
	result, err := db.Exec(updateQuery_TagType,
		row.TagType,
		row.Description,
		row.DatasetID,
		row.TagTypeID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("TagType update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagType update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagType update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_TagType = `DELETE FROM
 tag_types
WHERE
 dataset_id=$1 AND 
 tag_type_id=$2`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t TagTypeTable) Delete(
	db DBi,
	DatasetID string,
	TagTypeID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_TagType,
		DatasetID,
		TagTypeID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("TagType.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_TagType = `SELECT
 dataset_id,
 tag_type_id,
 tag_type,
 description
FROM
 tag_types
WHERE
 dataset_id=$1 AND 
 tag_type_id=$2`

// Get returns the `TagType` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t TagTypeTable) Get(
	db DBi,
	DatasetID string,
	TagTypeID string,
) (
	row TagType,
	err error,
) {
	src := db.QueryRow(getQuery_TagType,
		DatasetID,
		TagTypeID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetByTagType
// ----------------------------------------------------------------------------

const getQuery_TagType_byTagType = `SELECT
 dataset_id,
 tag_type_id,
 tag_type,
 description
FROM
  tag_types
WHERE
 dataset_id=$1 AND 
 tag_type=$2`

// GetByTagType return the TagType object by a natural key.
func (t TagTypeTable) GetByTagType(
	db DBi,
	DatasetID string,
	TagType string,
) (
	row TagType,
	err error,
) {

	src := db.QueryRow(getQuery_TagType_byTagType,
		DatasetID,
		TagType)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `TagType` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t TagTypeTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []TagType,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("TagType.List failed: %w", err).
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
			Wrap("TagType.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_TagType = `CREATE TABLE tag_types(
dataset_id,
tag_type_id,
tag_type,
description
)`

func (t TagTypeTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_TagType)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagType.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_TagType)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagType.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagType.List failed: %w", err).
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
			row.TagTypeID,
			row.TagType,
			row.Description,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("TagType.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("TagType.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
