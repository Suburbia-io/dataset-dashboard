package migrations

// language=SQL
const CreateLocationsTable = `
CREATE TABLE locations (
  loc_id       		UUID PRIMARY KEY,
  loc_hash     		TEXT NOT NULL UNIQUE,
  loc_string 	    TEXT NOT NULL,
  detected_city     TEXT NOT NULL,
  detected_country  TEXT NOT NULL,
  approved		    BOOLEAN NULL,	
  created_at        TIMESTAMPTZ NOT NULL,
  updated_at        TIMESTAMPTZ NOT NULL,
  archived_at       TIMESTAMPTZ
);

ALTER TABLE events
  ADD COLUMN on_location UUID REFERENCES locations(loc_id);

CREATE INDEX events_on_location_index
  ON events (on_location,name);
`
