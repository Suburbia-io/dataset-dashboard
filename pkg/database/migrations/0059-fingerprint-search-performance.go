package migrations

// language=SQL
const FingerprintSearchPerformanceUpdates = `
CREATE EXTENSION IF NOT EXISTS btree_gin;

DROP INDEX fingerprints_raw_text_idx;
DROP INDEX fingerprints_annotations_idx;

CREATE INDEX ON fingerprints USING gin(dataset_id, raw_text gin_trgm_ops);
CREATE INDEX ON fingerprints USING gin(dataset_id, annotations gin_trgm_ops);

ALTER TABLE tag_app_tags DROP CONSTRAINT tag_app_tags_pkey;
ALTER TABLE tag_app_tags ADD PRIMARY KEY(dataset_id, tag_app_id, tag_type_id, fingerprint);
`
