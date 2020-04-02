package migrations

// language=SQL
const DropOldFingerprintTables = `
DROP VIEW fingerprint_tag_historical_view;
DROP TABLE fingerprint_tag_historicals;
DROP VIEW fingerprint_tag_view;
DROP TABLE fingerprint_tags;
`
