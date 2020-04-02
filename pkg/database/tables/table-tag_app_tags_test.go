package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestTagAppTag_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagAppTagForTesting(db)
	if err := TagAppTags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := TagAppTags.Get(
		db,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagAppID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := TagAppTags.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestTagAppTag_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagAppTagForTesting(db)
	if err := TagAppTags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := TagAppTags.Delete(
		db,
		row.DatasetID, row.Fingerprint, row.TagTypeID, row.TagAppID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = TagAppTags.Get(
		db,
		row.DatasetID, row.Fingerprint, row.TagTypeID, row.TagAppID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestTagAppTag_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM tag_app_tags WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewTagAppTagForTesting(db)
		if err := TagAppTags.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := TagAppTags.List(
		db,
		`SELECT `+TagAppTags.SelectCols()+` FROM `+TagAppTags.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
