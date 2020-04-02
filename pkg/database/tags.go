package database

import (
	"database/sql"
	"fmt"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type Tag = tables.Tag

func (db *DBAL) TagUpsert(tag *Tag) error {
	return tables.Tags.Upsert(db, tag)
}

func (db *DBAL) TagLookup(datasetID, tagTypeID, tag string) (t Tag, err error) {
	return tables.Tags.GetByTag(db, datasetID, tagTypeID, tag)
}

func (db *DBAL) TagUpdate(t Tag) error {
	return tables.Tags.Update(db, &t)
}

const humanTagAppTagsExistsQuery = `SELECT EXISTS (SELECT * FROM tag_app_tags
  WHERE dataset_id=$1 AND
    tag_app_id=$2 AND
    tag_type_id=$3 AND
    tag_id=$4)`

const tagAppTagsDeleteQuery = `DELETE FROM tag_app_tags
WHERE
 dataset_id=$1 AND
 tag_type_id=$2 AND
 tag_id=$3`

const consensusTagsDeleteQuery = `DELETE FROM consensus_tags
WHERE
 dataset_id=$1 AND
 tag_type_id=$2 AND
 tag_id=$3`

func (db *DBAL) TagDelete(datasetID, tagTypeID, tagID string) error {
	return WithTx(db.DB, func(tx *sql.Tx) error {
		var humanAppTagsExist bool
		row := db.QueryRow(humanTagAppTagsExistsQuery,
			datasetID,
			TAG_APP_HUMAN_POOL,
			tagTypeID,
			tagID)

		err := row.Scan(&humanAppTagsExist)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to run count query on tag_app_tags table: %w", err).
				Alert()
		}

		if humanAppTagsExist {
			return fmt.Errorf(
				"This tag is currently used by the Humun Tag Pool app and cannot be" +
					" deleted.")
		}

		_, err = db.Exec(tagAppTagsDeleteQuery,
			datasetID,
			tagTypeID,
			tagID)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to delete from tag_app_tags table: %w", err).
				Alert()
		}

		_, err = db.Exec(consensusTagsDeleteQuery,
			datasetID,
			tagTypeID,
			tagID)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to delete from consensus_tags table: %w", err).
				Alert()
		}

		return tables.Tags.Delete(db, datasetID, tagTypeID, tagID)
	})
}

func (db *DBAL) TagGet(datasetID, tagTypeID, tagID string) (t Tag, err error) {
	return tables.Tags.Get(db, datasetID, tagTypeID, tagID)
}

type TagListArgs struct {
	Search     string `json:"search" schema:"search"`
	IsIncluded *bool  `json:"isIncluded" schema:"is-included"`
}

func (db *DBAL) TagList(datasetID, tagTypeID string, args TagListArgs) (tags []Tag, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Tags.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Tags.View())
	builder.Write(`WHERE dataset_id=$1 AND tag_type_id=$2`, datasetID, tagTypeID)

	if args.IsIncluded != nil {
		if *args.IsIncluded {
			builder.Write(`AND is_included IS TRUE`)
		} else {
			builder.Write(`AND is_included IS FALSE`)
		}
	}

	if args.Search != "" {
		builder.Write(`AND (tag ILIKE $1 OR
												description ILIKE $1 OR
                        internal_notes ILIKE $1)`, "%"+args.Search+"%")
	}

	builder.Write(`ORDER BY tag ASC`)

	query, queryArgs := builder.MustBuild()
	return tables.Tags.List(db, query, queryArgs...)
}

const tagCountQuery = `
SELECT
  t.dataset_id,
  t.tag_type_id,
  t.tag_id,
  count(ct.fingerprint) AS num_fingerprints,
  sum(f.count) AS num_line_items
FROM tags AS t
  JOIN consensus_tags AS ct ON
    ct.dataset_id=t.dataset_id AND
    ct.tag_type_id=t.tag_type_id AND
    ct.tag_id=t.tag_id
  JOIN fingerprints AS f ON
    ct.dataset_id = f.dataset_id AND
    ct.fingerprint = f.fingerprint
WHERE
  t.dataset_id = $1 AND
  t.tag_type_id = $2
GROUP BY t.dataset_id, t.tag_type_id, t.tag_id, t.tag;
`

func (db *DBAL) UpdateTagCounts(datasetID, tagTypeID string) error {
	rows, err := db.Query(tagCountQuery, datasetID, tagTypeID)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to list tags: %w", err).
			Alert()
	}
	defer rows.Close()

	return WithTx(db.DB, func(tx *sql.Tx) error {
		for rows.Next() {
			row := Tag{}
			err = rows.Scan(
				&row.DatasetID,
				&row.TagTypeID,
				&row.TagID,
				&row.NumFingerprints,
				&row.NumLineItems)
			if err != nil {
				return err
			}

			err = tables.Tags.UpdateNumFingerprints(
				tx, datasetID, tagTypeID, row.TagID, row.NumFingerprints)
			if err != nil {
				return err
			}

			err = tables.Tags.UpdateNumLineItems(
				tx, datasetID, tagTypeID, row.TagID, row.NumLineItems)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

const tagExportSQLiteCreateQuery = `
CREATE TABLE tags(
  tag_type_id TEXT NOT NULL,
  tag_id      TEXT NOT NULL,
  tag         TEXT NOT NULL,
  description TEXT NOT NULL,
  grade       INT NOT NULL,
  is_included BOOLEAN NOT NULL,

  PRIMARY KEY(tag_type_id, tag_id)
);
`

const tagExportSQLiteInsertQuery = `
INSERT INTO tags(tag_type_id, tag_id, tag, description, grade, is_included)
  VALUES (?,?,?,?,?,?)
`

func (db *DBAL) tagExportSQLite(datasetID string, sqlite *sql.DB) error {
	if _, err := sqlite.Exec(tagExportSQLiteCreateQuery); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create tag type table: %w", err).
			Alert()
	}

	query := `SELECT ` + tables.Tags.SelectCols() + ` FROM ` + tables.Tags.View() +
		` WHERE dataset_id=$1`

	l, err := tables.Tags.List(db, query, datasetID)
	if err != nil {
		return err
	}

	return dbTX(sqlite, func(tx *sql.Tx) error {
		insertStmt, err := tx.Prepare(tagExportSQLiteInsertQuery)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to prepare tag type insert statement: %w", err).
				Alert()
		}

		for _, tag := range l {
			_, err := insertStmt.Exec(
				tag.TagTypeID, tag.TagID, tag.Tag, tag.Description, tag.Grade, tag.IsIncluded)
			if err != nil {
				return errors.Unexpected.
					Wrap("Failed to insert tag type: %w", err).
					Alert()
			}
		}

		return nil
	})
}
