package database

import (
	"database/sql"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
)

type CustomerDataset = tables.CustomerDataset

func (db *DBAL) CustomerDatasetMappingGet(customerID string, datasetID string) (customerDataset CustomerDataset, err error) {
	return tables.CustomerDatasets.Get(db, customerID, datasetID)
}

func (db *DBAL) CustomerDatasetMappingList(customerID string) (customerDatasets []CustomerDataset, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.CustomerDatasets.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.CustomerDatasets.View())
	builder.Write(`WHERE customer_entity = $1`, customerID)

	query, queryArgs := builder.MustBuild()
	return tables.CustomerDatasets.List(db, query, queryArgs...)
}

func (db *DBAL) CustomerDatasetMappingSet(customerID string, datasetIDs []string) (err error) {
	builder := qb.Builder{}
	builder.Write(`DELETE FROM customer_datasets`)
	builder.Write(`WHERE customer_entity = $1`, customerID)
	query, queryArgs := builder.MustBuild()

	return WithTx(db.DB, func(tx *sql.Tx) (err error) {
		_, err = tx.Exec(query, queryArgs...)

		for _, datasetID := range datasetIDs {
			err = tables.CustomerDatasets.Insert(tx, &CustomerDataset{
				CustomerDatasetID: crypto.NewUUID(),
				DatasetEntity:     datasetID,
				CustomerEntity:    customerID,
				CreatedAt:         time.Now(),
			})
			if err != nil {
				break
			}
		}

		return err
	})
}
