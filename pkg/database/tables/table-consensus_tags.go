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

type ConsensusTag struct {
	DatasetID   string    `json:"datasetID"`
	Fingerprint string    `json:"fingerprint"`
	TagTypeID   string    `json:"tagTypeID"`
	TagID       string    `json:"tagID"`
	Confidence  float64   `json:"confidence"`
	SourceCount int64     `json:"sourceCount"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ConsensusTagTable struct{}

var ConsensusTags = ConsensusTagTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row ConsensusTag) Equals(rhs ConsensusTag) bool {
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.Fingerprint != rhs.Fingerprint {
		return false
	}
	if row.TagTypeID != rhs.TagTypeID {
		return false
	}
	if row.TagID != rhs.TagID {
		return false
	}
	if row.Confidence != rhs.Confidence {
		return false
	}
	if row.SourceCount != rhs.SourceCount {
		return false
	}
	if !row.UpdatedAt.Equal(rhs.UpdatedAt) {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `ConsensusTag` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t ConsensusTagTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row ConsensusTag,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.Fingerprint,
		&row.TagTypeID,
		&row.TagID,
		&row.Confidence,
		&row.SourceCount,
		&row.UpdatedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan ConsensusTag: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t ConsensusTagTable) Table() string {
	return `consensus_tags`
}

// View returns the table's view (for reading). May be the same as Table().
func (t ConsensusTagTable) View() string {
	return `consensus_tags`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t ConsensusTagTable) SelectCols() string {
	return `dataset_id,fingerprint,tag_type_id,tag_id,confidence,source_count,updated_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_ConsensusTag = `INSERT INTO consensus_tags(
dataset_id,
fingerprint,
tag_type_id,
tag_id,
confidence,
source_count,
updated_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7
)`

// Insert will validate and insert a new `ConsensusTag`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t ConsensusTagTable) Insert(
	db DBi,
	row *ConsensusTag,
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

	// Validate Confidence.
	if err := validate.Confidence(row.Confidence); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Confidence.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_ConsensusTag,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagID,
		row.Confidence,
		row.SourceCount,
		row.UpdatedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("ConsensusTag.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_ConsensusTag = `INSERT INTO consensus_tags(
 dataset_id,
 fingerprint,
 tag_type_id,
 tag_id,
 confidence,
 source_count,
 updated_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7
)
ON CONFLICT (dataset_id,fingerprint,tag_type_id)
DO UPDATE SET
 tag_id=EXCLUDED.tag_id,
 confidence=EXCLUDED.confidence,
 source_count=EXCLUDED.source_count,
 updated_at=EXCLUDED.updated_at`

func (t ConsensusTagTable) Upsert(
	db DBi,
	row *ConsensusTag,
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

	// Validate Confidence.
	if err := validate.Confidence(row.Confidence); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Confidence.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_ConsensusTag,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagID,
		row.Confidence,
		row.SourceCount,
		row.UpdatedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("ConsensusTag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("ConsensusTag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("ConsensusTag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_ConsensusTag = `UPDATE
 consensus_tags
SET
 tag_id=$1,
 confidence=$2,
 source_count=$3,
 updated_at=$4
WHERE
 dataset_id=$5 AND 
 fingerprint=$6 AND 
 tag_type_id=$7`

// Update updates the following column values:
//   - TagID
//   - Confidence
//   - SourceCount
//   - UpdatedAt
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t ConsensusTagTable) Update(
	db DBi,
	row *ConsensusTag,
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

	// Validate Confidence.
	if err := validate.Confidence(row.Confidence); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Confidence.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_ConsensusTag,
		row.TagID,
		row.Confidence,
		row.SourceCount,
		row.UpdatedAt,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("ConsensusTag update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("ConsensusTag update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("ConsensusTag update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_ConsensusTag = `DELETE FROM
 consensus_tags
WHERE
 dataset_id=$1 AND 
 fingerprint=$2 AND 
 tag_type_id=$3`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t ConsensusTagTable) Delete(
	db DBi,
	DatasetID string,
	Fingerprint string,
	TagTypeID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_ConsensusTag,
		DatasetID,
		Fingerprint,
		TagTypeID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("ConsensusTag.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_ConsensusTag = `SELECT
 dataset_id,
 fingerprint,
 tag_type_id,
 tag_id,
 confidence,
 source_count,
 updated_at
FROM
 consensus_tags
WHERE
 dataset_id=$1 AND 
 fingerprint=$2 AND 
 tag_type_id=$3`

// Get returns the `ConsensusTag` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t ConsensusTagTable) Get(
	db DBi,
	DatasetID string,
	Fingerprint string,
	TagTypeID string,
) (
	row ConsensusTag,
	err error,
) {
	src := db.QueryRow(getQuery_ConsensusTag,
		DatasetID,
		Fingerprint,
		TagTypeID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `ConsensusTag` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t ConsensusTagTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []ConsensusTag,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("ConsensusTag.List failed: %w", err).
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
			Wrap("ConsensusTag.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_ConsensusTag = `CREATE TABLE consensus_tags(
dataset_id,
fingerprint,
tag_type_id,
tag_id,
confidence,
source_count,
updated_at
)`

func (t ConsensusTagTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_ConsensusTag)
	if err != nil {
		return errors.Unexpected.
			Wrap("ConsensusTag.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_ConsensusTag)
	if err != nil {
		return errors.Unexpected.
			Wrap("ConsensusTag.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("ConsensusTag.List failed: %w", err).
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
			row.TagID,
			row.Confidence,
			row.SourceCount,
			row.UpdatedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("ConsensusTag.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("ConsensusTag.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
