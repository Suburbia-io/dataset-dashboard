package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestSession_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewSessionForTesting(db)
	if err := Sessions.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Sessions.Get(
		db,
		row.Token,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := Sessions.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestSession_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewSessionForTesting(db)
	if err := Sessions.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := Sessions.Delete(
		db,
		row.Token,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = Sessions.Get(
		db,
		row.Token,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestSession_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM sessions WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewSessionForTesting(db)
		if err := Sessions.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := Sessions.List(
		db,
		`SELECT `+Sessions.SelectCols()+` FROM `+Sessions.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
