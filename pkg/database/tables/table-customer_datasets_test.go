package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestCustomerDataset_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCustomerDatasetForTesting(db)
	if err := CustomerDatasets.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CustomerDatasets.Get(
		db,
		row.CustomerEntity,
		row.DatasetEntity,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := CustomerDatasets.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestCustomerDataset_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCustomerDatasetForTesting(db)
	if err := CustomerDatasets.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := CustomerDatasets.Delete(
		db,
		row.CustomerEntity, row.DatasetEntity,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = CustomerDatasets.Get(
		db,
		row.CustomerEntity, row.DatasetEntity,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestCustomerDataset_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM customer_datasets WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewCustomerDatasetForTesting(db)
		if err := CustomerDatasets.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := CustomerDatasets.List(
		db,
		`SELECT `+CustomerDatasets.SelectCols()+` FROM `+CustomerDatasets.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
