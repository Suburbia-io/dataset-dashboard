package migrations

// language=SQL
const LocationsDropUpdatedAt = `
ALTER TABLE locations
  DROP COLUMN updated_at;
`
