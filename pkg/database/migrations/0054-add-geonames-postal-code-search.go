package migrations

// language=SQL
const AddGeoanmesPostalCodeSearch = `
DROP VIEW location_view;

ALTER TABLE locations
	ADD COLUMN geonames_postal_codes JSON NOT NULL DEFAULT '{}',
	ADD COLUMN parsed_country_code TEXT NOT NULL DEFAULT '',
	ADD COLUMN parsed_postal_code TEXT NOT NULL DEFAULT '';

ALTER TABLE locations
  RENAME COLUMN hierarchy TO geonames_hierarchy;

CREATE VIEW location_view AS SELECT
  locations.*,
  COALESCE(CAST(geonames_hierarchy->'geonames'->>-1 AS JSON)->>'toponymName', '') AS name,
  COALESCE(CAST(geonames_hierarchy->'geonames'->>-1 AS JSON)->>'population', '') AS population,
  COALESCE(CAST(geonames_hierarchy->'geonames'->>-2 AS JSON)->>'toponymName', '') AS parent_name,
  COALESCE(CAST(geonames_hierarchy->'geonames'->>-2 AS JSON)->>'population', '') AS parent_population,
  COALESCE(CAST(geonames_hierarchy->'geonames'->>-2 AS JSON)->>'geonameId', '') AS parent_geoname_id,
  COALESCE(CAST(geonames_hierarchy->'geonames'->>-1 AS JSON)->>'countryCode', '') AS country_code
FROM locations;
`
