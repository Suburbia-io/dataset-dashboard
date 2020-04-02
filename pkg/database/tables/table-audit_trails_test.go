package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestAuditTrail_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewAuditTrailForTesting(db)
	if err := AuditTrails.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := AuditTrails.Get(
		db,
		row.AuditTrailID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := AuditTrails.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestAuditTrail_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewAuditTrailForTesting(db)
	if err := AuditTrails.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := AuditTrails.Delete(
		db,
		row.AuditTrailID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = AuditTrails.Get(
		db,
		row.AuditTrailID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestAuditTrail_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM audit_trails WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewAuditTrailForTesting(db)
		if err := AuditTrails.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := AuditTrails.List(
		db,
		`SELECT `+AuditTrails.SelectCols()+` FROM `+AuditTrails.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
