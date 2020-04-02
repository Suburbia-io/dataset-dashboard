package database

import (
	"strings"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
)

type User struct {
	tables.User
	Customer tables.Customer `schema:"-" json:"customer"`
}

func GenerateSFTPUsername(email string) string {
	return strings.ReplaceAll(email, "@", ".")
}

func (db *DBAL) UserUpsert(user *User) error {
	if user.UserID == "" {
		user.CreatedAt = time.Now()
		user.SFTPUsername = GenerateSFTPUsername(user.Email)
	}
	return tables.Users.Upsert(db, &user.User)
}

func (db *DBAL) UserGet(userID string) (user User, err error) {
	user.User, err = tables.Users.GetByUserID(db, userID)
	if err != nil {
		return user, err
	}
	user.Customer, err = tables.Customers.Get(db, user.CustomerID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *DBAL) UserGetByEmail(email string) (user User, err error) {
	user.User, err = tables.Users.GetByEmail(db, email)
	if err != nil {
		return user, err
	}
	user.Customer, err = tables.Customers.Get(db, user.CustomerID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *DBAL) UserGetByAPIKey(apiKey string) (user User, err error) {
	user.User, err = tables.Users.GetByAPIKey(db, &apiKey)
	if err != nil {
		return user, err
	}
	user.Customer, err = tables.Customers.Get(db, user.CustomerID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *DBAL) UserGetBySFTPUsername(SFTPUsername string) (user User, err error) {
	user.User, err = tables.Users.GetBySFTPUsername(db, SFTPUsername)
	if err != nil {
		return user, err
	}
	user.Customer, err = tables.Customers.Get(db, user.CustomerID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *DBAL) UserGenerateAPIKey(customerID, userID string) (apiKey string, err error) {
	apiKey = crypto.RandAlphaNum(32)

	err = tables.Users.UpdateAPIKey(db, customerID, userID, &apiKey)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

func (db *DBAL) UserRevokeAPIKey(customerID, userID string) (err error) {
	err = tables.Users.UpdateAPIKey(db, customerID, userID, nil)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBAL) UserCreateLoginToken(customerID, userID string, lifetime time.Duration) (loginToken string, err error) {
	loginToken = crypto.RandAlphaNum(32)
	expiresAt := time.Now().Add(lifetime)

	err = tables.Users.UpdateLoginToken(db, customerID, userID, &loginToken)
	if err != nil {
		return "", err
	}

	err = tables.Users.UpdateLoginTokenExpiresAt(db, customerID, userID, &expiresAt)
	if err != nil {
		return "", err
	}

	return loginToken, nil
}

func (db *DBAL) UserClearLoginToken(customerID string, userID string) (err error) {
	err = tables.Users.UpdateLoginToken(db, customerID, userID, nil)
	if err != nil {
		return err
	}

	err = tables.Users.UpdateLoginTokenExpiresAt(db, customerID, userID, nil)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBAL) UserGetByValidLoginToken(loginToken string) (user User, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Users.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Users.View())
	builder.Write(`WHERE archived_at IS NULL`)
	builder.Write(`AND login_token = $1`, loginToken)
	builder.Write(`AND login_token_expires_at > NOW()`)

	query, queryArgs := builder.MustBuild()
	row := db.QueryRow(query, queryArgs...)

	user.User, err = tables.Users.Scan(row)
	if err != nil {
		return user, err
	}

	user.Customer, err = tables.Customers.Get(db, user.CustomerID)
	if err != nil {
		return user, err
	}

	return user, nil
}

type UserListArgs struct {
	Search       string `json:"search" schema:"search"`
	Archived     *bool  `json:"archived" schema:"archived"`
	IsAdmin      *bool  `json:"isAdmin" schema:"is-admin"`
	IsSuperAdmin *bool  `json:"isSuperAdmin" schema:"is-super-admin"`
	IsLabeler    *bool  `json:"isLabeler" schema:"is-labeler"`
	Limit        int    `json:"limit" schema:"limit"`
	Offset       int    `json:"offset" schema:"offset"`
}

func (db *DBAL) UserList(args UserListArgs) (users []tables.User, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Users.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Users.View())
	builder.Write(`WHERE TRUE`)

	if args.Archived != nil {
		if *args.Archived {
			builder.Write(`AND archived_at IS NOT NULL`)
		} else {
			builder.Write(`AND archived_at IS NULL`)
		}
	}

	if args.IsAdmin != nil {
		if *args.IsAdmin {
			builder.Write(`AND is_role_admin IS TRUE`)
		} else {
			builder.Write(`AND is_role_admin IS FALSE`)
		}
	}

	if args.IsSuperAdmin != nil {
		if *args.IsSuperAdmin {
			builder.Write(`AND is_role_super_admin IS TRUE`)
		} else {
			builder.Write(`AND is_role_super_admin IS FALSE`)
		}
	}

	if args.IsLabeler != nil {
		if *args.IsLabeler {
			builder.Write(`AND is_role_labeler IS TRUE`)
		} else {
			builder.Write(`AND is_role_labeler IS FALSE`)
		}
	}

	if args.Search != "" {
		search := "%" + args.Search + "%"
		builder.Write(`AND (email ILIKE $1 OR name ILIKE $1 OR customer_name ILIKE $1)`, search)
	}

	builder.Write(`ORDER BY name ASC`)
	builder.Write(`LIMIT $1 OFFSET $2`, args.Limit, args.Offset)

	query, queryArgs := builder.MustBuild()
	return tables.Users.List(db, query, queryArgs...)
}

func (db *DBAL) UserTouchLastActiveAt(customerID, userID string) (err error) {
	now := time.Now()
	return tables.Users.UpdateLastActiveAt(db, customerID, userID, &now)
}

func (db *DBAL) UserSetStrongPasswordForAdmin(userID, password string) (err error) {
	user, err := db.UserGet(userID)
	if err != nil {
		return errors.DBNotFound
	}

	err = validate.StrongPassword(password, []string{user.Name, user.Email, user.UserID})
	if err != nil {
		return err
	}

	hash, err := crypto.HashPassword(password)
	if err != nil {
		return errors.UnexpectedError(err, "Failed hashing password")
	}

	return tables.Users.UpdateHash(db, user.CustomerID, user.UserID, hash)
}

func (db *DBAL) UserUpdateIsRoleCustomerUser(customerID, userID string, isRoleCustomerUser bool) (err error) {
	return tables.Users.UpdateIsRoleCustomerUser(db, customerID, userID, isRoleCustomerUser)
}

func (db *DBAL) UserUpdateIsRoleAdmin(customerID, userID string, isRoleAdmin bool) (err error) {
	return tables.Users.UpdateIsRoleAdmin(db, customerID, userID, isRoleAdmin)
}

func (db *DBAL) UserUpdateIsRoleSuperAdmin(customerID, userID string, isRoleSuperAdmin bool) (err error) {
	return tables.Users.UpdateIsRoleSuperAdmin(db, customerID, userID, isRoleSuperAdmin)
}

func (db *DBAL) UserUpdateIsRoleLabeler(customerID, userID string, isRoleLabeler bool) (err error) {
	return tables.Users.UpdateIsRoleLabeler(db, customerID, userID, isRoleLabeler)
}
