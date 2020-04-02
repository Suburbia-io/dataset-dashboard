package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/sanitize"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type Location struct {
	DatasetID           string          `json:"datasetID"`
	LocationHash        string          `json:"locationHash"`
	LocationString      string          `json:"locationString"`
	ParsedCountryCode   string          `json:"parsedCountryCode"`
	ParsedPostalCode    string          `json:"parsedPostalCode"`
	GeonamesPostalCodes json.RawMessage `json:"geonamesPostalCodes"`
	GeonameID           *int            `json:"geonameID"`
	GeonamesHierarchy   json.RawMessage `json:"geonamesHierarchy"`
	Approved            *bool           `json:"approved"`
	CreatedAt           time.Time       `json:"createdAt"`
	Name                string          `json:"name"`
	Population          string          `json:"population"`
	CountryCode         string          `json:"countryCode"`
	ParentName          string          `json:"parentName"`
	ParentPopulation    string          `json:"parentPopulation"`
	ParentGeonameID     string          `json:"parentGeonameID"`
}

type LocationTable struct{}

var Locations = LocationTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row Location) Equals(rhs Location) bool {
	if row.DatasetID != rhs.DatasetID {
		return false
	}
	if row.LocationHash != rhs.LocationHash {
		return false
	}
	if row.LocationString != rhs.LocationString {
		return false
	}
	if row.ParsedCountryCode != rhs.ParsedCountryCode {
		return false
	}
	if row.ParsedPostalCode != rhs.ParsedPostalCode {
		return false
	}
	if !bytes.Equal(row.GeonamesPostalCodes, rhs.GeonamesPostalCodes) {
		return false
	}

	if row.GeonameID != nil || rhs.GeonameID != nil {
		if row.GeonameID == nil || rhs.GeonameID == nil {
			return false
		}
		if *row.GeonameID != *rhs.GeonameID {
			return false
		}
	}

	if !bytes.Equal(row.GeonamesHierarchy, rhs.GeonamesHierarchy) {
		return false
	}

	if row.Approved != nil || rhs.Approved != nil {
		if row.Approved == nil || rhs.Approved == nil {
			return false
		}
		if *row.Approved != *rhs.Approved {
			return false
		}
	}

	if !row.CreatedAt.Equal(rhs.CreatedAt) {
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `Location` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t LocationTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row Location,
	err error,
) {
	err = src.Scan(
		&row.DatasetID,
		&row.LocationHash,
		&row.LocationString,
		&row.ParsedCountryCode,
		&row.ParsedPostalCode,
		&row.GeonamesPostalCodes,
		&row.GeonameID,
		&row.GeonamesHierarchy,
		&row.Approved,
		&row.CreatedAt,
		&row.Name,
		&row.Population,
		&row.CountryCode,
		&row.ParentName,
		&row.ParentPopulation,
		&row.ParentGeonameID)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan Location: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t LocationTable) Table() string {
	return `locations`
}

// View returns the table's view (for reading). May be the same as Table().
func (t LocationTable) View() string {
	return `location_view`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t LocationTable) SelectCols() string {
	return `dataset_id,location_hash,location_string,parsed_country_code,parsed_postal_code,geonames_postal_codes,geoname_id,geonames_hierarchy,approved,created_at,name,population,country_code,parent_name,parent_population,parent_geoname_id`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_Location = `INSERT INTO locations(
dataset_id,
location_hash,
location_string,
parsed_country_code,
parsed_postal_code,
geonames_postal_codes,
geoname_id,
geonames_hierarchy,
approved,
created_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
)`

// Insert will validate and insert a new `Location`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t LocationTable) Insert(
	db DBi,
	row *Location,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize LocationHash.
	row.LocationHash = sanitize.SingleLineString(row.LocationHash)

	// Validate LocationHash.
	if err := validate.NonEmptyString(row.LocationHash); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on LocationHash.")
	}

	// Sanitize LocationString.
	row.LocationString = sanitize.SingleLineString(row.LocationString)

	// Validate LocationString.
	if err := validate.NonEmptyString(row.LocationString); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on LocationString.")
	}

	// Execute query.
	_, err = db.Exec(insertQuery_Location,
		row.DatasetID,
		row.LocationHash,
		row.LocationString,
		row.ParsedCountryCode,
		row.ParsedPostalCode,
		row.GeonamesPostalCodes,
		row.GeonameID,
		row.GeonamesHierarchy,
		row.Approved,
		row.CreatedAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Location.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_Location = `INSERT INTO locations(
 dataset_id,
 location_hash,
 location_string,
 parsed_country_code,
 parsed_postal_code,
 geonames_postal_codes,
 geoname_id,
 geonames_hierarchy,
 approved,
 created_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
)
ON CONFLICT (dataset_id,location_hash)
DO UPDATE SET
 parsed_country_code=EXCLUDED.parsed_country_code,
 parsed_postal_code=EXCLUDED.parsed_postal_code,
 geonames_postal_codes=EXCLUDED.geonames_postal_codes,
 geoname_id=EXCLUDED.geoname_id,
 geonames_hierarchy=EXCLUDED.geonames_hierarchy,
 approved=EXCLUDED.approved`

func (t LocationTable) Upsert(
	db DBi,
	row *Location,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize LocationHash.
	row.LocationHash = sanitize.SingleLineString(row.LocationHash)

	// Validate LocationHash.
	if err := validate.NonEmptyString(row.LocationHash); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on LocationHash.")
	}

	// Sanitize LocationString.
	row.LocationString = sanitize.SingleLineString(row.LocationString)

	// Validate LocationString.
	if err := validate.NonEmptyString(row.LocationString); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on LocationString.")
	}

	// Execute query.
	result, err := db.Exec(upsertQuery_Location,
		row.DatasetID,
		row.LocationHash,
		row.LocationString,
		row.ParsedCountryCode,
		row.ParsedPostalCode,
		row.GeonamesPostalCodes,
		row.GeonameID,
		row.GeonamesHierarchy,
		row.Approved,
		row.CreatedAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Location update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Location update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Location update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_Location = `UPDATE
 locations
SET
 parsed_country_code=$1,
 parsed_postal_code=$2,
 geonames_postal_codes=$3,
 geoname_id=$4,
 geonames_hierarchy=$5,
 approved=$6
WHERE
 dataset_id=$7 AND 
 location_hash=$8`

// Update updates the following column values:
//   - ParsedCountryCode
//   - ParsedPostalCode
//   - GeonamesPostalCodes
//   - GeonameID
//   - GeonamesHierarchy
//   - Approved
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t LocationTable) Update(
	db DBi,
	row *Location,
) (
	err error,
) {

	// Validate DatasetID.
	if err := validate.UUID(row.DatasetID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on DatasetID.")
	}

	// Sanitize LocationHash.
	row.LocationHash = sanitize.SingleLineString(row.LocationHash)

	// Validate LocationHash.
	if err := validate.NonEmptyString(row.LocationHash); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on LocationHash.")
	}

	// Sanitize LocationString.
	row.LocationString = sanitize.SingleLineString(row.LocationString)

	// Validate LocationString.
	if err := validate.NonEmptyString(row.LocationString); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on LocationString.")
	}

	// Execute query.
	result, err := db.Exec(updateQuery_Location,
		row.ParsedCountryCode,
		row.ParsedPostalCode,
		row.GeonamesPostalCodes,
		row.GeonameID,
		row.GeonamesHierarchy,
		row.Approved,
		row.DatasetID,
		row.LocationHash)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("Location update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("Location update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("Location update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_Location = `DELETE FROM
 locations
WHERE
 dataset_id=$1 AND 
 location_hash=$2`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t LocationTable) Delete(
	db DBi,
	DatasetID string,
	LocationHash string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_Location,
		DatasetID,
		LocationHash)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("Location.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_Location = `SELECT
 dataset_id,
 location_hash,
 location_string,
 parsed_country_code,
 parsed_postal_code,
 geonames_postal_codes,
 geoname_id,
 geonames_hierarchy,
 approved,
 created_at,
 name,
 population,
 country_code,
 parent_name,
 parent_population,
 parent_geoname_id
FROM
 location_view
WHERE
 dataset_id=$1 AND 
 location_hash=$2`

// Get returns the `Location` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t LocationTable) Get(
	db DBi,
	DatasetID string,
	LocationHash string,
) (
	row Location,
	err error,
) {
	src := db.QueryRow(getQuery_Location,
		DatasetID,
		LocationHash)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetByLocationHash
// ----------------------------------------------------------------------------

const getQuery_Location_byLocationHash = `SELECT
 dataset_id,
 location_hash,
 location_string,
 parsed_country_code,
 parsed_postal_code,
 geonames_postal_codes,
 geoname_id,
 geonames_hierarchy,
 approved,
 created_at,
 name,
 population,
 country_code,
 parent_name,
 parent_population,
 parent_geoname_id
FROM
  location_view
WHERE
 location_hash=$1`

// GetByLocationHash return the Location object by a natural key.
func (t LocationTable) GetByLocationHash(
	db DBi,
	LocationHash string,
) (
	row Location,
	err error,
) {

	src := db.QueryRow(getQuery_Location_byLocationHash,
		LocationHash)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `Location` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t LocationTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []Location,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("Location.List failed: %w", err).
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
			Wrap("Location.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_Location = `CREATE TABLE locations(
dataset_id,
location_hash,
location_string,
parsed_country_code,
parsed_postal_code,
geonames_postal_codes,
geoname_id,
geonames_hierarchy,
approved,
created_at
)`

func (t LocationTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_Location)
	if err != nil {
		return errors.Unexpected.
			Wrap("Location.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_Location)
	if err != nil {
		return errors.Unexpected.
			Wrap("Location.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Location.List failed: %w", err).
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
			row.LocationHash,
			row.LocationString,
			row.ParsedCountryCode,
			row.ParsedPostalCode,
			row.GeonamesPostalCodes,
			row.GeonameID,
			row.GeonamesHierarchy,
			row.Approved,
			row.CreatedAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("Location.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("Location.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
