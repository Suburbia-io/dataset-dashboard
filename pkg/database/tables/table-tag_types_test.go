package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestTagType_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagTypeForTesting(db)
	if err := TagTypes.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := TagTypes.Get(
		db,
		row.DatasetID,
		row.TagTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := TagTypes.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestTagType_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagTypeForTesting(db)
	if err := TagTypes.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewTagTypeForTesting(db)
	row.TagType = newRow.TagType
	row.Description = newRow.Description

	if err := TagTypes.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := TagTypes.Get(
		db,
		row.DatasetID,
		row.TagTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTagType_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagTypeForTesting(db)
	if err := TagTypes.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewTagTypeForTesting(db)
	row.TagType = newRow.TagType
	row.Description = newRow.Description

	if err := TagTypes.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := TagTypes.Get(
		db,
		row.DatasetID,
		row.TagTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTagType_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagTypeForTesting(db)
	if err := TagTypes.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := TagTypes.Delete(
		db,
		row.DatasetID, row.TagTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = TagTypes.Get(
		db,
		row.DatasetID, row.TagTypeID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestTagType_GetByTagType(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagTypeForTesting(db)
	if err := TagTypes.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := TagTypes.GetByTagType(
		db,
		row.DatasetID, row.TagType,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTagType_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM tag_types WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewTagTypeForTesting(db)
		if err := TagTypes.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := TagTypes.List(
		db,
		`SELECT `+TagTypes.SelectCols()+` FROM `+TagTypes.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
