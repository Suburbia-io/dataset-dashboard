package database

import (
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type Corporation = tables.Corporation

func (db *DBAL) CorporationUpsert(c *Corporation) error {
	if c.CorporationID == "" {
		c.CreatedAt = time.Now()
	}
	return tables.Corporations.Upsert(db, c)
}

func (db *DBAL) CorporationGet(corporationID string) (Corporation, error) {
	return tables.Corporations.Get(db, corporationID)
}

type CorporationListArgs struct {
	DatasetID string `json:"datasetID"`
}

func (db *DBAL) CorporationList(args CorporationListArgs) (cs []Corporation, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Corporations.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Corporations.Table())
	builder.Write(`WHERE dataset_id=$1`, args.DatasetID)
	builder.Write(`ORDER BY name ASC`)

	query, queryArgs := builder.MustBuild()
	return tables.Corporations.List(db, query, queryArgs...)
}

func (db *DBAL) CorporationDelete(corporationID string) error {
	return tables.Corporations.Delete(db, corporationID)
}

func (db *DBAL) CorporationCreate(datasetID string, name string, slug string, exchange string, code string, isin string, cusip string) (Corporation, error) {
	now := time.Now()
	c := Corporation{
		DatasetID:     datasetID,
		CorporationID: crypto.NewUUID(),
		Name:          name,
		Slug:          slug,
		Exchange:      exchange,
		Code:          code,
		Isin:          isin,
		Cusip:         cusip,
		CreatedAt:     now,
		UpdatedAt:     now,
		ArchivedAt:    nil,
	}
	return c, tables.Corporations.Insert(db, &c)
}
