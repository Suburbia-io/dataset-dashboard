package migrations

// language=SQL
const FingerprintCorpsConstraints = `
ALTER TABLE corporation_types 
  ADD CONSTRAINT type_unique_per_dataset UNIQUE(dataset_id,corporation_type);

ALTER TABLE corporation_mappings
  DROP COLUMN brand_entity;

ALTER TABLE corporation_mappings
  RENAME COLUMN dataset_entity TO dataset_id;

ALTER TABLE corporation_mappings 
  ADD CONSTRAINT unique_per_tag UNIQUE(tag_type_id,tag_id,dataset_id);
`
