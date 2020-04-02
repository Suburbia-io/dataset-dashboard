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

type CorpMappingRule struct {
	CorpMappingRuleID string    `json:"corpMappingRuleID"`
	CorpMappingID     string    `json:"corpMappingID"`
	CorpID            string    `json:"corpID"`
	InternalNotes     string    `json:"internalNotes"`
	ExternalNotes     string    `json:"externalNotes"`
	FromDate          time.Time `json:"fromDate"`
	Country           string    `json:"country"`
}

type CorpMappingRuleTable struct{}

var CorpMappingRules = CorpMappingRuleTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row CorpMappingRule) Equals(rhs CorpMappingRule) bool {
	if row.CorpMappingRuleID != rhs.CorpMappingRuleID {
		return false
	}
	if row.CorpMappingID != rhs.CorpMappingID {
		return false
	}
	if row.CorpID != rhs.CorpID {
		return false
	}
	if row.InternalNotes != rhs.InternalNotes {
		return false
	}
	if row.ExternalNotes != rhs.ExternalNotes {
		return false
	}
	if !row.FromDate.Equal(rhs.FromDate) {
		return false
	}

	if row.Country != rhs.Country {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `CorpMappingRule` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t CorpMappingRuleTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row CorpMappingRule,
	err error,
) {
	err = src.Scan(
		&row.CorpMappingRuleID,
		&row.CorpMappingID,
		&row.CorpID,
		&row.InternalNotes,
		&row.ExternalNotes,
		&row.FromDate,
		&row.Country)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan CorpMappingRule: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t CorpMappingRuleTable) Table() string {
	return `corp_mapping_rules`
}

// View returns the table's view (for reading). May be the same as Table().
func (t CorpMappingRuleTable) View() string {
	return `corp_mapping_rules`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t CorpMappingRuleTable) SelectCols() string {
	return `corp_mapping_rule_id,corp_mapping_id,corp_id,internal_notes,external_notes,from_date,country`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_CorpMappingRule = `INSERT INTO corp_mapping_rules(
corp_mapping_rule_id,
corp_mapping_id,
corp_id,
internal_notes,
external_notes,
from_date,
country
) VALUES (
 $1,$2,$3,$4,$5,$6,$7
)`

// Insert will validate and insert a new `CorpMappingRule`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorpMappingRuleTable) Insert(
	db DBi,
	row *CorpMappingRule,
) (
	err error,
) {

	if row.CorpMappingRuleID == "" {
		row.CorpMappingRuleID = crypto.NewUUID()
	}

	// Validate CorpMappingRuleID.
	if err := validate.UUID(row.CorpMappingRuleID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingRuleID.")
	}

	// Validate CorpMappingID.
	if err := validate.UUID(row.CorpMappingID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingID.")
	}

	// Validate CorpID.
	if err := validate.UUID(row.CorpID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpID.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_CorpMappingRule,
		row.CorpMappingRuleID,
		row.CorpMappingID,
		row.CorpID,
		row.InternalNotes,
		row.ExternalNotes,
		row.FromDate,
		row.Country)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorpMappingRule.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_CorpMappingRule = `INSERT INTO corp_mapping_rules(
 corp_mapping_rule_id,
 corp_mapping_id,
 corp_id,
 internal_notes,
 external_notes,
 from_date,
 country
) VALUES (
 $1,$2,$3,$4,$5,$6,$7
)
ON CONFLICT (corp_mapping_rule_id)
DO UPDATE SET
 corp_mapping_id=EXCLUDED.corp_mapping_id,
 corp_id=EXCLUDED.corp_id,
 internal_notes=EXCLUDED.internal_notes,
 external_notes=EXCLUDED.external_notes,
 from_date=EXCLUDED.from_date,
 country=EXCLUDED.country`

func (t CorpMappingRuleTable) Upsert(
	db DBi,
	row *CorpMappingRule,
) (
	err error,
) {

	if row.CorpMappingRuleID == "" {
		row.CorpMappingRuleID = crypto.NewUUID()
	}

	// Validate CorpMappingRuleID.
	if err := validate.UUID(row.CorpMappingRuleID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingRuleID.")
	}

	// Validate CorpMappingID.
	if err := validate.UUID(row.CorpMappingID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingID.")
	}

	// Validate CorpID.
	if err := validate.UUID(row.CorpID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpID.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_CorpMappingRule,
		row.CorpMappingRuleID,
		row.CorpMappingID,
		row.CorpID,
		row.InternalNotes,
		row.ExternalNotes,
		row.FromDate,
		row.Country)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("CorpMappingRule update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("CorpMappingRule update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorpMappingRule update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_CorpMappingRule = `UPDATE
 corp_mapping_rules
SET
 corp_mapping_id=$1,
 corp_id=$2,
 internal_notes=$3,
 external_notes=$4,
 from_date=$5,
 country=$6
WHERE
 corp_mapping_rule_id=$7`

// Update updates the following column values:
//   - CorpMappingID
//   - CorpID
//   - InternalNotes
//   - ExternalNotes
//   - FromDate
//   - Country
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t CorpMappingRuleTable) Update(
	db DBi,
	row *CorpMappingRule,
) (
	err error,
) {

	// Validate CorpMappingRuleID.
	if err := validate.UUID(row.CorpMappingRuleID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingRuleID.")
	}

	// Validate CorpMappingID.
	if err := validate.UUID(row.CorpMappingID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpMappingID.")
	}

	// Validate CorpID.
	if err := validate.UUID(row.CorpID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CorpID.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_CorpMappingRule,
		row.CorpMappingID,
		row.CorpID,
		row.InternalNotes,
		row.ExternalNotes,
		row.FromDate,
		row.Country,
		row.CorpMappingRuleID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("CorpMappingRule update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("CorpMappingRule update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("CorpMappingRule update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_CorpMappingRule = `DELETE FROM
 corp_mapping_rules
WHERE
 corp_mapping_rule_id=$1`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t CorpMappingRuleTable) Delete(
	db DBi,
	CorpMappingRuleID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_CorpMappingRule,
		CorpMappingRuleID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("CorpMappingRule.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_CorpMappingRule = `SELECT
 corp_mapping_rule_id,
 corp_mapping_id,
 corp_id,
 internal_notes,
 external_notes,
 from_date,
 country
FROM
 corp_mapping_rules
WHERE
 corp_mapping_rule_id=$1`

// Get returns the `CorpMappingRule` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t CorpMappingRuleTable) Get(
	db DBi,
	CorpMappingRuleID string,
) (
	row CorpMappingRule,
	err error,
) {
	src := db.QueryRow(getQuery_CorpMappingRule,
		CorpMappingRuleID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `CorpMappingRule` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t CorpMappingRuleTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []CorpMappingRule,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("CorpMappingRule.List failed: %w", err).
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
			Wrap("CorpMappingRule.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_CorpMappingRule = `CREATE TABLE corp_mapping_rules(
corp_mapping_rule_id,
corp_mapping_id,
corp_id,
internal_notes,
external_notes,
from_date,
country
)`

func (t CorpMappingRuleTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_CorpMappingRule)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorpMappingRule.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_CorpMappingRule)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorpMappingRule.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("CorpMappingRule.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.CorpMappingRuleID,
			row.CorpMappingID,
			row.CorpID,
			row.InternalNotes,
			row.ExternalNotes,
			row.FromDate,
			row.Country,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("CorpMappingRule.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("CorpMappingRule.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
