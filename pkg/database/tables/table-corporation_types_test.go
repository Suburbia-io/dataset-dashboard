package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestCorporationType_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationTypeForTesting(db)
	if err := CorporationTypes.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorporationTypes.Get(
		db,
		row.CorporationTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := CorporationTypes.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestCorporationType_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationTypeForTesting(db)
	if err := CorporationTypes.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorporationTypeForTesting(db)
	row.CorporationType = newRow.CorporationType
	row.Description = newRow.Description

	if err := CorporationTypes.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorporationTypes.Get(
		db,
		row.CorporationTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorporationType_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationTypeForTesting(db)
	if err := CorporationTypes.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorporationTypeForTesting(db)
	row.CorporationType = newRow.CorporationType
	row.Description = newRow.Description

	if err := CorporationTypes.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorporationTypes.Get(
		db,
		row.CorporationTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorporationType_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorporationTypeForTesting(db)
	if err := CorporationTypes.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := CorporationTypes.Delete(
		db,
		row.CorporationTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = CorporationTypes.Get(
		db,
		row.CorporationTypeID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestCorporationType_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM corporation_types WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewCorporationTypeForTesting(db)
		if err := CorporationTypes.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := CorporationTypes.List(
		db,
		`SELECT `+CorporationTypes.SelectCols()+` FROM `+CorporationTypes.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
