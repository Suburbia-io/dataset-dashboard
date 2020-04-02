package database

import (
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/Suburbia-io/dashboard/pkg/helpers/testdb"
)

var tdb *testdb.Manager

func TestMain(m *testing.M) {
	var config struct{ DB Config }
	_, err := toml.DecodeFile(os.Getenv("SUBURBIA_DASHBOARD_CONFIG"), &config)
	if err != nil {
		log.Fatal(err, "database test config not found", os.Getenv("SUBURBIA_DASHBOARD_CONFIG"))
	}

	tdb, err = testdb.NewManager(
		"db_testing", 2*runtime.NumCPU(),
		config.DB.DBHost,
		config.DB.DBUser,
		config.DB.DBPassword,
		config.DB.DBPort,
		config.DB.DBSSLMode,
	)
	if err != nil {
		panic(err)
	}

	var status int

	defer func() {
		recover()
		tdb.TearDown()
		os.Exit(status)
	}()

	status = m.Run()
}

func NewTestDBAL() (dbal *DBAL, close func()) {
	conn, close, err := tdb.NewConn()
	if err != nil {
		panic(err)
	}

	dbal = &DBAL{DB: conn}
	err = dbal.Fresh()
	if err != nil {
		panic(err)
	}
	return dbal, close
}

func testQuery(db *DBAL, query string, args ...interface{}) bool {
	ok := false
	err := db.QueryRow(query, args...).Scan(&ok)
	if err != nil {
		panic(err)
	}
	return ok
}
