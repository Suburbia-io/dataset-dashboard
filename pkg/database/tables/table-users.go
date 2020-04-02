package tables

// ----------------------------------------------------------------------------
// THIS FILE IS GENERATED.
// ----------------------------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/sanitize"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type User struct {
	CustomerID          string     `json:"customerID"`
	UserID              string     `json:"userID"`
	CustomerName        string     `json:"customerName"`
	Name                string     `json:"name"`
	Email               string     `json:"email"`
	IsRoleAdmin         bool       `json:"isRoleAdmin"`
	IsRoleSuperAdmin    bool       `json:"isRoleSuperAdmin"`
	IsRoleCustomerUser  bool       `json:"isRoleCustomerUser"`
	IsRoleLabeler       bool       `json:"isRoleLabeler"`
	Hash                string     `json:"-"`
	APIKey              *string    `json:"apiKey"`
	SFTPUsername        string     `json:"sftpUsername"`
	SFTPPubKey          string     `json:"sftpPubKey"`
	CreatedAt           time.Time  `json:"createdAt"`
	ArchivedAt          *time.Time `json:"archivedAt"`
	LastActiveAt        *time.Time `json:"lastActiveAt"`
	LoginToken          *string    `json:"loginToken"`
	LoginTokenExpiresAt *time.Time `json:"loginTokenExpiresAt"`
}

type UserTable struct{}

var Users = UserTable{}

// ----------------------------------------------------------------------------
// Equality
// ----------------------------------------------------------------------------

// Check equality based on insertable columns. Columns that are view-only are
// ignored.
func (row User) Equals(rhs User) bool {
	if row.CustomerID != rhs.CustomerID {
		return false
	}
	if row.UserID != rhs.UserID {
		return false
	}
	if row.Name != rhs.Name {
		return false
	}
	if row.Email != rhs.Email {
		return false
	}
	if row.IsRoleAdmin != rhs.IsRoleAdmin {
		return false
	}
	if row.IsRoleSuperAdmin != rhs.IsRoleSuperAdmin {
		return false
	}
	if row.IsRoleCustomerUser != rhs.IsRoleCustomerUser {
		return false
	}
	if row.IsRoleLabeler != rhs.IsRoleLabeler {
		return false
	}
	if row.Hash != rhs.Hash {
		return false
	}

	if row.APIKey != nil || rhs.APIKey != nil {
		if row.APIKey == nil || rhs.APIKey == nil {
			return false
		}
		if *row.APIKey != *rhs.APIKey {
			return false
		}
	}

	if row.SFTPUsername != rhs.SFTPUsername {
		return false
	}
	if row.SFTPPubKey != rhs.SFTPPubKey {
		return false
	}
	if !row.CreatedAt.Equal(rhs.CreatedAt) {
		return false
	}

	if row.ArchivedAt != nil || rhs.ArchivedAt != nil {
		if row.ArchivedAt == nil || rhs.ArchivedAt == nil {
			return false
		}
		if !row.ArchivedAt.Equal(*rhs.ArchivedAt) {
			return false
		}
	}

	if row.LastActiveAt != nil || rhs.LastActiveAt != nil {
		if row.LastActiveAt == nil || rhs.LastActiveAt == nil {
			return false
		}
		if !row.LastActiveAt.Equal(*rhs.LastActiveAt) {
			return false
		}
	}

	if row.LoginToken != nil || rhs.LoginToken != nil {
		if row.LoginToken == nil || rhs.LoginToken == nil {
			return false
		}
		if *row.LoginToken != *rhs.LoginToken {
			return false
		}
	}

	if row.LoginTokenExpiresAt != nil || rhs.LoginTokenExpiresAt != nil {
		if row.LoginTokenExpiresAt == nil || rhs.LoginTokenExpiresAt == nil {
			return false
		}
		if !row.LoginTokenExpiresAt.Equal(*rhs.LoginTokenExpiresAt) {
			return false
		}
	}

	return true
}

// ----------------------------------------------------------------------------
// Scan
// ----------------------------------------------------------------------------

// Scan a database row into a `User` object. The selected columns
// should be those returned by the `SelectCols` function.
func (t UserTable) Scan(
	src interface {
		Scan(args ...interface{}) error
	},
) (
	row User,
	err error,
) {
	err = src.Scan(
		&row.CustomerID,
		&row.UserID,
		&row.CustomerName,
		&row.Name,
		&row.Email,
		&row.IsRoleAdmin,
		&row.IsRoleSuperAdmin,
		&row.IsRoleCustomerUser,
		&row.IsRoleLabeler,
		&row.Hash,
		&row.APIKey,
		&row.SFTPUsername,
		&row.SFTPPubKey,
		&row.CreatedAt,
		&row.ArchivedAt,
		&row.LastActiveAt,
		&row.LoginToken,
		&row.LoginTokenExpiresAt)

	if err == nil {
		return
	}

	if err == sql.ErrNoRows {
		return row, errors.DBNotFound
	}

	return row, errors.Unexpected.
		Wrap("Failed to scan User: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Query helpers.
// ----------------------------------------------------------------------------

// Table returns the table name.
func (t UserTable) Table() string {
	return `users`
}

// View returns the table's view (for reading). May be the same as Table().
func (t UserTable) View() string {
	return `user_view`
}

// SelectCols returns a list of columns to select. This should be used when
// building a query in order to use this class's `List` or `Scan` functions.
func (t UserTable) SelectCols() string {
	return `customer_id,user_id,customer_name,name,email,is_role_admin,is_role_super_admin,is_role_customer_user,is_role_labeler,hash,api_key,sftp_username,sftp_pub_key,created_at,archived_at,last_active_at,login_token,login_token_expires_at`
}

// ----------------------------------------------------------------------------
// Insert
// ----------------------------------------------------------------------------

const insertQuery_User = `INSERT INTO users(
customer_id,
user_id,
name,
email,
is_role_admin,
is_role_super_admin,
is_role_customer_user,
is_role_labeler,
hash,
api_key,
sftp_username,
sftp_pub_key,
created_at,
archived_at,
last_active_at,
login_token,
login_token_expires_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17
)`

// Insert will validate and insert a new `User`.
// It may return the following errors:
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) Insert(
	db DBi,
	row *User,
) (
	err error,
) {

	if row.UserID == "" {
		row.UserID = crypto.NewUUID()
	}

	// Validate CustomerID.
	if err := validate.UUID(row.CustomerID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CustomerID.")
	}

	// Validate UserID.
	if err := validate.UUID(row.UserID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on UserID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.HumanName(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Sanitize Email.
	row.Email = sanitize.Email(row.Email)

	// Validate Email.
	if err := validate.Email(row.Email); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Email.")
	}

	// Sanitize Hash.
	row.Hash = sanitize.SingleLineString(row.Hash)

	// Sanitize APIKey.
	row.APIKey = sanitize.SingleLineStringPtr(row.APIKey)

	// Sanitize SFTPUsername.
	row.SFTPUsername = sanitize.SingleLineString(row.SFTPUsername)

	// Sanitize SFTPPubKey.
	row.SFTPPubKey = sanitize.SingleLineString(row.SFTPPubKey)

	// Sanitize LoginToken.
	row.LoginToken = sanitize.SingleLineStringPtr(row.LoginToken)

	// Execute query.
	_, err = db.Exec(insertQuery_User,
		row.CustomerID,
		row.UserID,
		row.Name,
		row.Email,
		row.IsRoleAdmin,
		row.IsRoleSuperAdmin,
		row.IsRoleCustomerUser,
		row.IsRoleLabeler,
		row.Hash,
		row.APIKey,
		row.SFTPUsername,
		row.SFTPPubKey,
		row.CreatedAt,
		row.ArchivedAt,
		row.LastActiveAt,
		row.LoginToken,
		row.LoginTokenExpiresAt)

	if err == nil {
		return nil
	}

	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User.Insert failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Upsert
// ----------------------------------------------------------------------------

const upsertQuery_User = `INSERT INTO users(
 customer_id,
 user_id,
 name,
 email,
 is_role_admin,
 is_role_super_admin,
 is_role_customer_user,
 is_role_labeler,
 hash,
 api_key,
 sftp_username,
 sftp_pub_key,
 created_at,
 archived_at,
 last_active_at,
 login_token,
 login_token_expires_at
) VALUES (
 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17
)
ON CONFLICT (customer_id,user_id)
DO UPDATE SET
 name=EXCLUDED.name,
 email=EXCLUDED.email,
 sftp_pub_key=EXCLUDED.sftp_pub_key,
 archived_at=EXCLUDED.archived_at`

func (t UserTable) Upsert(
	db DBi,
	row *User,
) (
	err error,
) {

	if row.UserID == "" {
		row.UserID = crypto.NewUUID()
	}

	// Validate CustomerID.
	if err := validate.UUID(row.CustomerID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CustomerID.")
	}

	// Validate UserID.
	if err := validate.UUID(row.UserID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on UserID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.HumanName(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Sanitize Email.
	row.Email = sanitize.Email(row.Email)

	// Validate Email.
	if err := validate.Email(row.Email); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Email.")
	}

	// Sanitize Hash.
	row.Hash = sanitize.SingleLineString(row.Hash)

	// Sanitize APIKey.
	row.APIKey = sanitize.SingleLineStringPtr(row.APIKey)

	// Sanitize SFTPUsername.
	row.SFTPUsername = sanitize.SingleLineString(row.SFTPUsername)

	// Sanitize SFTPPubKey.
	row.SFTPPubKey = sanitize.SingleLineString(row.SFTPPubKey)

	// Sanitize LoginToken.
	row.LoginToken = sanitize.SingleLineStringPtr(row.LoginToken)

	// Execute query.
	result, err := db.Exec(upsertQuery_User,
		row.CustomerID,
		row.UserID,
		row.Name,
		row.Email,
		row.IsRoleAdmin,
		row.IsRoleSuperAdmin,
		row.IsRoleCustomerUser,
		row.IsRoleLabeler,
		row.Hash,
		row.APIKey,
		row.SFTPUsername,
		row.SFTPPubKey,
		row.CreatedAt,
		row.ArchivedAt,
		row.LastActiveAt,
		row.LoginToken,
		row.LoginTokenExpiresAt)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

const updateQuery_User = `UPDATE
 users
SET
 name=$1,
 email=$2,
 sftp_pub_key=$3,
 archived_at=$4
WHERE
 customer_id=$5 AND 
 user_id=$6`

// Update updates the following column values:
//   - Name
//   - Email
//   - SFTPPubKey
//   - ArchivedAt
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) Update(
	db DBi,
	row *User,
) (
	err error,
) {

	// Validate CustomerID.
	if err := validate.UUID(row.CustomerID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on CustomerID.")
	}

	// Validate UserID.
	if err := validate.UUID(row.UserID); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on UserID.")
	}

	// Sanitize Name.
	row.Name = sanitize.SingleLineString(row.Name)

	// Validate Name.
	if err := validate.HumanName(row.Name); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Name.")
	}

	// Sanitize Email.
	row.Email = sanitize.Email(row.Email)

	// Validate Email.
	if err := validate.Email(row.Email); err != nil {
		e := err.(errors.Error)
		return e.Wrap("Validation failed on Email.")
	}

	// Sanitize Hash.
	row.Hash = sanitize.SingleLineString(row.Hash)

	// Sanitize APIKey.
	row.APIKey = sanitize.SingleLineStringPtr(row.APIKey)

	// Sanitize SFTPUsername.
	row.SFTPUsername = sanitize.SingleLineString(row.SFTPUsername)

	// Sanitize SFTPPubKey.
	row.SFTPPubKey = sanitize.SingleLineString(row.SFTPPubKey)

	// Sanitize LoginToken.
	row.LoginToken = sanitize.SingleLineStringPtr(row.LoginToken)

	// Execute query.
	result, err := db.Exec(updateQuery_User,
		row.Name,
		row.Email,
		row.SFTPPubKey,
		row.ArchivedAt,
		row.CustomerID,
		row.UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateIsRoleAdmin
// ----------------------------------------------------------------------------

const updateQuery_User_IsRoleAdmin = `UPDATE
 users
SET
 is_role_admin=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateIsRoleAdmin will attempt to update the IsRoleAdmin column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateIsRoleAdmin(
	db DBi,
	CustomerID string,
	UserID string,
	IsRoleAdmin bool,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_User_IsRoleAdmin,
		IsRoleAdmin,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateIsRoleSuperAdmin
// ----------------------------------------------------------------------------

const updateQuery_User_IsRoleSuperAdmin = `UPDATE
 users
SET
 is_role_super_admin=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateIsRoleSuperAdmin will attempt to update the IsRoleSuperAdmin column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateIsRoleSuperAdmin(
	db DBi,
	CustomerID string,
	UserID string,
	IsRoleSuperAdmin bool,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_User_IsRoleSuperAdmin,
		IsRoleSuperAdmin,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateIsRoleCustomerUser
// ----------------------------------------------------------------------------

const updateQuery_User_IsRoleCustomerUser = `UPDATE
 users
SET
 is_role_customer_user=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateIsRoleCustomerUser will attempt to update the IsRoleCustomerUser column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateIsRoleCustomerUser(
	db DBi,
	CustomerID string,
	UserID string,
	IsRoleCustomerUser bool,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_User_IsRoleCustomerUser,
		IsRoleCustomerUser,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateIsRoleLabeler
// ----------------------------------------------------------------------------

const updateQuery_User_IsRoleLabeler = `UPDATE
 users
SET
 is_role_labeler=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateIsRoleLabeler will attempt to update the IsRoleLabeler column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateIsRoleLabeler(
	db DBi,
	CustomerID string,
	UserID string,
	IsRoleLabeler bool,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_User_IsRoleLabeler,
		IsRoleLabeler,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateHash
// ----------------------------------------------------------------------------

const updateQuery_User_Hash = `UPDATE
 users
SET
 hash=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateHash will attempt to update the Hash column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateHash(
	db DBi,
	CustomerID string,
	UserID string,
	Hash string,
) (
	err error,
) {
	Hash = sanitize.SingleLineString(Hash)

	result, err := db.Exec(updateQuery_User_Hash,
		Hash,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateAPIKey
// ----------------------------------------------------------------------------

const updateQuery_User_APIKey = `UPDATE
 users
SET
 api_key=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateAPIKey will attempt to update the APIKey column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateAPIKey(
	db DBi,
	CustomerID string,
	UserID string,
	APIKey *string,
) (
	err error,
) {
	APIKey = sanitize.SingleLineStringPtr(APIKey)

	result, err := db.Exec(updateQuery_User_APIKey,
		APIKey,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateLastActiveAt
// ----------------------------------------------------------------------------

const updateQuery_User_LastActiveAt = `UPDATE
 users
SET
 last_active_at=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateLastActiveAt will attempt to update the LastActiveAt column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateLastActiveAt(
	db DBi,
	CustomerID string,
	UserID string,
	LastActiveAt *time.Time,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_User_LastActiveAt,
		LastActiveAt,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateLoginToken
// ----------------------------------------------------------------------------

const updateQuery_User_LoginToken = `UPDATE
 users
SET
 login_token=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateLoginToken will attempt to update the LoginToken column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateLoginToken(
	db DBi,
	CustomerID string,
	UserID string,
	LoginToken *string,
) (
	err error,
) {
	LoginToken = sanitize.SingleLineStringPtr(LoginToken)

	result, err := db.Exec(updateQuery_User_LoginToken,
		LoginToken,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// UpdateLoginTokenExpiresAt
// ----------------------------------------------------------------------------

const updateQuery_User_LoginTokenExpiresAt = `UPDATE
 users
SET
 login_token_expires_at=$1
WHERE
 customer_id=$2 AND 
 user_id=$3`

// UpdateLoginTokenExpiresAt will attempt to update the LoginTokenExpiresAt column in the row
// corresponding to the given primary key.
//
// It may return the following errors:
//   - DBNotFound
//   - DBDuplicate
//   - DBFKey
//   - DBNullConstraint
// Additionally it may return validation errors.
func (t UserTable) UpdateLoginTokenExpiresAt(
	db DBi,
	CustomerID string,
	UserID string,
	LoginTokenExpiresAt *time.Time,
) (
	err error,
) {

	result, err := db.Exec(updateQuery_User_LoginTokenExpiresAt,
		LoginTokenExpiresAt,
		CustomerID,
		UserID)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			return errors.Unexpected.
				Wrap("User update rows affected failed: %w", err).
				Alert()
		}
		switch n {
		case 1:
			return nil

		case 0:
			return errors.DBNotFound

		default:
			return errors.Unexpected.
				Wrap("User update affected %d rows.", n).
				Alert()
		}
	}

	// Return application error.
	if err := translateDBError(err); err != nil {
		return err
	}

	return errors.Unexpected.
		Wrap("User update failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Delete
// ----------------------------------------------------------------------------

const deleteQuery_User = `DELETE FROM
 users
WHERE
 customer_id=$1 AND 
 user_id=$2`

// Delete attempts to delete a row specified by its primary key. It does not
// cascade, and may return errors.DBFKey if the delete fails.
func (t UserTable) Delete(
	db DBi,
	CustomerID string,
	UserID string,
) (
	err error,
) {
	_, err = db.Exec(deleteQuery_User,
		CustomerID,
		UserID)

	if err == nil {
		return nil
	}
	if err := translateDBError(err); err != nil {
		return err
	}

	// Not a known error.
	return errors.Unexpected.
		Wrap("User.Delete failed: %w", err).
		Alert()
}

// ----------------------------------------------------------------------------
// Get
// ----------------------------------------------------------------------------

const getQuery_User = `SELECT
 customer_id,
 user_id,
 customer_name,
 name,
 email,
 is_role_admin,
 is_role_super_admin,
 is_role_customer_user,
 is_role_labeler,
 hash,
 api_key,
 sftp_username,
 sftp_pub_key,
 created_at,
 archived_at,
 last_active_at,
 login_token,
 login_token_expires_at
FROM
 user_view
WHERE
 customer_id=$1 AND 
 user_id=$2`

// Get returns the `User` object specified by its primary key. May
// return a DBNotFound error if the row isn't found.
func (t UserTable) Get(
	db DBi,
	CustomerID string,
	UserID string,
) (
	row User,
	err error,
) {
	src := db.QueryRow(getQuery_User,
		CustomerID,
		UserID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetByUserID
// ----------------------------------------------------------------------------

const getQuery_User_byUserID = `SELECT
 customer_id,
 user_id,
 customer_name,
 name,
 email,
 is_role_admin,
 is_role_super_admin,
 is_role_customer_user,
 is_role_labeler,
 hash,
 api_key,
 sftp_username,
 sftp_pub_key,
 created_at,
 archived_at,
 last_active_at,
 login_token,
 login_token_expires_at
FROM
  user_view
WHERE
 user_id=$1`

// GetByUserID return the User object by a natural key.
func (t UserTable) GetByUserID(
	db DBi,
	UserID string,
) (
	row User,
	err error,
) {

	src := db.QueryRow(getQuery_User_byUserID,
		UserID)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetByEmail
// ----------------------------------------------------------------------------

const getQuery_User_byEmail = `SELECT
 customer_id,
 user_id,
 customer_name,
 name,
 email,
 is_role_admin,
 is_role_super_admin,
 is_role_customer_user,
 is_role_labeler,
 hash,
 api_key,
 sftp_username,
 sftp_pub_key,
 created_at,
 archived_at,
 last_active_at,
 login_token,
 login_token_expires_at
FROM
  user_view
WHERE
 email=$1`

// GetByEmail return the User object by a natural key.
func (t UserTable) GetByEmail(
	db DBi,
	Email string,
) (
	row User,
	err error,
) {

	src := db.QueryRow(getQuery_User_byEmail,
		Email)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetByAPIKey
// ----------------------------------------------------------------------------

const getQuery_User_byAPIKey = `SELECT
 customer_id,
 user_id,
 customer_name,
 name,
 email,
 is_role_admin,
 is_role_super_admin,
 is_role_customer_user,
 is_role_labeler,
 hash,
 api_key,
 sftp_username,
 sftp_pub_key,
 created_at,
 archived_at,
 last_active_at,
 login_token,
 login_token_expires_at
FROM
  user_view
WHERE
 api_key=$1`

// GetByAPIKey return the User object by a natural key.
func (t UserTable) GetByAPIKey(
	db DBi,
	APIKey *string,
) (
	row User,
	err error,
) {

	if APIKey == nil {
		return row, errors.DBNotFound
	}

	src := db.QueryRow(getQuery_User_byAPIKey,
		APIKey)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// GetBySFTPUsername
// ----------------------------------------------------------------------------

const getQuery_User_bySFTPUsername = `SELECT
 customer_id,
 user_id,
 customer_name,
 name,
 email,
 is_role_admin,
 is_role_super_admin,
 is_role_customer_user,
 is_role_labeler,
 hash,
 api_key,
 sftp_username,
 sftp_pub_key,
 created_at,
 archived_at,
 last_active_at,
 login_token,
 login_token_expires_at
FROM
  user_view
WHERE
 sftp_username=$1`

// GetBySFTPUsername return the User object by a natural key.
func (t UserTable) GetBySFTPUsername(
	db DBi,
	SFTPUsername string,
) (
	row User,
	err error,
) {

	src := db.QueryRow(getQuery_User_bySFTPUsername,
		SFTPUsername)

	return t.Scan(src)
}

// ----------------------------------------------------------------------------
// List
// ----------------------------------------------------------------------------

// List will execute the given query (with arguments) and scan the results into
// a list of `User` objects.
//
// The query should select from the `View` columns returned by the `SelectCols`
// function.
func (t UserTable) List(
	db DBi,
	query string,
	args ...interface{},
) (
	l []User,
	err error,
) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return l, errors.Unexpected.
			Wrap("User.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return l, err
		}
		l = append(l, row)
	}

	if err := rows.Err(); err != nil {
		return l, errors.Unexpected.
			Wrap("User.List iteration failed: %w", err).
			Alert()
	}

	return l, nil
}

// ----------------------------------------------------------------------------
// Dump To Sqlite
// ----------------------------------------------------------------------------

const createTableQuery_User = `CREATE TABLE users(
customer_id,
user_id,
name,
email,
is_role_admin,
is_role_super_admin,
is_role_customer_user,
is_role_labeler,
hash,
api_key,
sftp_username,
sftp_pub_key,
created_at,
archived_at,
last_active_at,
login_token,
login_token_expires_at
)`

func (t UserTable) DumpToSqlite(
	db DBi,
	sqlite DBi,
	selectQuery string,
	args ...interface{},
) (
	err error,
) {
	_, err = sqlite.Exec(createTableQuery_User)
	if err != nil {
		return errors.Unexpected.
			Wrap("User.Create in sqlite failed: %w", err).
			Alert()
	}

	stmt, err := sqlite.Prepare(insertQuery_User)
	if err != nil {
		return errors.Unexpected.
			Wrap("User.Insert to sqlite failed: %w", err).
			Alert()
	}

	rows, err := db.Query(selectQuery, args...)
	if err != nil {
		return errors.Unexpected.
			Wrap("User.List failed: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := t.Scan(rows)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			row.CustomerID,
			row.UserID,
			row.Name,
			row.Email,
			row.IsRoleAdmin,
			row.IsRoleSuperAdmin,
			row.IsRoleCustomerUser,
			row.IsRoleLabeler,
			row.Hash,
			row.APIKey,
			row.SFTPUsername,
			row.SFTPPubKey,
			row.CreatedAt,
			row.ArchivedAt,
			row.LastActiveAt,
			row.LoginToken,
			row.LoginTokenExpiresAt,
		)

		if err != nil {
			return errors.Unexpected.
				Wrap("User.Insert failed: %w", err).
				Alert()
		}
	}

	if err := rows.Err(); err != nil {
		return errors.Unexpected.
			Wrap("User.List iteration failed: %w", err).
			Alert()
	}

	return nil

}
