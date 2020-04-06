package application

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/Suburbia-io/dashboard/pkg/errors"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/shared"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/helpers/mailer"
	"github.com/Suburbia-io/dashboard/pkg/views"
)

// -----------------------------------------------------------------------------

func DBService(fresh bool) Service {
	return func(app *App) (err error) {
		app.DBAL, err = database.Bootstrap(app.Config.DB, fresh)
		return err
	}
}

func NewTestDBService(conn *sql.DB) Service {
	return func(app *App) (err error) {
		app.DBAL = &database.DBAL{DB: conn}
		return app.DBAL.Fresh()
	}
}

// -----------------------------------------------------------------------------
func DevSetupService(app *App) (err error) {
	if app.Config.Env != "dev" {
		panic("cannot use DevSetupService in production")
	}

	// Create the Suburbia Internal customer.
	customer, err := app.DBAL.CustomerGet(shared.SuburbiaInternalCustomerUUID)
	if errors.DBNotFound.Is(err) {
		customer.CustomerID = shared.SuburbiaInternalCustomerUUID
		defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Customers.Table(), customer.CustomerID, "suburbiaInternalCreate", customer)
	}
	customer.Name = "Suburbia Internal"
	err = app.DBAL.CustomerUpsert(&customer)
	if err != nil {
		return err
	}

	// Create and setup the admin super user.
	user, err := app.DBAL.UserGetByEmail(app.Config.SuperAdminEmail)
	if errors.DBNotFound.Is(err) {
		user.Name = "Default Admin"
		user.Email = app.Config.SuperAdminEmail
		user.CustomerID = shared.SuburbiaInternalCustomerUUID
		user.ArchivedAt = nil
		err = app.DBAL.UserUpsert(&user)
		if err != nil {
			return err
		}
		defer app.DBAL.AuditTrailBySystemInsertAsync(tables.Users.Table(), user.UserID, "defaultAdminCreate", user)
	}
	err = app.DBAL.UserUpdateIsRoleAdmin(user.CustomerID, user.UserID, true)
	if err != nil {
		return err
	}
	err = app.DBAL.UserUpdateIsRoleSuperAdmin(user.CustomerID, user.UserID, true)
	if err != nil {
		return err
	}
	err = app.DBAL.UserUpdateIsRoleLabeler(user.CustomerID, user.UserID, true)
	if err != nil {
		return err
	}

	if !crypto.CheckPasswordHash(user.Hash, app.Config.SuperAdminPassword) {
		if err := app.DBAL.UserSetStrongPasswordForAdmin(user.UserID, app.Config.SuperAdminPassword); err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
func MailerService(app *App) (err error) {
	app.Mailer = mailer.NewMailer(app.Config.Mailer)
	return nil
}

func MailerTestingService(app *App) (err error) {
	app.Mailer = mailer.NewMailCatcher()
	return nil
}

// -----------------------------------------------------------------------------
func ViewService(app *App) (err error) {
	manifest := map[string]string{}
	manifestBytes, err := ioutil.ReadFile(filepath.Join(app.Config.StaticDir, "manifest.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return err
	}

	app.Views = views.NewViews(app.Config.Views)
	app.Views.RegisterFunc("Manifest", func(key string) string {
		return manifest[key]
	})
	app.Views.RegisterFunc("IsProd", func() bool {
		return app.Config.Env != "dev"
	})

	return nil
}

func ViewTestingService(app *App) (err error) {
	app.Views = views.NewViews(app.Config.Views)
	app.Views.RegisterFunc("Manifest", func(key string) string {
		return key
	})
	app.Views.RegisterFunc("IsProd", func() bool {
		return app.Config.Env == "prod"
	})
	return nil
}
