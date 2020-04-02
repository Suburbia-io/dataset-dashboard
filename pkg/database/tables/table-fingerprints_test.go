package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestFingerprint_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewFingerprintForTesting(db)
	if err := Fingerprints.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Fingerprints.Get(
		db,
		row.DatasetID,
		row.Fingerprint,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := Fingerprints.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestFingerprint_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewFingerprintForTesting(db)
	if err := Fingerprints.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewFingerprintForTesting(db)
	row.RawText = newRow.RawText
	row.UpdatedAt = newRow.UpdatedAt
	row.Count = newRow.Count

	if err := Fingerprints.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Fingerprints.Get(
		db,
		row.DatasetID,
		row.Fingerprint,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestFingerprint_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewFingerprintForTesting(db)
	if err := Fingerprints.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewFingerprintForTesting(db)
	row.RawText = newRow.RawText
	row.UpdatedAt = newRow.UpdatedAt
	row.Count = newRow.Count

	if err := Fingerprints.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Fingerprints.Get(
		db,
		row.DatasetID,
		row.Fingerprint,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestFingerprint_UpdateAnnotations(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewFingerprintForTesting(db)
	if err := Fingerprints.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.Annotations = NewFingerprintForTesting(db).Annotations

	err := Fingerprints.UpdateAnnotations(
		db,
		row.DatasetID, row.Fingerprint,
		row.Annotations)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Fingerprints.Get(
		db,
		row.DatasetID, row.Fingerprint,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestFingerprint_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewFingerprintForTesting(db)
	if err := Fingerprints.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := Fingerprints.Delete(
		db,
		row.DatasetID, row.Fingerprint,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = Fingerprints.Get(
		db,
		row.DatasetID, row.Fingerprint,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestFingerprint_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM fingerprints WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewFingerprintForTesting(db)
		if err := Fingerprints.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := Fingerprints.List(
		db,
		`SELECT `+Fingerprints.SelectCols()+` FROM `+Fingerprints.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
