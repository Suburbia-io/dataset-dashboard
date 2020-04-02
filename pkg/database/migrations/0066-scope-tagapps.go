package migrations

// language=SQL
const ScopeTagAppss = `
	ALTER TABLE corp_mapping_rules
	DROP CONSTRAINT unique_corp_mapping_rule;

	ALTER TABLE corp_mapping_rules
	ADD CONSTRAINT unique_corp_mapping_rule UNIQUE (corp_mapping_id, corp_id, from_date, country);
`
