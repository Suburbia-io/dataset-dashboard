package application

import (
	"errors"
	"net/http"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminUserList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		ListArgs database.UserListArgs
		Users    []tables.User
		Session  shared.SharedSession
	}{Session: s}
	err = decoder.Decode(&ctx.ListArgs, r.Form)
	if err != nil {
		return err
	}

	// No pagination at the moment.
	ctx.ListArgs.Limit = 100000
	ctx.ListArgs.Offset = 0

	ctx.Users, err = app.DBAL.UserList(ctx.ListArgs)
	if err != nil {
		return err
	}

	err = app.Views.ExecuteSST(w, "admin-user-list.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminUserForm(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Session                      shared.SharedSession
		SuburbiaInternalCustomerUUID string
		User                         struct {
			database.User
			Archived bool
		}
		Customers []database.Customer
	}{
		Session:                      s,
		SuburbiaInternalCustomerUUID: shared.SuburbiaInternalCustomerUUID,
	}
	err = decoder.Decode(&ctx.User, r.Form)
	if err != nil {
		return err
	}

	if ctx.User.UserID != "" {
		ctx.User.User, err = app.DBAL.UserGet(ctx.User.UserID)
		if err != nil {
			return err
		}
	} else {
		archived := false
		listArgs := database.CustomerListArgs{
			// Pagination doesn't make sense here but we might need to do something else if the list becomes too large.
			Limit:    100000,
			Archived: &archived,
		}
		ctx.Customers, err = app.DBAL.CustomerList(listArgs)
		if err != nil {
			return err
		}
	}

	err = app.Views.ExecuteSST(w, "admin-user-form.gohtml", ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) AdminUserFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	args := struct {
		Session shared.SharedSession
		User    struct {
			database.User
			Archived bool
		}
	}{Session: s}
	err = decoder.Decode(&args.User, r.Form)
	if err != nil {
		return err
	}

	if (args.User.ArchivedAt != nil) != args.User.Archived {
		if args.User.Archived {
			t := time.Now()
			args.User.ArchivedAt = &t
		} else {
			args.User.ArchivedAt = nil
		}
	}

	if err = app.DBAL.UserUpsert(&args.User.User); err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(
		s.Session, tables.Users.Table(), args.User.UserID, "upsert", args.User)

	http.Redirect(w, r, "/admin/users/", http.StatusSeeOther)

	return nil
}

func (app *App) UserAPIKeyFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	args := struct {
		Session shared.SharedSession
		APIKey  struct {
			CustomerID string
			UserID     string
			Action     string
		}
	}{Session: s}
	err = decoder.Decode(&args.APIKey, r.Form)
	if err != nil {
		return err
	}

	var apiKey string
	if args.APIKey.Action == "generate" {
		apiKey, err = app.DBAL.UserGenerateAPIKey(args.APIKey.CustomerID, args.APIKey.UserID)
	} else if args.APIKey.Action == "revoke" {
		err = app.DBAL.UserRevokeAPIKey(args.APIKey.CustomerID, args.APIKey.UserID)
	} else {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return nil
	}

	if err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(
		s.Session, tables.Users.Table(), args.APIKey.UserID, "apiKey "+args.APIKey.Action, apiKey)

	http.Redirect(w, r, "/admin/users/"+args.APIKey.UserID, http.StatusSeeOther)

	return nil
}

func (app *App) UserPasswordFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	args := struct {
		Session shared.SharedSession
		Pwd     struct {
			CustomerID string
			UserID     string
			Password   string
		}
	}{Session: s}
	err = decoder.Decode(&args.Pwd, r.Form)
	if err != nil {
		return err
	}

	if args.Pwd.CustomerID != shared.SuburbiaInternalCustomerUUID {
		return errors.New("only users that belong to the Suburbia Internal can have passwords")
	}

	err = app.DBAL.UserSetStrongPasswordForAdmin(args.Pwd.UserID, args.Pwd.Password)
	if err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(
		s.Session, tables.Users.Table(), args.Pwd.UserID, "setPasswordForAdmin", "")

	http.Redirect(w, r, "/admin/users/"+args.Pwd.UserID, http.StatusSeeOther)

	return nil
}

func (app *App) UserRolesFormSubmit(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	args := struct {
		Session shared.SharedSession
		Roles   struct {
			CustomerID         string
			UserID             string
			IsRoleCustomerUser bool
			IsRoleAdmin        bool
			IsRoleSuperAdmin   bool
			IsRoleLabeler      bool
		}
	}{Session: s}
	err = decoder.Decode(&args.Roles, r.Form)
	if err != nil {
		return err
	}

	if args.Roles.IsRoleSuperAdmin && !args.Roles.IsRoleAdmin {
		return errors.New("IsRoleSuperAdmin cannot be set when IsRoleAdmin is not set")
	}

	if args.Roles.IsRoleAdmin && args.Roles.CustomerID != shared.SuburbiaInternalCustomerUUID {
		return errors.New("only users that belong to the Suburbia Internal can be admins")
	}

	if args.Roles.IsRoleLabeler && args.Roles.CustomerID != shared.SuburbiaInternalCustomerUUID {
		return errors.New("only users that belong to the Suburbia Internal can be labelers")
	}

	// This is not the most efficient way to do this but having the roles outside of the normal user update allows us to
	// use the RoleSuperUser permissions check on the route for this handler.
	err = app.DBAL.UserUpdateIsRoleCustomerUser(args.Roles.CustomerID, args.Roles.UserID, args.Roles.IsRoleCustomerUser)
	if err != nil {
		return err
	}
	err = app.DBAL.UserUpdateIsRoleAdmin(args.Roles.CustomerID, args.Roles.UserID, args.Roles.IsRoleAdmin)
	if err != nil {
		return err
	}
	err = app.DBAL.UserUpdateIsRoleSuperAdmin(args.Roles.CustomerID, args.Roles.UserID, args.Roles.IsRoleSuperAdmin)
	if err != nil {
		return err
	}
	err = app.DBAL.UserUpdateIsRoleLabeler(args.Roles.CustomerID, args.Roles.UserID, args.Roles.IsRoleLabeler)
	if err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(
		s.Session, tables.Users.Table(), args.Roles.UserID, "rolesUpdate", args.Roles)

	http.Redirect(w, r, "/admin/users/"+args.Roles.UserID, http.StatusSeeOther)

	return nil
}
