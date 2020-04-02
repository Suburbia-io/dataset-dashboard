package testdb

import (
	"database/sql"
	"fmt"
	"strconv"
)

type NamedConn struct {
	conn *sql.DB
	name string
}

type Manager struct {
	masterConn    *sql.DB
	connString    string
	namedConnChan chan *NamedConn
}

func NewManager(Prefix string, MaxConn int, DBHost, DBUser, DBPassword, DBPort, DBSSLMode string) (m *Manager, err error) {
	m = &Manager{}

	maxConnections := MaxConn
	m.namedConnChan = make(chan *NamedConn, maxConnections)
	for i := 0; i < maxConnections; i++ {
		m.namedConnChan <- &NamedConn{name: Prefix + "_" + strconv.Itoa(i)}

		m.connString = fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=%s",
			DBHost, DBUser, DBPassword, DBPort, DBSSLMode,
		)

		m.masterConn, err = sql.Open("postgres", m.connString)

		if err := m.drop(Prefix); err != nil {
			return m, err
		}
	}

	return m, err
}

func (tdb *Manager) getNextConn() (nc *NamedConn, err error) {
	nc = <-tdb.namedConnChan
	if nc.conn != nil {
		return nc, nil
	}

	if _, err := tdb.masterConn.Exec(fmt.Sprintf(`create database %s`, nc.name)); err != nil {
		return nc, err
	}

	nc.conn, err = sql.Open("postgres", tdb.connString+" dbname="+nc.name)
	if err != nil {
		return nc, err
	}

	return nc, nil
}

func (tdb *Manager) TearDown() {
	close(tdb.namedConnChan)

	for nc := range tdb.namedConnChan {
		if nc.conn != nil {
			nc.conn.Close()
			tdb.masterConn.Exec(fmt.Sprintf(`drop database %s`, nc.name))
		}
	}

	tdb.masterConn.Close()
}

func (tdb *Manager) NewConn() (db *sql.DB, close func(), err error) {
	nc, err := tdb.getNextConn()
	if err != nil {
		return db, close, err
	}

	return nc.conn, func() {
		tdb.namedConnChan <- nc
	}, nil
}

func (tdb *Manager) drop(Prefix string) error {
	rows, err := tdb.masterConn.Query(`SELECT 'DROP DATABASE IF EXISTS "' || datname || '";' FROM pg_database WHERE datname like '` + Prefix + `_%';`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var query string
		err := rows.Scan(&query)
		if err != nil {
			return err
		}
		fmt.Println("CLEANUP: dropping db: " + query)
		if _, err := tdb.masterConn.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
