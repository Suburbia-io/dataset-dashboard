package migrations

// language=SQL
const RemoveGeonamesTables = `
DROP VIEW geonames_view;
DROP TABLE geonames, geonames_admin1_codes, geonames_admin2_codes, geonames_import;
`
