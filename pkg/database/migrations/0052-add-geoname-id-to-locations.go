package migrations

// language=SQL
const AddGeonamesIDToLocations = `
ALTER TABLE locations
  ADD COLUMN geoname_id INT,
	ADD COLUMN hierarchy JSON NOT NULL DEFAULT '{}',
  DROP COLUMN detected_city,
  DROP COLUMN detected_country;

UPDATE locations
  SET approved = NULL;

CREATE OR REPLACE VIEW location_view AS SELECT
  locations.*,
  COALESCE(CAST(hierarchy->'geonames'->>-1 AS JSON)->>'toponymName', '') AS name,
  COALESCE(CAST(hierarchy->'geonames'->>-1 AS JSON)->>'population', '') AS population,
  COALESCE(CAST(hierarchy->'geonames'->>-2 AS JSON)->>'toponymName', '') AS parent_name,
  COALESCE(CAST(hierarchy->'geonames'->>-2 AS JSON)->>'population', '') AS parent_population,
  COALESCE(CAST(hierarchy->'geonames'->>-2 AS JSON)->>'geonameId', '') AS parent_geoname_id,
  COALESCE(CAST(hierarchy->'geonames'->>-1 AS JSON)->>'countryCode', '') AS country_code
FROM locations;
`
