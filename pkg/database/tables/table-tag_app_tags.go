package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type TagAppTag struct {
	DatasetID   string    `json:"datasetID"`
	Fingerprint string    `json:"fingerprint"`
	TagTypeID   string    `json:"tagTypeID"`
	TagAppID    string    `json:"tagAppID"`
	TagID       string    `json:"tagID"`
	Confidence  float64   `json:"confidence"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UserID      string    `json:"userID"`
}

type TagAppTagTable struct{}

var TagAppTags = TagAppTagTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row TagAppTag) Equals(rhs TagAppTag) bool {
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.Fingerprint != rhs.Fingerprint {
		return false
	}
	if row.TagTypeID != rhs.TagTypeID {
		return false
	}
	if row.TagAppID != rhs.TagAppID {
		return false
	}
	if row.TagID != rhs.TagID {
		return false
	}
	if row.Confidence != rhs.Confidence {
		return false
	}
	if !row.UpdatedAt.Equal(rhs.UpdatedAt) {
		return false
	}

	if row.UserID != rhs.UserID {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `TagAppTag` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t TagAppTagTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row TagAppTag,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.Fingerprint,
		&row.TagTypeID,
		&row.TagAppID,
		&row.TagID,
		&row.Confidence,
		&row.UpdatedAt,
		&row.UserID)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan TagAppTag: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t TagAppTagTable) Table() string {
	return `tag_app_tags`
}

// View returns the table's view (for reading). May be the same as Table().
func (t TagAppTagTable) View() string {
	return `tag_app_tags`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t TagAppTagTable) SelectCols() string {
	return `dataset_id,fingerprint,tag_type_id,tag_app_id,tag_id,confidence,updated_at,user_id`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_TagAppTag = `INSERT INTO tag_app_tags(
dataset_id,
fingerprint,
tag_type_id,
tag_app_id,
tag_id,
confidence,
updated_at,
user_id
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8
)`

// Insert will validate and insert a new `TagAppTag`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagAppTagTable) Insert(
	db DBi,
	row *TagAppTag,
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

	// Validate TagAppID.
	if err := validate.UUID(row.TagAppID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagAppID.")
	}

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Validate UserID.
	if err := validate.UUID(row.UserID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on UserID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_TagAppTag,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagAppID,
		row.TagID,
		row.Confidence,
		row.UpdatedAt,
		row.UserID)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagAppTag.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_TagAppTag = `INSERT INTO tag_app_tags(
 dataset_id,
 fingerprint,
 tag_type_id,
 tag_app_id,
 tag_id,
 confidence,
 updated_at,
 user_id
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8
)
ON CONFLICT (dataset_id,fingerprint,tag_type_id,tag_app_id)
DO UPDATE SET
 tag_id=EXCLUDED.tag_id,
 confidence=EXCLUDED.confidence,
 updated_at=EXCLUDED.updated_at,
 user_id=EXCLUDED.user_id`

func (t TagAppTagTable) Upsert(
	db DBi,
	row *TagAppTag,
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

	// Validate TagAppID.
	if err := validate.UUID(row.TagAppID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagAppID.")
	}

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Validate UserID.
	if err := validate.UUID(row.UserID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on UserID.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_TagAppTag,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagAppID,
		row.TagID,
		row.Confidence,
		row.UpdatedAt,
		row.UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("TagAppTag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagAppTag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagAppTag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_TagAppTag = `UPDATE
 tag_app_tags
SET
 tag_id=$1,
 confidence=$2,
 updated_at=$3,
 user_id=$4
WHERE
 dataset_id=$5 AND 
 fingerprint=$6 AND 
 tag_type_id=$7 AND 
 tag_app_id=$8`

// Update updates the following column values:
//   - TagID
//   - Confidence
//   - UpdatedAt
//   - UserID
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagAppTagTable) Update(
	db DBi,
	row *TagAppTag,
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

	// Validate TagAppID.
	if err := validate.UUID(row.TagAppID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagAppID.")
	}

	// Validate TagID.
	if err := validate.UUID(row.TagID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on TagID.")
	}

	// Validate UserID.
	if err := validate.UUID(row.UserID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on UserID.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_TagAppTag,
		row.TagID,
		row.Confidence,
		row.UpdatedAt,
		row.UserID,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagAppID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("TagAppTag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagAppTag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagAppTag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_TagAppTag = `DELETE FROM
 tag_app_tags
WHERE
 dataset_id=$1 AND 
 fingerprint=$2 AND 
 tag_type_id=$3 AND 
 tag_app_id=$4`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t TagAppTagTable) Delete(
	db DBi,
	DatasetID string,
	Fingerprint string,
	TagTypeID string,
	TagAppID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_TagAppTag,
		DatasetID,
		Fingerprint,
		TagTypeID,
		TagAppID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("TagAppTag.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_TagAppTag = `SELECT
 dataset_id,
 fingerprint,
 tag_type_id,
 tag_app_id,
 tag_id,
 confidence,
 updated_at,
 user_id
FROM
 tag_app_tags
WHERE
 dataset_id=$1 AND 
 fingerprint=$2 AND 
 tag_type_id=$3 AND 
 tag_app_id=$4`

// Get returns the `TagAppTag` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t TagAppTagTable) Get(
	db DBi,
	DatasetID string,
	Fingerprint string,
	TagTypeID string,
	TagAppID string,
) (
	row TagAppTag,
	err error,
) {
	src := db.QueryRow(getQuery_TagAppTag,
		DatasetID,
		Fingerprint,
		TagTypeID,
		TagAppID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `TagAppTag` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t TagAppTagTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []TagAppTag,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("TagAppTag.List failed: %w", err).
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
			Wrap("TagAppTag.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_TagAppTag = `CREATE TABLE tag_app_tags(
dataset_id,
fingerprint,
tag_type_id,
tag_app_id,
tag_id,
confidence,
updated_at,
user_id
)`

func (t TagAppTagTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_TagAppTag)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagAppTag.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_TagAppTag)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagAppTag.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagAppTag.List failed: %w", err).
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
			row.TagTypeID,
			row.TagAppID,
			row.TagID,
			row.Confidence,
			row.UpdatedAt,
			row.UserID,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("TagAppTag.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("TagAppTag.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
