package database

import (
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type CorporationType = tables.CorporationType

func (db *DBAL) CorporationTypeUpsert(ct *CorporationType) error {
	return tables.CorporationTypes.Upsert(db, ct)
}

func (db *DBAL) CorporationTypeGet(corporationTypeID string) (CorporationType, error) {
	return tables.CorporationTypes.Get(db, corporationTypeID)
}

type CorporationTypeListArgs struct {
	DatasetID string `json:"datasetID"`
}

func (db *DBAL) CorporationTypeList(args CorporationTypeListArgs) (cts []CorporationType, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.CorporationTypes.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.CorporationTypes.Table())
	builder.Write(`WHERE dataset_id=$1`, args.DatasetID)
	builder.Write(`ORDER BY corporation_type ASC`)

	query, queryArgs := builder.MustBuild()
	return tables.CorporationTypes.List(db, query, queryArgs...)
}

func (db *DBAL) CorporationTypeDelete(corporationTypeID string) error {
	return tables.CorporationTypes.Delete(db, corporationTypeID)
}

func (db *DBAL) CorporationTypeCreate(datasetID string, corporationType string, description string) (CorporationType, error) {
	ct := CorporationType{
		DatasetID:         datasetID,
		CorporationTypeID: crypto.NewUUID(),
		CorporationType:   corporationType,
		Description:       description,
	}
	return ct, tables.CorporationTypes.Insert(db, &ct)
}
