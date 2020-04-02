package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/sanitize"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"userID"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type SessionTable struct{}

var Sessions = SessionTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row Session) Equals(rhs Session) bool {
	if row.Token != rhs.Token {
		return false
	}
	if row.UserID != rhs.UserID {
		return false
	}
	if !row.ExpiresAt.Equal(rhs.ExpiresAt) {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `Session` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t SessionTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row Session,
	err error,
) {
	err = src.Scan(
		&row.Token,
		&row.UserID,
		&row.ExpiresAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan Session: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t SessionTable) Table() string {
	return `sessions`
}

// View returns the table's view (for reading). May be the same as Table().
func (t SessionTable) View() string {
	return `sessions`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t SessionTable) SelectCols() string {
	return `token,user_id,expires_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_Session = `INSERT INTO sessions(
token,
user_id,
expires_at
) VALUES (
 $1,$2,$3
)`

// Insert will validate and insert a new `Session`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t SessionTable) Insert(
	db DBi,
	row *Session,
) (
	err error,
) {

	// Sanitize Token.
	row.Token = sanitize.SingleLineString(row.Token)

	// Validate Token.
	if err := validate.NonEmptyString(row.Token); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Token.")
	}

	// Validate UserID.
	if err := validate.UUID(row.UserID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on UserID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_Session,
		row.Token,
		row.UserID,
		row.ExpiresAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Session.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_Session = `DELETE FROM
 sessions
WHERE
 token=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t SessionTable) Delete(
	db DBi,
	Token string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_Session,
		Token)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("Session.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_Session = `SELECT
 token,
 user_id,
 expires_at
FROM
 sessions
WHERE
 token=$1`

// Get returns the `Session` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t SessionTable) Get(
	db DBi,
	Token string,
) (
	row Session,
	err error,
) {
	src := db.QueryRow(getQuery_Session,
		Token)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `Session` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t SessionTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []Session,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("Session.List failed: %w", err).
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
			Wrap("Session.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_Session = `CREATE TABLE sessions(
token,
user_id,
expires_at
)`

func (t SessionTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_Session)
	if err != nil {
		return errors.Unexpected.
			Wrap("Session.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_Session)
	if err != nil {
		return errors.Unexpected.
			Wrap("Session.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Session.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.Token,
			row.UserID,
			row.ExpiresAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("Session.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("Session.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
