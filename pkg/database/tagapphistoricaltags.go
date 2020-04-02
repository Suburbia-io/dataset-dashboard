package database

import (
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type TagAppHistoricalTag = tables.TagAppHistoricalTag

func (db *DBAL) tagAppHistoricalTagInsert(tx DBi, t *TagAppTag) error {
	return tables.TagAppHistoricalTags.Insert(tx, &TagAppHistoricalTag{
		DatasetID:   t.DatasetID,
		Fingerprint: t.Fingerprint,
		TagTypeID:   t.TagTypeID,
		TagAppID:    t.TagAppID,
		TagID:       t.TagID,
		Confidence:  t.Confidence,
		UpdatedAt:   t.UpdatedAt,
		UserID:      t.UserID,
	})
}

type TagAppHistoricalTagListArgs struct {
	DatasetID   string
	Fingerprint string
	TagTypeID   string
	TagAppID    *string
}

func (db *DBAL) TagAppHistoricalTagList(
	args TagAppHistoricalTagListArgs,
) (
	l []TagAppHistoricalTag,
	err error,
) {
	b := qb.Builder{}
	b.Write(`SELECT`)
	b.Write(tables.TagAppHistoricalTags.SelectCols())
	b.Write(`FROM`)
	b.Write(tables.TagAppHistoricalTags.View())
	b.Write(`WHERE dataset_id=$1`, args.DatasetID)
	b.Write(`AND fingerprint=$1`, args.Fingerprint)
	b.Write(`AND tag_type_id=$1`, args.TagTypeID)
	if args.TagAppID != nil {
		b.Write(`AND tag_app_id=$1`, args.TagAppID)
	}
	b.Write(`ORDER BY updated_at DESC`)

	query, queryArgs := b.MustBuild()
	return tables.TagAppHistoricalTags.List(db, query, queryArgs...)
}
