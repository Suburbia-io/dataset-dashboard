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

type Tag struct {
	DatasetID       string `json:"datasetID"`
	TagTypeID       string `json:"tagTypeID"`
	TagID           string `json:"tagID"`
	Tag             string `json:"tag"`
	Description     string `json:"description"`
	InternalNotes   string `json:"internalNotes"`
	IsIncluded      bool   `json:"isIncluded"`
	Grade           int    `json:"grade"`
	NumFingerprints int    `json:"numFingerprints"`
	NumLineItems    int    `json:"numLineItems"`
}

type TagTable struct{}

var Tags = TagTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row Tag) Equals(rhs Tag) bool {
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.TagTypeID != rhs.TagTypeID {
		return false
	}
	if row.TagID != rhs.TagID {
		return false
	}
	if row.Tag != rhs.Tag {
		return false
	}
	if row.Description != rhs.Description {
		return false
	}
	if row.InternalNotes != rhs.InternalNotes {
		return false
	}
	if row.IsIncluded != rhs.IsIncluded {
		return false
	}
	if row.Grade != rhs.Grade {
		return false
	}
	if row.NumFingerprints != rhs.NumFingerprints {
		return false
	}
	if row.NumLineItems != rhs.NumLineItems {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `Tag` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t TagTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row Tag,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.TagTypeID,
		&row.TagID,
		&row.Tag,
		&row.Description,
		&row.InternalNotes,
		&row.IsIncluded,
		&row.Grade,
		&row.NumFingerprints,
		&row.NumLineItems)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan Tag: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t TagTable) Table() string {
	return `tags`
}

// View returns the table's view (for reading). May be the same as Table().
func (t TagTable) View() string {
	return `tags`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t TagTable) SelectCols() string {
	return `dataset_id,tag_type_id,tag_id,tag,description,internal_notes,is_included,grade,num_fingerprints,num_line_items`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_Tag = `INSERT INTO tags(
dataset_id,
tag_type_id,
tag_id,
tag,
description,
internal_notes,
is_included,
grade,
num_fingerprints,
num_line_items
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
)`

// Insert will validate and insert a new `Tag`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagTable) Insert(
	db DBi,
	row *Tag,
) (
	err error,
) {

	if row.TagID == "" {
		row.TagID = crypto.NewUUID()
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

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Sanitize Tag.
	row.Tag = sanitize.SingleLineString(row.Tag)

	// Validate Tag.
	if err := validate.Tag(row.Tag); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Tag.")
	}

	// Sanitize Description.
	row.Description = sanitize.TrimSpace(row.Description)

	// Sanitize InternalNotes.
	row.InternalNotes = sanitize.TrimSpace(row.InternalNotes)

	// Validate Grade.
	if err := validate.TagGrade(row.Grade); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Grade.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_Tag,
		row.DatasetID,
		row.TagTypeID,
		row.TagID,
		row.Tag,
		row.Description,
		row.InternalNotes,
		row.IsIncluded,
		row.Grade,
		row.NumFingerprints,
		row.NumLineItems)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Tag.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_Tag = `INSERT INTO tags(
 dataset_id,
 tag_type_id,
 tag_id,
 tag,
 description,
 internal_notes,
 is_included,
 grade,
 num_fingerprints,
 num_line_items
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
)
ON CONFLICT (dataset_id,tag_type_id,tag_id)
DO UPDATE SET
 tag=EXCLUDED.tag,
 description=EXCLUDED.description,
 internal_notes=EXCLUDED.internal_notes,
 is_included=EXCLUDED.is_included,
 grade=EXCLUDED.grade`

func (t TagTable) Upsert(
	db DBi,
	row *Tag,
) (
	err error,
) {

	if row.TagID == "" {
		row.TagID = crypto.NewUUID()
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

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Sanitize Tag.
	row.Tag = sanitize.SingleLineString(row.Tag)

	// Validate Tag.
	if err := validate.Tag(row.Tag); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Tag.")
	}

	// Sanitize Description.
	row.Description = sanitize.TrimSpace(row.Description)

	// Sanitize InternalNotes.
	row.InternalNotes = sanitize.TrimSpace(row.InternalNotes)

	// Validate Grade.
	if err := validate.TagGrade(row.Grade); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Grade.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_Tag,
		row.DatasetID,
		row.TagTypeID,
		row.TagID,
		row.Tag,
		row.Description,
		row.InternalNotes,
		row.IsIncluded,
		row.Grade,
		row.NumFingerprints,
		row.NumLineItems)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Tag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Tag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Tag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_Tag = `UPDATE
 tags
SET
 tag=$1,
 description=$2,
 internal_notes=$3,
 is_included=$4,
 grade=$5
WHERE
 dataset_id=$6 AND 
 tag_type_id=$7 AND 
 tag_id=$8`

// Update updates the following column values:
//   - Tag
//   - Description
//   - InternalNotes
//   - IsIncluded
//   - Grade
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagTable) Update(
	db DBi,
	row *Tag,
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

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Sanitize Tag.
	row.Tag = sanitize.SingleLineString(row.Tag)

	// Validate Tag.
	if err := validate.Tag(row.Tag); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Tag.")
	}

	// Sanitize Description.
	row.Description = sanitize.TrimSpace(row.Description)

	// Sanitize InternalNotes.
	row.InternalNotes = sanitize.TrimSpace(row.InternalNotes)

	// Validate Grade.
	if err := validate.TagGrade(row.Grade); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Grade.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_Tag,
		row.Tag,
		row.Description,
		row.InternalNotes,
		row.IsIncluded,
		row.Grade,
		row.DatasetID,
		row.TagTypeID,
		row.TagID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Tag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Tag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Tag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateNumFingerprints
// ----------------------------------------------------------------------------

const updateQuery_Tag_NumFingerprints = `UPDATE
 tags
SET
 num_fingerprints=$1
WHERE
 dataset_id=$2 AND 
 tag_type_id=$3 AND 
 tag_id=$4`

// UpdateNumFingerprints will attempt to update the NumFingerprints column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagTable) UpdateNumFingerprints(
	db DBi,
	DatasetID string,
	TagTypeID string,
	TagID string,
	NumFingerprints int,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_Tag_NumFingerprints,
		NumFingerprints,
		DatasetID,
		TagTypeID,
		TagID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Tag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Tag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Tag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateNumLineItems
// ----------------------------------------------------------------------------

const updateQuery_Tag_NumLineItems = `UPDATE
 tags
SET
 num_line_items=$1
WHERE
 dataset_id=$2 AND 
 tag_type_id=$3 AND 
 tag_id=$4`

// UpdateNumLineItems will attempt to update the NumLineItems column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagTable) UpdateNumLineItems(
	db DBi,
	DatasetID string,
	TagTypeID string,
	TagID string,
	NumLineItems int,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_Tag_NumLineItems,
		NumLineItems,
		DatasetID,
		TagTypeID,
		TagID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Tag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Tag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Tag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_Tag = `DELETE FROM
 tags
WHERE
 dataset_id=$1 AND 
 tag_type_id=$2 AND 
 tag_id=$3`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t TagTable) Delete(
	db DBi,
	DatasetID string,
	TagTypeID string,
	TagID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_Tag,
		DatasetID,
		TagTypeID,
		TagID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("Tag.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_Tag = `SELECT
 dataset_id,
 tag_type_id,
 tag_id,
 tag,
 description,
 internal_notes,
 is_included,
 grade,
 num_fingerprints,
 num_line_items
FROM
 tags
WHERE
 dataset_id=$1 AND 
 tag_type_id=$2 AND 
 tag_id=$3`

// Get returns the `Tag` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t TagTable) Get(
	db DBi,
	DatasetID string,
	TagTypeID string,
	TagID string,
) (
	row Tag,
	err error,
) {
	src := db.QueryRow(getQuery_Tag,
		DatasetID,
		TagTypeID,
		TagID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetByTag
// ----------------------------------------------------------------------------

const getQuery_Tag_byTag = `SELECT
 dataset_id,
 tag_type_id,
 tag_id,
 tag,
 description,
 internal_notes,
 is_included,
 grade,
 num_fingerprints,
 num_line_items
FROM
  tags
WHERE
 dataset_id=$1 AND 
 tag_type_id=$2 AND 
 tag=$3`

// GetByTag return the Tag object by a natural key.
func (t TagTable) GetByTag(
	db DBi,
	DatasetID string,
	TagTypeID string,
	Tag string,
) (
	row Tag,
	err error,
) {

	src := db.QueryRow(getQuery_Tag_byTag,
		DatasetID,
		TagTypeID,
		Tag)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `Tag` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t TagTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []Tag,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("Tag.List failed: %w", err).
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
			Wrap("Tag.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_Tag = `CREATE TABLE tags(
dataset_id,
tag_type_id,
tag_id,
tag,
description,
internal_notes,
is_included,
grade,
num_fingerprints,
num_line_items
)`

func (t TagTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_Tag)
	if err != nil {
		return errors.Unexpected.
			Wrap("Tag.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_Tag)
	if err != nil {
		return errors.Unexpected.
			Wrap("Tag.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Tag.List failed: %w", err).
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
			row.TagID,
			row.Tag,
			row.Description,
			row.InternalNotes,
			row.IsIncluded,
			row.Grade,
			row.NumFingerprints,
			row.NumLineItems,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("Tag.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("Tag.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
