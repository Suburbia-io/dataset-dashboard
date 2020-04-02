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

type AuditTrail struct {
	AuditTrailID string    `json:"auditTrailID"`
	ByUser       *string   `json:"byUser"`
	Type         string    `json:"type"`
	RelatedTable string    `json:"relatedTable"`
	RelatedID    string    `json:"relatedID"`
	Payload      string    `json:"payload"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AuditTrailTable struct{}

var AuditTrails = AuditTrailTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row AuditTrail) Equals(rhs AuditTrail) bool {
	if row.AuditTrailID != rhs.AuditTrailID {
		return false
	}

	if row.ByUser != nil || rhs.ByUser != nil {
		if row.ByUser == nil || rhs.ByUser == nil {
			return false
		}
		if *row.ByUser != *rhs.ByUser {
			return false
		}
	}

	if row.Type != rhs.Type {
		return false
	}
	if row.RelatedTable != rhs.RelatedTable {
		return false
	}
	if row.RelatedID != rhs.RelatedID {
		return false
	}
	if row.Payload != rhs.Payload {
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

// Scan a database row into a `AuditTrail` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t AuditTrailTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row AuditTrail,
	err error,
) {
	err = src.Scan(
		&row.AuditTrailID,
		&row.ByUser,
		&row.Type,
		&row.RelatedTable,
		&row.RelatedID,
		&row.Payload,
		&row.CreatedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan AuditTrail: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t AuditTrailTable) Table() string {
	return `audit_trails`
}

// View returns the table's view (for reading). May be the same as Table().
func (t AuditTrailTable) View() string {
	return `audit_trails`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t AuditTrailTable) SelectCols() string {
	return `audit_trail_id,by_user,type,related_table,related_id,payload,created_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_AuditTrail = `INSERT INTO audit_trails(
audit_trail_id,
by_user,
type,
related_table,
related_id,
payload,
created_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7
)`

// Insert will validate and insert a new `AuditTrail`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t AuditTrailTable) Insert(
	db DBi,
	row *AuditTrail,
) (
	err error,
) {

	if row.AuditTrailID == "" {
		row.AuditTrailID = crypto.NewUUID()
	}

	// Validate ByUser.
	if err := validate.UUIDPtr(row.ByUser); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on ByUser.")
	}

	// Validate RelatedID.
	if err := validate.UUID(row.RelatedID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on RelatedID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_AuditTrail,
		row.AuditTrailID,
		row.ByUser,
		row.Type,
		row.RelatedTable,
		row.RelatedID,
		row.Payload,
		row.CreatedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("AuditTrail.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_AuditTrail = `DELETE FROM
 audit_trails
WHERE
 audit_trail_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t AuditTrailTable) Delete(
	db DBi,
	AuditTrailID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_AuditTrail,
		AuditTrailID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("AuditTrail.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_AuditTrail = `SELECT
 audit_trail_id,
 by_user,
 type,
 related_table,
 related_id,
 payload,
 created_at
FROM
 audit_trails
WHERE
 audit_trail_id=$1`

// Get returns the `AuditTrail` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t AuditTrailTable) Get(
	db DBi,
	AuditTrailID string,
) (
	row AuditTrail,
	err error,
) {
	src := db.QueryRow(getQuery_AuditTrail,
		AuditTrailID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `AuditTrail` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t AuditTrailTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []AuditTrail,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("AuditTrail.List failed: %w", err).
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
			Wrap("AuditTrail.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_AuditTrail = `CREATE TABLE audit_trails(
audit_trail_id,
by_user,
type,
related_table,
related_id,
payload,
created_at
)`

func (t AuditTrailTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_AuditTrail)
	if err != nil {
		return errors.Unexpected.
			Wrap("AuditTrail.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_AuditTrail)
	if err != nil {
		return errors.Unexpected.
			Wrap("AuditTrail.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("AuditTrail.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.AuditTrailID,
			row.ByUser,
			row.Type,
			row.RelatedTable,
			row.RelatedID,
			row.Payload,
			row.CreatedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("AuditTrail.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("AuditTrail.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
