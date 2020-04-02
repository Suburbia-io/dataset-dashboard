package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestTagAppHistoricalTag_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagAppHistoricalTagForTesting(db)
	if err := TagAppHistoricalTags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := TagAppHistoricalTags.Get(
		db,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
		row.TagAppID,
		row.UpdatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := TagAppHistoricalTags.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestTagAppHistoricalTag_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewTagAppHistoricalTagForTesting(db)
	if err := TagAppHistoricalTags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := TagAppHistoricalTags.Delete(
		db,
		row.DatasetID, row.Fingerprint, row.TagTypeID, row.TagAppID, row.UpdatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = TagAppHistoricalTags.Get(
		db,
		row.DatasetID, row.Fingerprint, row.TagTypeID, row.TagAppID, row.UpdatedAt,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestTagAppHistoricalTag_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM tag_app_historical_tags WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewTagAppHistoricalTagForTesting(db)
		if err := TagAppHistoricalTags.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := TagAppHistoricalTags.List(
		db,
		`SELECT `+TagAppHistoricalTags.SelectCols()+` FROM `+TagAppHistoricalTags.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
