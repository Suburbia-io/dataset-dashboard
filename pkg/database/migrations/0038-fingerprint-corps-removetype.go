package migrations

// language=SQL
const FingerprintCorpsRemoveType = `
ALTER TABLE corporations
  DROP COLUMN corporation_type_id;
`
