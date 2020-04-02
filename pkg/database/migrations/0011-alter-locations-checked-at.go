package migrations

// language=SQL
const AlterLocationsAddCheckedAt = `
ALTER TABLE locations 
    ADD COLUMN checked_at TIMESTAMPTZ;

CREATE INDEX locations_checked
  ON locations (checked_at);
`
