package migrations

// language=SQL
const AddCorpMappings = `
	DROP TABLE IF EXISTS fingerprint_corporations;

	CREATE TABLE corp_mappings(
	  corp_mapping_id  UUID NOT NULL PRIMARY KEY,
		corp_type_id UUID NOT NULL REFERENCES corporation_types(corporation_type_id),
		tag_type_id UUID NOT NULL REFERENCES tag_types(tag_type_id),
		tag_id UUID NOT NULL REFERENCES tags(tag_id)
	);

	CREATE TABLE corp_mapping_rules(
		corp_mapping_rule_id UUID NOT NULL PRIMARY KEY,
		corp_mapping_id UUID NOT NULL REFERENCES corp_mappings(corp_mapping_id),
		corp_id UUID NOT NULL REFERENCES corporations(corporation_id),
		external_notes TEXT NOT NULL,
		internal_notes TEXT NOT NULL,
		from_date TIMESTAMPTZ NOT NULL,
		country TEXT NOT NULL
	);
`
