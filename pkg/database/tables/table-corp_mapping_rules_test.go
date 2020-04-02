package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestCorpMappingRule_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingRuleForTesting(db)
	if err := CorpMappingRules.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorpMappingRules.Get(
		db,
		row.CorpMappingRuleID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := CorpMappingRules.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestCorpMappingRule_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingRuleForTesting(db)
	if err := CorpMappingRules.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorpMappingRuleForTesting(db)
	row.CorpMappingID = newRow.CorpMappingID
	row.CorpID = newRow.CorpID
	row.InternalNotes = newRow.InternalNotes
	row.ExternalNotes = newRow.ExternalNotes
	row.FromDate = newRow.FromDate
	row.Country = newRow.Country

	if err := CorpMappingRules.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorpMappingRules.Get(
		db,
		row.CorpMappingRuleID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorpMappingRule_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingRuleForTesting(db)
	if err := CorpMappingRules.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewCorpMappingRuleForTesting(db)
	row.CorpMappingID = newRow.CorpMappingID
	row.CorpID = newRow.CorpID
	row.InternalNotes = newRow.InternalNotes
	row.ExternalNotes = newRow.ExternalNotes
	row.FromDate = newRow.FromDate
	row.Country = newRow.Country

	if err := CorpMappingRules.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := CorpMappingRules.Get(
		db,
		row.CorpMappingRuleID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestCorpMappingRule_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewCorpMappingRuleForTesting(db)
	if err := CorpMappingRules.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := CorpMappingRules.Delete(
		db,
		row.CorpMappingRuleID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = CorpMappingRules.Get(
		db,
		row.CorpMappingRuleID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestCorpMappingRule_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM corp_mapping_rules WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewCorpMappingRuleForTesting(db)
		if err := CorpMappingRules.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := CorpMappingRules.List(
		db,
		`SELECT `+CorpMappingRules.SelectCols()+` FROM `+CorpMappingRules.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
