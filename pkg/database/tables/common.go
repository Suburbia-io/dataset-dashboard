package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED
// ----------------------------------------------------------------------------

import (
	"database/sql"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/lib/pq"
)

type DBi interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Prepare(query string) (*sql.Stmt, error)
}

func translateDBError(err error) error {
	pqErr := &pq.Error{}
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return errors.DBDuplicate

		case "23503":
			return errors.DBFKey

		case "23502":
			return errors.DBNullConstraint
		}
	}
	return nil
}
