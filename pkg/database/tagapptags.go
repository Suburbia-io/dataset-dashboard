package database

import (
	"database/sql"
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
)

type TagAppTag = tables.TagAppTag

func (db *DBAL) TagAppTagUpsertCSV(
	datasetID string,
	tagAppID string,
	userID string,
	reader *csv.Reader,
) (
	err error,
) {
	const (
		FP         = 0 // The fingerprint.
		TAG_TYPE   = 1 // Tag type slug. Must exist.
		TAG        = 2 // Tag slug. Must exist.
		CONFIDENCE = 3 // Confidence (0,1]
	)

	header, err := reader.Read()
	if err != nil {
		return errors.Unexpected.Wrap("Failed to read header.")
	}
	if header[FP] != "fingerprint" || header[TAG_TYPE] != "tag_type" || header[TAG] != "tag" || header[CONFIDENCE] != "confidence" {
		return errors.Unexpected.Wrap("Invalid header.")
	}

	// Cache slug -> ID mappings.
	tagTypeIDs := map[string]string{}
	tagIDs := map[string]string{}

	return WithTx(db.DB, func(tx *sql.Tx) error {
		// Return io.EOF when done.
		for {
			rec, err := reader.Read()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return errors.Unexpected.Wrap("Failed to read CSV: %w", err)
			}

			// Get the tag type ID.
			tagTypeID, ok := tagTypeIDs[rec[TAG_TYPE]]
			if !ok {
				tt, err := tables.TagTypes.GetByTagType(tx, datasetID, rec[TAG_TYPE])
				if err != nil {
					return err
				}
				tagTypeID = tt.TagTypeID
				tagTypeIDs[rec[TAG_TYPE]] = tagTypeID
			}

			// Could be a delete.
			if rec[TAG] == "" {
				if err := db.tagAppTagDelete(tx, datasetID, rec[FP], tagTypeID, tagAppID); err != nil {
					return err
				}
				continue
			}

			// Get the tag ID.
			tagIDKey := rec[TAG_TYPE] + rec[TAG]
			tagID, ok := tagIDs[tagIDKey]
			if !ok {
				t, err := tables.Tags.GetByTag(tx, datasetID, tagTypeID, rec[TAG])
				if err != nil {
					return err
				}
				tagID = t.TagID
				tagIDs[tagIDKey] = tagID
			}

			confidence, err := strconv.ParseFloat(rec[CONFIDENCE], 64)
			if err != nil {
				return errors.Unexpected.Wrap("Failed to parse confidence: %w", err)
			}

			// Insert the fingerprint tag.
			taTag := TagAppTag{
				DatasetID:   datasetID,
				Fingerprint: rec[FP],
				TagTypeID:   tagTypeID,
				TagAppID:    tagAppID,
				TagID:       tagID,
				Confidence:  confidence,
				UserID:      userID,
			}

			if err := db.tagAppTagUpsert(tx, &taTag); err != nil {
				return err
			}
		}
	})
}

func (db *DBAL) TagAppTagUpsert(t *TagAppTag) error {
	return WithTx(db.DB, func(tx *sql.Tx) error {
		return db.tagAppTagUpsert(tx, t)
	})
}

func (db *DBAL) tagAppTagUpsert(tx *sql.Tx, t *TagAppTag) error {
	t.UpdatedAt = time.Now().UTC()
	if err := tables.TagAppTags.Upsert(tx, t); err != nil {
		return err
	}
	if err := db.tagAppHistoricalTagInsert(tx, t); err != nil {
		return err
	}
	return db.updateConsensusTag(tx, t.DatasetID, t.Fingerprint, t.TagTypeID)
}

func (db *DBAL) TagAppTagDelete(
	datasetID,
	fingerprint,
	tagTypeID,
	tagAppID string,
) error {
	return WithTx(db.DB, func(tx *sql.Tx) error {
		return db.tagAppTagDelete(tx, datasetID, fingerprint, tagTypeID, tagAppID)
	})
}

const tagAppDeleteQuery = `DELETE FROM tag_app_tags
WHERE
 dataset_id=$1 AND
 tag_app_id=$2 AND
 tag_type_id=$3 AND
 fingerprint=$4`

func (db *DBAL) tagAppTagDelete(
	tx DBi,
	datasetID,
	fingerprint,
	tagTypeID,
	tagAppID string,
) error {
	result, err := tx.Exec(tagAppDeleteQuery,
		datasetID,
		tagAppID,
		tagTypeID,
		fingerprint)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to delete from tag_app_tags table: %w", err).
			Alert()
	}

	count, err := result.RowsAffected()
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to get tag_app_tags delete count: %w", err).
			Alert()
	}

	if count > 0 {
		return db.updateConsensusTag(tx, datasetID, fingerprint, tagTypeID)
	}
	return nil
}

type TagAppTagExport struct {
	Fingerprint string
	TagAppID    string
	TagType     string
	Tag         string
	Confidence  float64
}

const tagAppTagsExportCreateTableQuery = `CREATE TABLE tag_app_tags(
 fingerprint TEXT NOT NULL,
 tag_app_id TEXT NOT NULL,
 tag_type TEXT NOT NULL,
 tag TEXT NOT NULL,
 confidence TEXT NOT NULL);`

const tagAppTagsExportInsertQuery = `INSERT INTO tag_app_tags(
 fingerprint,
 tag_app_id,
 tag_type,
 tag,
 confidence
)VALUES(?,?,?,?,?)`

func (db *DBAL) tagAppTagsExportSQLite(
	datasetID string,
	sqlite *sql.DB,
) error {
	// Create table.
	if _, err := sqlite.Exec(tagAppTagsExportCreateTableQuery); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create tag_app_tags table: %w", err).
			Alert()
	}

	builder := qb.Builder{}
	builder.Write(`SELECT fingerprint, tag_app_id, tag_type, tag, confidence FROM tag_app_tags tat`)
	builder.Write(`JOIN tag_types tt ON tat.dataset_id=tt.dataset_id AND tat.tag_type_id=tt.tag_type_id`)
	builder.Write(`JOIN tags t on tat.dataset_id=t.dataset_id and tat.tag_type_id=t.tag_type_id AND tat.tag_id=t.tag_id`)
	builder.Write(`WHERE tat.dataset_id = $1`, datasetID)
	query, queryArgs := builder.MustBuild()

	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to list locations: %w", err).
			Alert()
	}
	defer rows.Close()

	return dbTX(sqlite, func(tx *sql.Tx) error {
		insertStmt, err := tx.Prepare(tagAppTagsExportInsertQuery)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to prepare tag_app_tags insert statement: %w", err).
				Alert()
		}

		for rows.Next() {
			row := TagAppTagExport{}
			err = rows.Scan(&row.Fingerprint, &row.TagAppID, &row.TagType, &row.Tag, &row.Confidence)
			if err != nil {
				return err
			}
			_, err = insertStmt.Exec(row.Fingerprint, row.TagAppID, row.TagType, row.Tag, row.Confidence)
			if err != nil {
				return errors.Unexpected.
					Wrap("Failed to insert tag_app_tags: %w", err).
					Alert()
			}
		}

		return nil
	})
}
