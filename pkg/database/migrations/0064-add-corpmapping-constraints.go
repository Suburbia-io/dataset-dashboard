package migrations

// language=SQL
const AddCorpMappingConstraints = `
	ALTER TABLE corp_mapping_rules
	ADD CONSTRAINT unique_corp_mapping_rule UNIQUE (corp_id, from_date, country);

	ALTER TABLE corp_mappings	
	ADD CONSTRAINT unique_corp_mapping UNIQUE (corp_type_id, tag_type_id, tag_id);
`
