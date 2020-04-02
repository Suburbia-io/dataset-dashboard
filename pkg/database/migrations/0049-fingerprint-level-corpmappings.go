package migrations

// language=SQL
const FpLevelCorpMappings = `
DROP TABLE corporation_mappings CASCADE;

CREATE TABLE fingerprint_corporations(
  fingerprint_corporation_id  UUID NOT NULL PRIMARY KEY,
  dataset_id   UUID NOT NULL REFERENCES datasets(dataset_id),
  corporation_type_id   UUID NOT NULL REFERENCES corporation_types(corporation_type_id),
  corporation_id UUID NOT NULL REFERENCES corporations(corporation_id),
  fingerprint TEXT NOT NULL,
  from_date TIMESTAMPTZ NOT NULL,
  location TEXT NULL,
	created_at TIMESTAMPTZ NOT NULL                                  
);
`
