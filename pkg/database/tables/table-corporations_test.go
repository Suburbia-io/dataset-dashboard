package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestCorporation_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationForTesting(db)
	if err := Corporations.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Corporations.Get(
		db,
		row.CorporationID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := Corporations.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestCorporation_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationForTesting(db)
	if err := Corporations.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorporationForTesting(db)
	row.Exchange = newRow.Exchange
	row.Code = newRow.Code
	row.Name = newRow.Name
	row.Slug = newRow.Slug
	row.Isin = newRow.Isin
	row.Cusip = newRow.Cusip

	if err := Corporations.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Corporations.Get(
		db,
		row.CorporationID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorporation_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationForTesting(db)
	if err := Corporations.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorporationForTesting(db)
	row.Exchange = newRow.Exchange
	row.Code = newRow.Code
	row.Name = newRow.Name
	row.Slug = newRow.Slug
	row.Isin = newRow.Isin
	row.Cusip = newRow.Cusip

	if err := Corporations.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Corporations.Get(
		db,
		row.CorporationID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorporation_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationForTesting(db)
	if err := Corporations.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := Corporations.Delete(
		db,
		row.CorporationID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = Corporations.Get(
		db,
		row.CorporationID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestCorporation_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM corporations WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewCorporationForTesting(db)
		if err := Corporations.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := Corporations.List(
		db,
		`SELECT `+Corporations.SelectCols()+` FROM `+Corporations.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
