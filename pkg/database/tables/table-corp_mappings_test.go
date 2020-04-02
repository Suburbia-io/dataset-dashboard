package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestCorpMapping_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingForTesting(db)
	if err := CorpMappings.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorpMappings.Get(
		db,
		row.CorpMappingID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := CorpMappings.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestCorpMapping_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingForTesting(db)
	if err := CorpMappings.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorpMappingForTesting(db)
	row.CorpTypeID = newRow.CorpTypeID
	row.TagTypeID = newRow.TagTypeID
	row.TagID = newRow.TagID

	if err := CorpMappings.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorpMappings.Get(
		db,
		row.CorpMappingID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorpMapping_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingForTesting(db)
	if err := CorpMappings.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorpMappingForTesting(db)
	row.CorpTypeID = newRow.CorpTypeID
	row.TagTypeID = newRow.TagTypeID
	row.TagID = newRow.TagID

	if err := CorpMappings.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorpMappings.Get(
		db,
		row.CorpMappingID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorpMapping_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingForTesting(db)
	if err := CorpMappings.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := CorpMappings.Delete(
		db,
		row.CorpMappingID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = CorpMappings.Get(
		db,
		row.CorpMappingID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestCorpMapping_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM corp_mappings WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewCorpMappingForTesting(db)
		if err := CorpMappings.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := CorpMappings.List(
		db,
		`SELECT `+CorpMappings.SelectCols()+` FROM `+CorpMappings.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
