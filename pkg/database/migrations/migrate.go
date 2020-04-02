package migrations

import (
	"database/sql"

	"github.com/pkg/errors"
)

var Migrations = []string{
	CreateAdminsTable,
	CreateAuthTokensTable,
	CreateAdminSessionsTable,
	CreateUsersTable,
	CreateUserSessionsTable,
	CreateBrandsTable,
	CreateIXRulesTable,
	CreateTickersTable,
	CreateCorporationsTable,
	CreateEventsTable,
	AlterTickersCorporationsTable,
	CreateCorpMappingTable,
	AlterCorpsTableAddISINCUSIP,
	AlterBrandsTableAddPublished,
	AlterUserTableAddHash,
	AlterUserTableAddApiKey,
	AlterAdminTableAddApiKey,
	AlterUserTableMakeApiKeyUnique,
	AlterEventsAddAuthMethod,
	AlterAuthTokensAddSalt,
	CreateLocationsTable,
	AlterLocationsAddCheckedAt,
	CreateDatasetsAndCustomersTable,
	UserBelongsToCustomer,
	AlterEventsAddOnDataset,
	CreateBrandCategoryTable,
	AddCategoryIDToBrands,
	AlterEventsAddOnCustomer,
	AlterEventsAddOnCustomerOnDatasetIndices,
	AddSftpSessionsTable,
	AlterUserTableAddSftpCreds,
	AddDatasetRelationships,
	AddDatasetManageFlag,
	AddScopedConstraintsToCorps,
	AddDataFingerprinting,
	AddTimestampsToFingerprintCorps,
	DataFingerprintingUpdates,
	DatasetsConsistencyUpdates,
	CreateUserCustomerView,
	TagTypeAndTagNaturalKeys,
	FingerprintTagView,
	UserTableCleanups,
	AddAuditTrail,
	ManageableTypoFix,
	UsersDropSuspendedAt,
	CustomersDropUpdatedAt,
	UsersFixSFTPPubKey,
	UsersFixSFTPPubKeyEncoding,
	FingerprintCorps,
	FingerprintCorpsConstraints,
	FingerprintCorpsRemoveType,
	FingerprintCorpsNotNull,
	MigrateAdminsToUsers,
	DropAdminTables,
	TagJobsAndTasks,
	ChangeCorpsConstraint,
	DataFingerprintingHistory,
	DataFingerprintViewsWithUsers,
	TagJobAndTaskRemove,
	LocationsTableCleanups,
	AddGeonamesTable,
	FpLevelCorpMappings,
	LocationsDropUpdatedAt,
	FingerprintTagConfidence,
	AddGeonamesIDToLocations,
	RemoveGeonamesTables,
	AddGeoanmesPostalCodeSearch,
	DeleteCPGLocations,
	AddUpdatedAtToFingerprints,
	InitialTagAppCreation,
	FingerprintUpdateCascade,
	FingerprintSearchPerformanceUpdates,
	IndexSortBrandConsConfidence,
	DropOldFingerprintTables,
	ExtraFieldsOnTag,
	AddCorpMappings,
	AddCorpMappingConstraints,
	FixCorpMappingConstraints,
	CustomerUserLogin,
}

func Drop(db *sql.DB) error {
	rows, err := db.Query(`SELECT 'DROP TABLE IF EXISTS "' || tablename || '" CASCADE;' FROM pg_tables WHERE schemaname='public';`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var query string
		err := rows.Scan(&query)
		if err != nil {
			return err
		}
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func Migrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;"); err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS migrations(version INTEGER NOT NULL);`)
	if err != nil {
		panic(err)
	}

	var version int
	row := tx.QueryRow(`SELECT version FROM migrations LIMIT 1;`)
	if err := row.Scan(&version); err != nil {
		if errors.Cause(err) != sql.ErrNoRows {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec(`INSERT INTO migrations(version) VALUES(0);`); err != nil {
			tx.Rollback()
			return err
		}

		version = 0
	}

	for i := version; i < len(Migrations); i++ {
		if _, err := tx.Exec(Migrations[i]); err != nil {
			tx.Rollback()
			return err
		}

		if _, err = tx.Exec(`UPDATE migrations SET version=$1;`, i+1); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
