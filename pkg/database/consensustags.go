package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
)

// Update a single consensus tag from tag_app_tags table.
const updateConsensusListQuery = `SELECT
  tags.tag_id     AS tag_id,
  (tags.confidence * apps.weight) AS confidence
FROM
  tag_app_tags AS tags
JOIN
  tag_apps AS apps
ON
  apps.tag_app_id=tags.tag_app_id
WHERE
  tags.dataset_id=$1 AND
  tags.fingerprint=$2 AND
  tags.tag_type_id=$3 AND
  apps.archived_at IS NULL
ORDER BY
  tags.tag_id`

func (db *DBAL) updateConsensusTag(
	tx DBi,
	datasetID string,
	fingerprint string,
	tagTypeID string,
) error {
	rows, err := tx.Query(
		updateConsensusListQuery,
		datasetID, fingerprint, tagTypeID)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to list tag app tags: %w", err).
			Alert()
	}
	defer rows.Close()

	conf := map[string]float64{}
	weights := map[string]float64{}
	confTotal := float64(0)
	sourceCount := int64(0)

	var tagID string
	var confidence float64

	for rows.Next() {
		if err := rows.Scan(&tagID, &confidence); err != nil {
			return errors.Unexpected.
				Wrap("Failed to scan (tag, confidence): %w", err).
				Alert()
		}
		if _, ok := conf[tagID]; !ok {
			conf[tagID] = 0
			weights[tagID] = 0
		}
		if confidence > 0 {
			sourceCount++
		}
		conf[tagID] += (1 - conf[tagID]) * confidence
		weights[tagID] += confidence
		confTotal += confidence
	}

	// No consensus tag found. Delete.
	if confTotal == 0 {
		return tables.ConsensusTags.Delete(
			tx,
			datasetID,
			fingerprint,
			tagTypeID)
	}

	// Find the consensus tag and confidence.
	consConf := float64(0)
	consTagID := ""

	for tagID, weight := range weights {
		c := conf[tagID] * (weight / confTotal)
		if c > consConf {
			consTagID = tagID
			consConf = c
		}
	}

	// Store the consensus tag.
	return tables.ConsensusTags.Upsert(tx, &tables.ConsensusTag{
		DatasetID:   datasetID,
		Fingerprint: fingerprint,
		TagTypeID:   tagTypeID,
		TagID:       consTagID,
		Confidence:  consConf,
		SourceCount: sourceCount,
		UpdatedAt:   time.Now(),
	})
}

// This is only called from an external tool, so it may be abit rough.
func (db *DBAL) UpdateConsensusForDatasetCmd(datasetID string) {
	tagTypes, err := db.TagTypeList(TagTypeListArgs{DatasetID: datasetID})
	if err != nil {
		log.Fatalf("Failed to list tag types: %v", err)
	}

	rows, err := db.Query(
		`SELECT fingerprint FROM fingerprints WHERE dataset_id=$1`,
		datasetID)
	if err != nil {
		log.Fatalf("Failed to list fingerprints: %v", err)
	}
	defer rows.Close()

	nextFP := func() (string, bool) {
		if !rows.Next() {
			return "", false
		}
		fp := ""
		if err := rows.Scan(&fp); err != nil {
			log.Fatalf("Failed to scan fingerprint: %v", err)
		}
		return fp, true
	}

	count := 0

	more := true
	for more {
		fp := ""
		_ = WithTx(db.DB, func(tx *sql.Tx) error {
			for i := 0; i < 2048; i++ {
				fp, more = nextFP()
				if !more {
					return nil
				}

				count++

				for _, tt := range tagTypes {
					err := db.updateConsensusTag(tx, datasetID, tt.TagTypeID, fp)
					if err != nil {
						log.Fatalf("Failed to update consensus tag: %v", err)
					}
				}
			}
			return nil
		})
		log.Printf("Completed: %d", count)
	}

}
