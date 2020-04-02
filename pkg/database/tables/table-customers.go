package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/sanitize"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type Customer struct {
	CustomerID string     `json:"customerID"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"createdAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

type CustomerTable struct{}

var Customers = CustomerTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row Customer) Equals(rhs Customer) bool {
	if row.CustomerID != rhs.CustomerID {
		return false
	}
	if row.Name != rhs.Name {
		return false
	}
	if !row.CreatedAt.Equal(rhs.CreatedAt) {
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

// Scan a database row into a `Customer` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t CustomerTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row Customer,
	err error,
) {
	err = src.Scan(
		&row.CustomerID,
		&row.Name,
		&row.CreatedAt,
		&row.ArchivedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan Customer: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t CustomerTable) Table() string {
	return `customers`
}

// View returns the table's view (for reading). May be the same as Table().
func (t CustomerTable) View() string {
	return `customers`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t CustomerTable) SelectCols() string {
	return `customer_id,name,created_at,archived_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_Customer = `INSERT INTO customers(
customer_id,
name,
created_at,
archived_at
) VALUES (
 $1,$2,$3,$4
)`

// Insert will validate and insert a new `Customer`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CustomerTable) Insert(
	db DBi,
	row *Customer,
) (
	err error,
) {

	if row.CustomerID == "" {
		row.CustomerID = crypto.NewUUID()
	}

	// Validate CustomerID.
	if err := validate.UUID(row.CustomerID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CustomerID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.NonEmptyString(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_Customer,
		row.CustomerID,
		row.Name,
		row.CreatedAt,
		row.ArchivedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Customer.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_Customer = `INSERT INTO customers(
 customer_id,
 name,
 created_at,
 archived_at
) VALUES (
 $1,$2,$3,$4
)
ON CONFLICT (customer_id)
DO UPDATE SET
 name=EXCLUDED.name,
 archived_at=EXCLUDED.archived_at`

func (t CustomerTable) Upsert(
	db DBi,
	row *Customer,
) (
	err error,
) {

	if row.CustomerID == "" {
		row.CustomerID = crypto.NewUUID()
	}

	// Validate CustomerID.
	if err := validate.UUID(row.CustomerID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CustomerID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.NonEmptyString(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_Customer,
		row.CustomerID,
		row.Name,
		row.CreatedAt,
		row.ArchivedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Customer update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Customer update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Customer update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_Customer = `UPDATE
 customers
SET
 name=$1,
 archived_at=$2
WHERE
 customer_id=$3`

// Update updates the following column values:
//   - Name
//   - ArchivedAt
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CustomerTable) Update(
	db DBi,
	row *Customer,
) (
	err error,
) {

	// Validate CustomerID.
	if err := validate.UUID(row.CustomerID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CustomerID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.NonEmptyString(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_Customer,
		row.Name,
		row.ArchivedAt,
		row.CustomerID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Customer update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Customer update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Customer update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_Customer = `DELETE FROM
 customers
WHERE
 customer_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t CustomerTable) Delete(
	db DBi,
	CustomerID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_Customer,
		CustomerID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("Customer.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_Customer = `SELECT
 customer_id,
 name,
 created_at,
 archived_at
FROM
 customers
WHERE
 customer_id=$1`

// Get returns the `Customer` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t CustomerTable) Get(
	db DBi,
	CustomerID string,
) (
	row Customer,
	err error,
) {
	src := db.QueryRow(getQuery_Customer,
		CustomerID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `Customer` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t CustomerTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []Customer,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("Customer.List failed: %w", err).
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
			Wrap("Customer.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_Customer = `CREATE TABLE customers(
customer_id,
name,
created_at,
archived_at
)`

func (t CustomerTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_Customer)
	if err != nil {
		return errors.Unexpected.
			Wrap("Customer.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_Customer)
	if err != nil {
		return errors.Unexpected.
			Wrap("Customer.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Customer.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.CustomerID,
			row.Name,
			row.CreatedAt,
			row.ArchivedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("Customer.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("Customer.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
