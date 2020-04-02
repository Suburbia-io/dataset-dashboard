package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestTag_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagForTesting(db)
	if err := Tags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Tags.Get(
		db,
		row.DatasetID,
		row.TagTypeID,
		row.TagID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := Tags.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestTag_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagForTesting(db)
	if err := Tags.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewTagForTesting(db)
	row.Tag = newRow.Tag
	row.Description = newRow.Description
	row.InternalNotes = newRow.InternalNotes
	row.IsIncluded = newRow.IsIncluded
	row.Grade = newRow.Grade

	if err := Tags.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Tags.Get(
		db,
		row.DatasetID,
		row.TagTypeID,
		row.TagID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTag_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagForTesting(db)
	if err := Tags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewTagForTesting(db)
	row.Tag = newRow.Tag
	row.Description = newRow.Description
	row.InternalNotes = newRow.InternalNotes
	row.IsIncluded = newRow.IsIncluded
	row.Grade = newRow.Grade

	if err := Tags.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Tags.Get(
		db,
		row.DatasetID,
		row.TagTypeID,
		row.TagID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTag_UpdateNumFingerprints(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagForTesting(db)
	if err := Tags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.NumFingerprints = NewTagForTesting(db).NumFingerprints

	err := Tags.UpdateNumFingerprints(
		db,
		row.DatasetID, row.TagTypeID, row.TagID,
		row.NumFingerprints)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Tags.Get(
		db,
		row.DatasetID, row.TagTypeID, row.TagID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTag_UpdateNumLineItems(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagForTesting(db)
	if err := Tags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.NumLineItems = NewTagForTesting(db).NumLineItems

	err := Tags.UpdateNumLineItems(
		db,
		row.DatasetID, row.TagTypeID, row.TagID,
		row.NumLineItems)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Tags.Get(
		db,
		row.DatasetID, row.TagTypeID, row.TagID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTag_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagForTesting(db)
	if err := Tags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := Tags.Delete(
		db,
		row.DatasetID, row.TagTypeID, row.TagID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = Tags.Get(
		db,
		row.DatasetID, row.TagTypeID, row.TagID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestTag_GetByTag(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagForTesting(db)
	if err := Tags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := Tags.GetByTag(
		db,
		row.DatasetID, row.TagTypeID, row.Tag,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestTag_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM tags WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewTagForTesting(db)
		if err := Tags.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := Tags.List(
		db,
		`SELECT `+Tags.SelectCols()+` FROM `+Tags.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
