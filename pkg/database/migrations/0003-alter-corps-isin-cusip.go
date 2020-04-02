package migrations

// language=SQL
const AlterCorpsTableAddISINCUSIP = `
ALTER TABLE corporations
  ADD COLUMN slug TEXT NOT NULL DEFAULT '',
  ADD COLUMN isin TEXT NOT NULL DEFAULT  '',
  ADD COLUMN cusip TEXT NOT NULL DEFAULT '';

CREATE UNIQUE INDEX corporations_slug
  ON corporations (slug) WHERE TRIM(slug) <> '';

CREATE UNIQUE INDEX corporations_isin
  ON corporations (isin) WHERE TRIM(isin) <> '';

CREATE UNIQUE INDEX corporations_cusip
  ON corporations (cusip) WHERE TRIM(cusip) <> '';
`
