package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestConsensusTag_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewConsensusTagForTesting(db)
	if err := ConsensusTags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := ConsensusTags.Get(
		db,
		row.DatasetID,
		row.Fingerprint,
		row.TagTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := ConsensusTags.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestConsensusTag_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewConsensusTagForTesting(db)
	if err := ConsensusTags.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := ConsensusTags.Delete(
		db,
		row.DatasetID, row.Fingerprint, row.TagTypeID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = ConsensusTags.Get(
		db,
		row.DatasetID, row.Fingerprint, row.TagTypeID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestConsensusTag_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM consensus_tags WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewConsensusTagForTesting(db)
		if err := ConsensusTags.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := ConsensusTags.List(
		db,
		`SELECT `+ConsensusTags.SelectCols()+` FROM `+ConsensusTags.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
