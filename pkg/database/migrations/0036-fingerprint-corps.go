package migrations

// language=SQL
const FingerprintCorps = `
ALTER TABLE corporations
  RENAME COLUMN dataset_entity TO dataset_id;

DROP TABLE fingerprint_corps;
DROP TABLE dataset_corp_types;

CREATE TABLE corporation_types (
  corporation_type_id    UUID PRIMARY KEY,
  dataset_id UUID REFERENCES datasets(dataset_id),
  corporation_type TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  created_at    TIMESTAMPTZ NOT NULL
);

ALTER TABLE corporations 
  ADD COLUMN corporation_type_id UUID REFERENCES corporation_types(corporation_type_id);

ALTER TABLE corp_mappings
  RENAME COLUMN mapping_id TO corporation_mapping_id;

ALTER TABLE corp_mappings
  RENAME COLUMN corp_entity TO corporation_id;

ALTER TABLE corp_mappings
	DROP COLUMN corp_parent;

ALTER TABLE corp_mappings
	DROP COLUMN updated_at;

ALTER TABLE corp_mappings
	DROP COLUMN archived_at;

ALTER TABLE corp_mappings 
  ADD COLUMN tag_type_id UUID REFERENCES tag_types(tag_type_id);

ALTER TABLE corp_mappings 
  ADD COLUMN tag_id UUID REFERENCES tags(tag_id);

ALTER TABLE corp_mappings 
  RENAME TO corporation_mappings;
`
