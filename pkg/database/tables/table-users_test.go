package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"testing"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func TestUser_Insert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID,
		row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}

	// Duplicate insert should give DBDuplicate.
	if err := Users.Insert(db, &row); !errors.DBDuplicate.Is(err) {
		t.Fatal(err)
	}
}

func TestUser_Upsert(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewUserForTesting(db)
	row.Name = newRow.Name
	row.Email = newRow.Email
	row.SFTPPubKey = newRow.SFTPPubKey
	row.ArchivedAt = newRow.ArchivedAt

	if err := Users.Upsert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID,
		row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_Update(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Update the row.
	newRow := NewUserForTesting(db)
	row.Name = newRow.Name
	row.Email = newRow.Email
	row.SFTPPubKey = newRow.SFTPPubKey
	row.ArchivedAt = newRow.ArchivedAt

	if err := Users.Update(db, &row); err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID,
		row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateIsRoleAdmin(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.IsRoleAdmin = NewUserForTesting(db).IsRoleAdmin

	err := Users.UpdateIsRoleAdmin(
		db,
		row.CustomerID, row.UserID,
		row.IsRoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateIsRoleSuperAdmin(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.IsRoleSuperAdmin = NewUserForTesting(db).IsRoleSuperAdmin

	err := Users.UpdateIsRoleSuperAdmin(
		db,
		row.CustomerID, row.UserID,
		row.IsRoleSuperAdmin)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateIsRoleCustomerUser(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.IsRoleCustomerUser = NewUserForTesting(db).IsRoleCustomerUser

	err := Users.UpdateIsRoleCustomerUser(
		db,
		row.CustomerID, row.UserID,
		row.IsRoleCustomerUser)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateIsRoleLabeler(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.IsRoleLabeler = NewUserForTesting(db).IsRoleLabeler

	err := Users.UpdateIsRoleLabeler(
		db,
		row.CustomerID, row.UserID,
		row.IsRoleLabeler)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateHash(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.Hash = NewUserForTesting(db).Hash

	err := Users.UpdateHash(
		db,
		row.CustomerID, row.UserID,
		row.Hash)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateAPIKey(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.APIKey = NewUserForTesting(db).APIKey

	err := Users.UpdateAPIKey(
		db,
		row.CustomerID, row.UserID,
		row.APIKey)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateLastActiveAt(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.LastActiveAt = NewUserForTesting(db).LastActiveAt

	err := Users.UpdateLastActiveAt(
		db,
		row.CustomerID, row.UserID,
		row.LastActiveAt)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateLoginToken(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.LoginToken = NewUserForTesting(db).LoginToken

	err := Users.UpdateLoginToken(
		db,
		row.CustomerID, row.UserID,
		row.LoginToken)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_UpdateLoginTokenExpiresAt(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Get random value for field from new row.
	row.LoginTokenExpiresAt = NewUserForTesting(db).LoginTokenExpiresAt

	err := Users.UpdateLoginTokenExpiresAt(
		db,
		row.CustomerID, row.UserID,
		row.LoginTokenExpiresAt)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure row from database matches inserted row.
	row2, err := Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_Delete(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Delete the row.
	err := Users.Delete(
		db,
		row.CustomerID, row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row shouldn't exist.
	_, err = Users.Get(
		db,
		row.CustomerID, row.UserID,
	)
	if !errors.DBNotFound.Is(err) {
		t.Fatal(err)
	}
}

func TestUser_GetByUserID(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := Users.GetByUserID(
		db,
		row.UserID,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_GetByEmail(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := Users.GetByEmail(
		db,
		row.Email,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_GetByAPIKey(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := Users.GetByAPIKey(
		db,
		row.APIKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_GetBySFTPUsername(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	// Insert a row.
	row := NewUserForTesting(db)
	if err := Users.Insert(db, &row); err != nil {
		t.Fatal(err)
	}

	// Row should exist.
	row2, err := Users.GetBySFTPUsername(
		db,
		row.SFTPUsername,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !row.Equals(row2) {
		t.Fatalf("%v != %v", row, row2)
	}
}

func TestUser_List(t *testing.T) {
	t.Parallel()
	db, done := NewDBForTesting()
	defer done()

	_, err := db.Exec(`DELETE FROM users WHERE TRUE`)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		row := NewUserForTesting(db)
		if err := Users.Insert(db, &row); err != nil {
			t.Fatal(err)
		}
	}

	l, err := Users.List(
		db,
		`SELECT `+Users.SelectCols()+` FROM `+Users.View())

	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 10 {
		t.Fatal(l)
	}
}
