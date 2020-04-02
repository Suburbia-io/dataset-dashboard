package migrations

// language=SQL
const FingerprintCorpsNotNull = `
ALTER TABLE corporations
  ALTER COLUMN exchange SET DEFAULT '';

ALTER TABLE corporations
  ALTER COLUMN exchange SET NOT NULL;

ALTER TABLE corporations
  ALTER COLUMN code SET DEFAULT '';

ALTER TABLE corporations
  ALTER COLUMN code SET NOT NULL;


ALTER TABLE corporations
  ALTER COLUMN isin SET DEFAULT '';

ALTER TABLE corporations
  ALTER COLUMN isin SET NOT NULL;


ALTER TABLE corporations
  ALTER COLUMN cusip SET DEFAULT '';

ALTER TABLE corporations
  ALTER COLUMN cusip SET NOT NULL;`
