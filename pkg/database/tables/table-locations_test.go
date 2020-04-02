package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestLocation_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewLocationForTesting(db)
	if err := Locations.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Locations.Get(
		db,
		row.DatasetID,
		row.LocationHash,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := Locations.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestLocation_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewLocationForTesting(db)
	if err := Locations.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewLocationForTesting(db)
	row.ParsedCountryCode = newRow.ParsedCountryCode
	row.ParsedPostalCode = newRow.ParsedPostalCode
	row.GeonamesPostalCodes = newRow.GeonamesPostalCodes
	row.GeonameID = newRow.GeonameID
	row.GeonamesHierarchy = newRow.GeonamesHierarchy
	row.Approved = newRow.Approved

	if err := Locations.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Locations.Get(
		db,
		row.DatasetID,
		row.LocationHash,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestLocation_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewLocationForTesting(db)
	if err := Locations.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewLocationForTesting(db)
	row.ParsedCountryCode = newRow.ParsedCountryCode
	row.ParsedPostalCode = newRow.ParsedPostalCode
	row.GeonamesPostalCodes = newRow.GeonamesPostalCodes
	row.GeonameID = newRow.GeonameID
	row.GeonamesHierarchy = newRow.GeonamesHierarchy
	row.Approved = newRow.Approved

	if err := Locations.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Locations.Get(
		db,
		row.DatasetID,
		row.LocationHash,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestLocation_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewLocationForTesting(db)
	if err := Locations.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := Locations.Delete(
		db,
		row.DatasetID, row.LocationHash,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = Locations.Get(
		db,
		row.DatasetID, row.LocationHash,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestLocation_GetByLocationHash(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewLocationForTesting(db)
	if err := Locations.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := Locations.GetByLocationHash(
		db,
		row.LocationHash,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestLocation_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM locations WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewLocationForTesting(db)
		if err := Locations.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := Locations.List(
		db,
		`SELECT `+Locations.SelectCols()+` FROM `+Locations.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
