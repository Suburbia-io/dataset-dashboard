package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestDataset_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Datasets.Get(
		db,
		row.DatasetID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := Datasets.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestDataset_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewDatasetForTesting(db)
	if err := Datasets.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewDatasetForTesting(db)
	row.Name = newRow.Name
	row.Slug = newRow.Slug
	row.Manageable = newRow.Manageable
	row.ArchivedAt = newRow.ArchivedAt

	if err := Datasets.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Datasets.Get(
		db,
		row.DatasetID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestDataset_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewDatasetForTesting(db)
	row.Name = newRow.Name
	row.Slug = newRow.Slug
	row.Manageable = newRow.Manageable
	row.ArchivedAt = newRow.ArchivedAt

	if err := Datasets.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Datasets.Get(
		db,
		row.DatasetID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestDataset_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := Datasets.Delete(
		db,
		row.DatasetID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = Datasets.Get(
		db,
		row.DatasetID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestDataset_GetBySlug(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewDatasetForTesting(db)
	if err := Datasets.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := Datasets.GetBySlug(
		db,
		row.Slug,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestDataset_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM datasets WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewDatasetForTesting(db)
		if err := Datasets.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := Datasets.List(
		db,
		`SELECT `+Datasets.SelectCols()+` FROM `+Datasets.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
