package migrations

// language=SQL
const AddGeonamesTable = `

-- More info here: http://download.geonames.org/export/dump/readme.txt
CREATE TABLE geonames (
  geoname_id INT NOT NULL PRIMARY KEY CHECK (geoname_id >= 0),
  name VARCHAR(200) NOT NULL,
  ascii_name VARCHAR(200) NOT NULL,
  alternate_names VARCHAR(10000) NOT NULL,
  latitude  VARCHAR(10) NOT NULL,
  longitude VARCHAR(10) NOT NULL,
  feature_class VARCHAR(1) NOT NULL,
  feature_code VARCHAR(10) NOT NULL,
  country_code VARCHAR(2) NOT NULL,
  admin1_code VARCHAR(20) NOT NULL,
  admin2_code VARCHAR(80) NOT NULL,
  admin3_code VARCHAR(20) NOT NULL,
  admin4_code VARCHAR(20) NOT NULL,
  population BIGINT NOT NULL CHECK (population >= 0),
  modification_date VARCHAR(10) NOT NULL
);

CREATE TABLE geonames_admin1_codes (
  country_code VARCHAR(2) NOT NULL,
  admin1_code VARCHAR(20) NOT NULL,
  name VARCHAR(200) NOT NULL,
  ascii_name VARCHAR(200) NOT NULL,
  geoname_id INT NOT NULL CHECK (geoname_id >= 0)
);

ALTER TABLE geonames_admin1_codes
  ADD PRIMARY KEY (country_code, admin1_code);

CREATE TABLE geonames_admin2_codes (
  country_code VARCHAR(2) NOT NULL,
  admin1_code VARCHAR(20) NOT NULL,
  admin2_code VARCHAR(80) NOT NULL,
  name VARCHAR(200) NOT NULL,
  ascii_name VARCHAR(200) NOT NULL,
  geoname_id INT NOT NULL CHECK (geoname_id >= 0)
);

ALTER TABLE geonames_admin2_codes
  ADD PRIMARY KEY (country_code, admin1_code, admin2_code);

CREATE OR REPLACE VIEW geonames_view AS SELECT
  geonames.*,
  geonames_admin1_codes.name AS admin1_name,
  geonames_admin1_codes.geoname_id AS admin1_geoname_id,
  geonames_admin2_codes.name AS admin2_name,
  geonames_admin2_codes.geoname_id AS admin2_geoname_id
FROM geonames LEFT JOIN geonames_admin1_codes ON geonames.country_code=geonames_admin1_codes.country_code AND
                                                 geonames.admin1_code=geonames_admin1_codes.admin1_code
              LEFT JOIN geonames_admin2_codes ON geonames.country_code=geonames_admin2_codes.country_code AND
                                                 geonames.admin1_code=geonames_admin2_codes.admin1_code AND
                                                 geonames.admin2_code=geonames_admin2_codes.admin2_code;

CREATE TABLE geonames_import (
  geonames_import_id  UUID NOT NULL PRIMARY KEY,
  started_at TIMESTAMPTZ NOT NULL,
  finished_at TIMESTAMPTZ,
  error TEXT NOT NULL,
  version TEXT NOT NULL
);
`
