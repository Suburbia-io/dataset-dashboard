package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/lib/pq"
)

// Types to use to check for null constraints.
type _is_null struct{}
type _is_not_null struct{}

var is_null _is_null
var is_not_null _is_not_null

func (db *DBAL) insert(table string, data interface{}) error {
	colNames := []string{}
	queryArgs := []interface{}{} // Passed to insert query.
	pHolders := []string{}

	dataV := reflect.ValueOf(data)
	dataT := reflect.TypeOf(data)

	for i := 0; i < dataV.NumField(); i++ {
		colName := dataT.Field(i).Tag.Get("db")
		if colName == "" {
			continue
		}

		colNames = append(colNames, colName)
		queryArgs = append(queryArgs, dataV.Field(i).Interface())
		pHolders = append(pHolders, fmt.Sprintf("$%d", len(queryArgs)))
	}

	query := `INSERT INTO ` + table + ` (` +
		strings.Join(colNames, ",") +
		`)VALUES(` +
		strings.Join(pHolders, ",") + `)`

	_, err := db.Exec(query, queryArgs...)
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		// Not a database error.
		// TODO: Alert!
		return errors.Unexpected.WithMsg("Insert failed: %v", err)
	}

	switch pqErr.Code {
	case "23505": // unique_violation
		return errors.DBDuplicate

	case "23503":
		return errors.DBFKey

	case "23502":
		return errors.DBNullConstraint

	default:
		// TODO: Alert
		return errors.Unexpected.WithMsg("Insert failed: %v", err)
	}
}

func (db *DBAL) getBy(
	dst interface{},
	table string,
	conds map[string]interface{},
) error {
	colNames := make([]string, 0, len(conds))
	scanTargets := []interface{}{}
	queryArgs := []interface{}{}
	whereConds := []string{}

	dstV := reflect.Indirect(reflect.ValueOf(dst))
	dstT := dstV.Type()

	for i := 0; i < dstV.NumField(); i++ {
		colName := dstT.Field(i).Tag.Get("db")
		if colName == "" {
			continue
		}

		colNames = append(colNames, colName)
		scanTargets = append(scanTargets, dstV.Field(i).Addr().Interface())
		if val, ok := conds[colName]; ok {
			switch val.(type) {
			case _is_null:
				whereConds = append(whereConds, fmt.Sprintf("%s IS NULL", colName))

			case _is_not_null:
				whereConds = append(whereConds, fmt.Sprintf("%s IS NOT NULL", colName))
			default:
				queryArgs = append(queryArgs, val)
				whereConds = append(whereConds, fmt.Sprintf("%s=$%d", colName, len(queryArgs)))
			}
		}
	}

	query := `SELECT ` + strings.Join(colNames, ",") + ` ` +
		`FROM ` + table + ` ` +
		`WHERE ` + strings.Join(whereConds, " AND ")

	err := db.QueryRow(query, queryArgs...).Scan(scanTargets...)
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return errors.DBNotFound
	}

	// TODO: Alert.
	return errors.Unexpected.WithMsg("Failed to get row from %s: %v", table, err)
}

func (db *DBAL) update(
	table string,
	values map[string]interface{},
	conds map[string]interface{},
) error {
	sets := []string{}
	whereConds := []string{}
	queryArgs := []interface{}{}

	for key, val := range values {
		queryArgs = append(queryArgs, val)
		sets = append(sets, fmt.Sprintf("%s=$%d", key, len(queryArgs)))
	}

	for key, val := range conds {
		switch val.(type) {
		case _is_null:
			whereConds = append(whereConds, fmt.Sprintf("%s IS NULL", key))

		case _is_not_null:
			whereConds = append(whereConds, fmt.Sprintf("%s IS NOT NULL", key))
		default:
			queryArgs = append(queryArgs, val)
			whereConds = append(whereConds, fmt.Sprintf("%s=$%d", key, len(queryArgs)))
		}
	}

	query := `UPDATE ` + table + ` ` +
		`SET ` + strings.Join(sets, ",") + ` ` +
		`WHERE ` + strings.Join(whereConds, " AND ")

	result, err := db.Exec(query, queryArgs...)
	if err == nil {
		n, err := result.RowsAffected()
		if err != nil {
			// TODO: Alert
			return errors.Unexpected.WithMsg("Failed to get rows affected: %v", err)
		}
		switch n {
		case 0:
			return errors.DBNotFound

		case 1:
			return nil

		default:
			// TODO: Alert
			panic(fmt.Sprintf("Updated affected %d rows.", n))
			return nil
		}
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		// Not a database error.
		// TODO: Alert!
		return errors.Unexpected.WithMsg("Insert failed: %v", err)
	}

	switch pqErr.Code {
	case "23505": // unique_violation
		return errors.DBDuplicate

	case "23503":
		return errors.DBFKey

	case "23502":
		return errors.DBNullConstraint

	default:
		// TODO: Alert
		return errors.Unexpected.WithMsg("Insert failed: %v", err)
	}
}
