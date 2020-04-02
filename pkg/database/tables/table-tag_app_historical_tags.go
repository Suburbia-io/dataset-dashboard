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

type TagAppHistoricalTag struct {
	DatasetID   string    `json:"datasetID"`
	Fingerprint string    `json:"fingerprint"`
	TagTypeID   string    `json:"tagTypeID"`
	TagAppID    string    `json:"tagAppID"`
	TagAppName  string    `json:"tagAppName"`
	TagID       string    `json:"tagID"`
	Tag         string    `json:"tag"`
	Confidence  float64   `json:"confidence"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UserID      string    `json:"userID"`
	UserName    string    `json:"userName"`
	UserEmail   string    `json:"userEmail"`
}

type TagAppHistoricalTagTable struct{}

var TagAppHistoricalTags = TagAppHistoricalTagTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row TagAppHistoricalTag) Equals(rhs TagAppHistoricalTag) bool {
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

// Scan a database row into a `TagAppHistoricalTag` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t TagAppHistoricalTagTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row TagAppHistoricalTag,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.Fingerprint,
		&row.TagTypeID,
		&row.TagAppID,
		&row.TagAppName,
		&row.TagID,
		&row.Tag,
		&row.Confidence,
		&row.UpdatedAt,
		&row.UserID,
		&row.UserName,
		&row.UserEmail)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan TagAppHistoricalTag: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t TagAppHistoricalTagTable) Table() string {
	return `tag_app_historical_tags`
}

// View returns the table's view (for reading). May be the same as Table().
func (t TagAppHistoricalTagTable) View() string {
	return `tag_app_historical_tag_view`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t TagAppHistoricalTagTable) SelectCols() string {
	return `dataset_id,fingerprint,tag_type_id,tag_app_id,tag_app_name,tag_id,tag,confidence,updated_at,user_id,user_name,user_email`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_TagAppHistoricalTag = `INSERT INTO tag_app_historical_tags(
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

// Insert will validate and insert a new `TagAppHistoricalTag`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagAppHistoricalTagTable) Insert(
	db DBi,
	row *TagAppHistoricalTag,
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
	_, err = db.Exec(insertQuery_TagAppHistoricalTag,
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
		Wrap("TagAppHistoricalTag.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_TagAppHistoricalTag = `INSERT INTO tag_app_historical_tags(
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
ON CONFLICT (dataset_id,fingerprint,tag_type_id,tag_app_id,updated_at)
DO UPDATE SET
 tag_id=EXCLUDED.tag_id,
 confidence=EXCLUDED.confidence,
 user_id=EXCLUDED.user_id`

func (t TagAppHistoricalTagTable) Upsert(
	db DBi,
	row *TagAppHistoricalTag,
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
	result, err := db.Exec(upsertQuery_TagAppHistoricalTag,
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
				Wrap("TagAppHistoricalTag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagAppHistoricalTag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagAppHistoricalTag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_TagAppHistoricalTag = `UPDATE
 tag_app_historical_tags
SET
 tag_id=$1,
 confidence=$2,
 user_id=$3
WHERE
 dataset_id=$4 AND 
 fingerprint=$5 AND 
 tag_type_id=$6 AND 
 tag_app_id=$7 AND 
 updated_at=$8`

// Update updates the following column values:
//   - TagID
//   - Confidence
//   - UserID
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t TagAppHistoricalTagTable) Update(
	db DBi,
	row *TagAppHistoricalTag,
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
	result, err := db.Exec(updateQuery_TagAppHistoricalTag,
		row.TagID,
		row.Confidence,
		row.UserID,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagAppID,
		row.UpdatedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("TagAppHistoricalTag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("TagAppHistoricalTag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("TagAppHistoricalTag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_TagAppHistoricalTag = `DELETE FROM
 tag_app_historical_tags
WHERE
 dataset_id=$1 AND 
 fingerprint=$2 AND 
 tag_type_id=$3 AND 
 tag_app_id=$4 AND 
 updated_at=$5`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t TagAppHistoricalTagTable) Delete(
	db DBi,
	DatasetID string,
	Fingerprint string,
	TagTypeID string,
	TagAppID string,
	UpdatedAt time.Time,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_TagAppHistoricalTag,
		DatasetID,
		Fingerprint,
		TagTypeID,
		TagAppID,
		UpdatedAt)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("TagAppHistoricalTag.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_TagAppHistoricalTag = `SELECT
 dataset_id,
 fingerprint,
 tag_type_id,
 tag_app_id,
 tag_app_name,
 tag_id,
 tag,
 confidence,
 updated_at,
 user_id,
 user_name,
 user_email
FROM
 tag_app_historical_tag_view
WHERE
 dataset_id=$1 AND 
 fingerprint=$2 AND 
 tag_type_id=$3 AND 
 tag_app_id=$4 AND 
 updated_at=$5`

// Get returns the `TagAppHistoricalTag` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t TagAppHistoricalTagTable) Get(
	db DBi,
	DatasetID string,
	Fingerprint string,
	TagTypeID string,
	TagAppID string,
	UpdatedAt time.Time,
) (
	row TagAppHistoricalTag,
	err error,
) {
	src := db.QueryRow(getQuery_TagAppHistoricalTag,
		DatasetID,
		Fingerprint,
		TagTypeID,
		TagAppID,
		UpdatedAt)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `TagAppHistoricalTag` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t TagAppHistoricalTagTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []TagAppHistoricalTag,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("TagAppHistoricalTag.List failed: %w", err).
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
			Wrap("TagAppHistoricalTag.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_TagAppHistoricalTag = `CREATE TABLE tag_app_historical_tags(
dataset_id,
fingerprint,
tag_type_id,
tag_app_id,
tag_id,
confidence,
updated_at,
user_id
)`

func (t TagAppHistoricalTagTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_TagAppHistoricalTag)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagAppHistoricalTag.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_TagAppHistoricalTag)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagAppHistoricalTag.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("TagAppHistoricalTag.List failed: %w", err).
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
				Wrap("TagAppHistoricalTag.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("TagAppHistoricalTag.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
