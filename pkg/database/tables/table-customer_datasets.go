package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

type CustomerDataset struct {
	CustomerDatasetID string    `json:"customerDatasetID"`
	CustomerEntity    string    `json:"customerEntity"`
	DatasetEntity     string    `json:"datasetEntity"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CustomerDatasetTable struct{}

var CustomerDatasets = CustomerDatasetTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row CustomerDataset) Equals(rhs CustomerDataset) bool {
	if row.CustomerDatasetID != rhs.CustomerDatasetID {
		return false
	}
	if row.CustomerEntity != rhs.CustomerEntity {
		return false
	}
	if row.DatasetEntity != rhs.DatasetEntity {
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

// Scan a database row into a `CustomerDataset` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t CustomerDatasetTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row CustomerDataset,
	err error,
) {
	err = src.Scan(
		&row.CustomerDatasetID,
		&row.CustomerEntity,
		&row.DatasetEntity,
		&row.CreatedAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan CustomerDataset: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t CustomerDatasetTable) Table() string {
	return `customer_datasets`
}

// View returns the table's view (for reading). May be the same as Table().
func (t CustomerDatasetTable) View() string {
	return `customer_datasets`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t CustomerDatasetTable) SelectCols() string {
	return `customer_dataset_id,customer_entity,dataset_entity,created_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_CustomerDataset = `INSERT INTO customer_datasets(
customer_dataset_id,
customer_entity,
dataset_entity,
created_at
) VALUES (
 $1,$2,$3,$4
)`

// Insert will validate and insert a new `CustomerDataset`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CustomerDatasetTable) Insert(
	db DBi,
	row *CustomerDataset,
) (
	err error,
) {

	// Execute query.
	_, err = db.Exec(insertQuery_CustomerDataset,
		row.CustomerDatasetID,
		row.CustomerEntity,
		row.DatasetEntity,
		row.CreatedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CustomerDataset.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_CustomerDataset = `DELETE FROM
 customer_datasets
WHERE
 customer_entity=$1 AND 
 dataset_entity=$2`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t CustomerDatasetTable) Delete(
	db DBi,
	CustomerEntity string,
	DatasetEntity string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_CustomerDataset,
		CustomerEntity,
		DatasetEntity)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("CustomerDataset.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_CustomerDataset = `SELECT
 customer_dataset_id,
 customer_entity,
 dataset_entity,
 created_at
FROM
 customer_datasets
WHERE
 customer_entity=$1 AND 
 dataset_entity=$2`

// Get returns the `CustomerDataset` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t CustomerDatasetTable) Get(
	db DBi,
	CustomerEntity string,
	DatasetEntity string,
) (
	row CustomerDataset,
	err error,
) {
	src := db.QueryRow(getQuery_CustomerDataset,
		CustomerEntity,
		DatasetEntity)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `CustomerDataset` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t CustomerDatasetTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []CustomerDataset,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("CustomerDataset.List failed: %w", err).
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
			Wrap("CustomerDataset.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_CustomerDataset = `CREATE TABLE customer_datasets(
customer_dataset_id,
customer_entity,
dataset_entity,
created_at
)`

func (t CustomerDatasetTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_CustomerDataset)
	if err != nil {
		return errors.Unexpected.
			Wrap("CustomerDataset.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_CustomerDataset)
	if err != nil {
		return errors.Unexpected.
			Wrap("CustomerDataset.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("CustomerDataset.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.CustomerDatasetID,
			row.CustomerEntity,
			row.DatasetEntity,
			row.CreatedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("CustomerDataset.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("CustomerDataset.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
