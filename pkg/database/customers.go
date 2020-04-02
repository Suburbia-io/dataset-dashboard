package database

import (
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validators"
)

type Customer = tables.Customer

func (db *DBAL) CustomerUpsert(c *Customer) error {
	if c.CustomerID == "" {
		c.CreatedAt = time.Now()
	}
	return tables.Customers.Upsert(db, c)
}

func (db *DBAL) CustomerGet(customerID string) (customer Customer, err error) {
	return tables.Customers.Get(db, customerID)
}

type CustomerListArgs struct {
	Search   string `json:"search" schema:"search"`
	Archived *bool  `json:"archived"  schema:"archived"`
	Limit    int    `json:"limit"  schema:"limit"`
	Offset   int    `json:"offset"  schema:"offset"`
}

func (db *DBAL) CustomerList(args CustomerListArgs) (customers []Customer, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Customers.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Customers.View())
	builder.Write(`WHERE TRUE`)

	if args.Archived != nil {
		if *args.Archived {
			builder.Write(`AND archived_at IS NOT NULL`)
		} else {
			builder.Write(`AND archived_at IS NULL`)
		}
	}

	if args.Search != "" {
		search := "%" + args.Search + "%"
		builder.Write(`AND name ILIKE $1`, search)
	}

	builder.Write(`ORDER BY name ASC`)
	builder.Write(`LIMIT $1 OFFSET $2`, args.Limit, args.Offset)

	query, queryArgs := builder.MustBuild()
	return tables.Customers.List(db, query, queryArgs...)
}

func (db *DBAL) CustomerSetName(customerID, name string) (err error) {
	if err := validators.UUIDDeleteMe(customerID); err != nil {
		return errors.CustomerNotFound
	}

	name, err = validators.CorpName(name)
	if err != nil {
		return err
	}
	if name == "" {
		return errors.CustomerInvalidName
	}

	stmt := `UPDATE customers SET name=$1, updated_at=$2 WHERE customer_id=$3 AND archived_at IS NULL`
	_, n, err := db.ExecOne(stmt, name, time.Now(), customerID)
	if err != nil {
		return errors.UnexpectedError(err, "Failed setting name for customer")
	} else if n == 0 {
		return errors.CustomerNotFound
	}

	return nil
}
