package application

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (app *App) Routes() http.Handler {
	// ----------------------------------------------------------------------------
	// Server-side rendered Admin routes
	// ----------------------------------------------------------------------------
	// Use a sub-router per section until the old Vue.js based dashboard is removed.
	// eventually move this to app.Mount()
	csrfMiddleware := csrf.Protect([]byte(app.Config.CSRFSecret), csrf.Path("/"))
	if app.Config.Env == "dev" {
		csrfMiddleware = csrf.Protect([]byte(app.Config.CSRFSecret), csrf.Path("/"), csrf.Secure(false))
	}

	adminRouter := app.Router.PathPrefix("/admin").Subrouter().StrictSlash(true)
	adminRouter.NewRoute().Path("/").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/auth/", http.StatusSeeOther)
	})
	adminRouter.Use(csrfMiddleware)
	adminAPIRouter := app.Router.PathPrefix("/admin/api/").Subrouter().StrictSlash(true)

	// Static file server for assets
	publicStatic := http.FileServer(http.Dir(app.Config.StaticDir))
	app.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static", publicStatic))

	// Admin Authentication Routes
	app.adminAuthPage(adminRouter, "/auth/", "admin-auth", http.MethodGet, app.AdminAuth)
	app.adminAuthPage(adminRouter, "/auth/", "admin-auth-submit", http.MethodPost, app.AdminAuthSubmit)
	app.adminAuthPage(adminRouter, "/auth/two-factor/", "admin-two-factor", http.MethodGet, app.AdminTwoFactor)
	app.adminAuthPage(adminRouter, "/auth/two-factor/", "admin-two-factor-submit", http.MethodPost, app.AdminTwoFactorSubmit)

	// User roles are checked in the handler for logout.
	app.pageWrapper(adminRouter, "/logout/", "logout", http.MethodPost, Public, app.Logout)

	// Admin User routes
	app.adminPage(adminRouter, "/users/", "user-list", http.MethodGet, app.AdminUserList)
	app.adminPage(adminRouter, "/users/create/", "user-create", http.MethodGet, app.AdminUserForm)
	app.adminPage(adminRouter, "/users/submit/", "user-form-submit", http.MethodPost, app.AdminUserFormSubmit)
	app.adminPage(adminRouter, "/users/{UserID}/", "user-edit", http.MethodGet, app.AdminUserForm)
	app.pageWrapper(adminRouter, "/users/api-key-submit/", "user-api-key-form-submit", http.MethodPost, RoleSuperAdmin, app.UserAPIKeyFormSubmit)
	app.pageWrapper(adminRouter, "/users/roles-submit/", "user-roles-form-submit", http.MethodPost, RoleSuperAdmin, app.UserRolesFormSubmit)
	app.pageWrapper(adminRouter, "/users/password-submit/", "user-password-form-submit", http.MethodPost, RoleSuperAdmin, app.UserPasswordFormSubmit)

	// Admin Dataset routes
	app.pageWrapper(adminRouter, "/datasets/", "dataset-list", http.MethodGet, RoleAdmin|RoleLabeler, app.AdminDatasetList)
	app.adminPage(adminRouter, "/datasets/create/", "dataset-create", http.MethodGet, app.AdminDatasetForm)
	app.adminPage(adminRouter, "/datasets/submit/", "dataset-form-submit", http.MethodPost, app.AdminDatasetFormSubmit)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/", "dataset-edit", http.MethodGet, RoleAdmin|RoleLabeler, app.AdminDatasetForm)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/dump", "dataset-dump-form", http.MethodGet, RoleAdmin|RoleSuperAdmin, app.AdminDumpDatasetForm)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/dump-submit", "dataset-dump-submit", http.MethodGet, RoleAdmin|RoleSuperAdmin, app.AdminDumpDatasetSubmit)

	// Admin Dataset Tag Type routes
	app.adminPage(adminRouter, "/datasets/{DatasetID}/tag-types/create/", "tag-type-create", http.MethodGet, app.AdminTagTypeForm)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/", "tag-type-edit", http.MethodGet, app.AdminTagTypeForm)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/tag-types/submit/", "tag-type-submit", http.MethodPost, app.AdminTagTypeSubmit)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/delete/", "tag-type-delete", http.MethodPost, app.AdminTagTypeDelete)

	// Admin Dataset Tag Type Tags routes
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/tags/", "tag-list", http.MethodGet, RoleAdmin|RoleLabeler, app.AdminTagList)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/tags/create/", "tag-create", http.MethodGet, RoleAdmin|RoleLabeler, app.AdminTagForm)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/update-tag-counts/", "update-tag-counts",
		http.MethodPost, RoleAdmin|RoleLabeler, app.AdminUpdateTagCountsFormSubmit)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/tags/{TagID}/", "tag-edit", http.MethodGet, RoleAdmin|RoleLabeler, app.AdminTagForm)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/tags/submit/", "tag-submit", http.MethodPost, RoleAdmin|RoleLabeler, app.AdminTagSubmit)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/tag-types/{TagTypeID}/tags/{TagID}/delete/", "tag-delete", http.MethodPost, RoleAdmin|RoleLabeler, app.AdminTagDelete)

	// Admin Dataset Corporation Type routes
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporation-types/create/", "corporation-type-create", http.MethodGet, app.AdminCorporationTypeForm)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporation-types/{CorporationTypeID}/", "corporation-type-edit", http.MethodGet, app.AdminCorporationTypeForm)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporation-types/submit/", "corporation-type-submit", http.MethodPost, app.AdminCorporationTypeSubmit)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporation-types/{CorporationTypeID}/delete/", "corporation-type-delete", http.MethodPost, app.AdminCorporationTypeDelete)

	// Admin Dataset Corporation routes
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporations/create/", "corporation-create", http.MethodGet, app.AdminCorporationForm)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporations/{CorporationID}/", "corporation-edit", http.MethodGet, app.AdminCorporationForm)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporations/submit/", "corporation-submit", http.MethodPost, app.AdminCorporationSubmit)
	app.adminPage(adminRouter, "/datasets/{DatasetID}/corporations/{CorporationID}/delete/", "corporation-delete", http.MethodPost, app.AdminCorporationDelete)

	// Admin Dataset Location routes
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/locations/", "location-list", http.MethodGet, RoleLabeler, app.AdminLocationList)
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/locations/search-geonames/", "location-search-geonames",
		http.MethodPost, RoleLabeler, app.AdminLocationSearchGeonamesFormSubmit)
	app.apiWrapper(adminAPIRouter, "/location-approve/", "location-approve", RoleLabeler, app.AdminApiLocationApprove)
	app.apiWrapper(adminAPIRouter, "/location-set-geoname-id/", "location-set-geoname-id", RoleLabeler, app.AdminApiLocationSetGeonameID)
	app.adminAPI(adminAPIRouter, "/locationUpsertCSV", "location-upsert-csv", app.AdminApiLocationUpsertCSV)

	// Admin Dataset CorpMapping routes
	app.pageWrapper(adminRouter, "/datasets/{DatasetID}/corpmappings/{CorpTypeID}", "corp-mapping-list", http.MethodGet, RoleAdmin|RoleSuperAdmin, app.AdminCorpMappings)
	app.apiWrapper(adminAPIRouter, "/listTagsForTagType", "tag-api-list-tags", RoleAdmin|RoleSuperAdmin, app.AdminApiListTagsForTagType)
	app.apiWrapper(adminAPIRouter, "/insertCorpMapping", "tag-api-insert-corpmapping", RoleAdmin|RoleSuperAdmin, app.AdminApiInsertCorpMapping)
	app.apiWrapper(adminAPIRouter, "/insertCorpMappingRule", "tag-api-insert-corpmapping-rule", RoleAdmin|RoleSuperAdmin, app.AdminApiInsertCorpMappingRule)
	app.apiWrapper(adminAPIRouter, "/deleteCorpMappingRule", "tag-api-delete-corpmapping-rule", RoleAdmin|RoleSuperAdmin, app.AdminApiDeleteCorpMappingRule)

	// Admin Customer routes
	app.adminPage(adminRouter, "/customers/", "customer-list", http.MethodGet, app.AdminCustomerList)
	app.adminPage(adminRouter, "/customers/create/", "customer-create", http.MethodGet, app.AdminCustomerForm)
	app.adminPage(adminRouter, "/customers/submit/", "customer-edit", http.MethodPost, app.AdminCustomerFormSubmit)
	app.adminPage(adminRouter, "/customers/{CustomerID}/", "customer-edit", http.MethodGet, app.AdminCustomerForm)

	// Admin Audit Trail routes
	app.adminPage(adminRouter, "/audit-trail/", "audittrail-list", http.MethodGet, app.AdminAuditTrailList)

	// Admin Fingerprints routes
	app.pageWrapper(adminRouter, "/fingerprints/{DatasetID}/", "fingerprint-list", http.MethodGet, RoleLabeler, app.AdminFingerprintList)

	app.apiWrapper(adminAPIRouter, "/fingerprintUpsertTags", "fingerprint-api-upsert-tags", RoleLabeler, app.AdminApiFingerprintUpsertTags)
	app.apiWrapper(adminAPIRouter, "/fingerprintListTags", "fingerprint-api-list-tags", RoleLabeler, app.AdminApiFingerprintListTags)
	app.apiWrapper(adminAPIRouter, "/fingerprintList", "fingerprint-api-list", RoleLabeler, app.AdminApiFingerprintList)
	app.apiWrapper(adminAPIRouter, "/fingerprintTagSuggestions", "fingerprint-api-tag-suggestions", RoleLabeler, app.AdminApiFingerprintTagSuggestions)
	app.apiWrapper(adminAPIRouter, "/fingerprintsListTagApps", "fingerprint-api-list-tagapps", RoleLabeler, app.AdminApiFingerprintListTagApps)

	app.adminAPI(adminAPIRouter, "/fingerprintUpsertCSV", "fingerprint-upsert-csv", app.AdminApiFingerprintUpsertCSV)
	app.adminAPI(adminAPIRouter, "/fingerprintAnnotationsUpsertCSV", "fingerprint-annotations-upsert-csv", app.AdminApiFingerprintAnnotationsUpsertCSV)
	app.adminAPI(adminAPIRouter, "/fingerprintTagUpsertCSV", "fingerprint-tag-upsert-csv", app.AdminApiFingerprintTagUpsertCSV)

	// --------------------------------------------------------------------------
	// Customer user pages
	// --------------------------------------------------------------------------
	customerUserPages := app.Router.PathPrefix("/").Subrouter().StrictSlash(true)
	customerUserPages.NewRoute().Path("/").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	})
	customerUserPages.Use(csrfMiddleware)

	// Customer user login
	app.pageWrapper(customerUserPages, "/login/", "customer-user-login", http.MethodGet, Public, app.CustomerUserLogin)
	app.pageWrapper(customerUserPages, "/login/submit/", "customer-user-login-submit", http.MethodPost, Public,
		app.CustomerUserLoginSubmit)

	// Customer user dashboard
	app.pageWrapper(customerUserPages, "/dashboard/", "customer-user-login", http.MethodGet, RoleCustomerUser,
		app.CustomerUserDashboard)

	// --------------------------------------------------------------------------
	// Legacy dataset specific export (v1 export)
	// See /datasets/{DatasetID}/dump for the new export which is a bit faster.
	// This is still here for backwards compatibility for tagapps that may use this old export
	// --------------------------------------------------------------------------
	app.pageWrapper(adminRouter, "/export/v1/dataset", "dataset-export-v1", http.MethodGet, RoleAdmin|RoleLabeler, app.AdminAPIExportDataset_v1)

	return middlewareGroup(
		app.catchPanicMdl,
		app.logMdl,
		app.secureHeadersMdl,
	)(app.Router)
}
