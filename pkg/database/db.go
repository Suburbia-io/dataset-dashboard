package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/slice"

	"github.com/Suburbia-io/dashboard/pkg/database/migrations"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
	"github.com/lib/pq"
)

// Interface can stand in for a database connection or a transaction.
type DBi = tables.DBi

type Config struct {
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBSSLMode  string
}

type DBAL struct {
	*sql.DB
}

func Bootstrap(config Config, fresh bool) (db *DBAL, err error) {
	db = &DBAL{}
	db.DB, err = sql.Open("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=%s",
		config.DBHost, config.DBName, config.DBUser, config.DBPassword, config.DBPort, config.DBSSLMode,
	))
	if err != nil {
		return db, err
	}

	if fresh {
		if err := db.Fresh(); err != nil {
			return nil, err
		}
	} else {
		if err := db.Migrate(); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *DBAL) Migrate() error {
	return migrations.Migrate(db.DB)
}

func (db *DBAL) Fresh() error {
	if err := migrations.Drop(db.DB); err != nil {
		return err
	}
	return db.Migrate()
}

func (db *DBAL) Close() error {
	return db.DB.Close()
}

// TODO: Remove
func dbIsNullConstraintErr(err error, column string) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23502" && pqErr.Column == column
	}
	return false
}

func (db *DBAL) ExecOne(query string, args ...interface{}) (result sql.Result, n int, err error) {
	return dbExecOne(db, query, args...)
}

type dbExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func dbExecOne(db dbExec, query string, args ...interface{}) (result sql.Result, n int, err error) {
	result, err = db.Exec(query, args...)
	if err != nil {
		return result, 0, err
	}

	if n, err := result.RowsAffected(); err != nil {
		return result, 0, err
	} else if n > 1 {
		panic("update too many rows: " + string(n))
	} else {
		return result, int(n), nil
	}
}

func (db *DBAL) DumpTableToSqlite(sqlite *sql.DB, datasetID string, table interface {
	View() string
	SelectCols() string
	DumpToSqlite(db DBi, sqlite DBi, selectQuery string, args ...interface{}) error
}) error {

	// get the psql table contents
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(table.SelectCols())
	builder.Write(`FROM`)
	builder.Write(table.View())
	builder.Write(`WHERE TRUE`)

	// test if dataset scope filter can be added
	colsList := strings.Split(table.SelectCols(), ",")
	if slice.ContainsString(colsList, "dataset_id") {
		builder.Write(`AND dataset_id = $1`, datasetID)
	}

	// test if archived at filter can be added
	if slice.ContainsString(colsList, "archived_at") {
		builder.Write(`AND archived_at IS NULL`)
	}

	query, args := builder.MustBuild()
	return table.DumpToSqlite(db.DB, sqlite, query, args...)
}

type OrderBy struct {
	Key        string `json:"key"`
	Desc       bool   `json:"desc"`
	NullsFirst bool   `json:"nullsFirst"`
}

func dbBuildOrderBy(sorts []OrderBy, keys map[string]string) (result string, ok bool) {
	if len(sorts) == 0 {
		return "", true
	}

	var parts []string
	for _, s := range sorts {
		part, ok := keys[s.Key]
		if !ok {
			return "", false
		}
		if s.Desc {
			part += " DESC"
		} else {
			part += " ASC"
		}
		if s.NullsFirst {
			part += " NULLS FIRST"
		} else {
			part += " NULLS LAST"
		}
		parts = append(parts, part)
	}

	return "ORDER BY " + strings.Join(parts, ", "), true
}

// TODO: Deprecated?
func dbTX(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

// Postgres doesn't use serializable transactions by default. This wrapper will
// run the enclosed function within a serializable transaction and
// automatically retry when appropriate.
func WithTx(db *sql.DB, fn func(*sql.Tx) error) error {
	wrapped := func() error {
		// Start a transaction.
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
		if err == nil {
			err = fn(tx)
		}

		if err == nil {
			err = tx.Commit()
		}

		if err != nil {
			_ = tx.Rollback()
		}

		return err
	}

	timeout := 100 * time.Millisecond

	for {
		err := wrapped()
		if err == nil {
			return nil
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "40001" {
				log.Printf(
					"Postgres serilization failure. Retrying in %v.", timeout)
				time.Sleep(timeout)
				if timeout < time.Second {
					timeout += 100 * time.Millisecond
				}
				continue
			}
		}

		return err
	}
}
