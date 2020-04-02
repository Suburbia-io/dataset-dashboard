package migrations

// language=SQL
const AddScopedConstraintsToCorps = `
DROP INDEX IF EXISTS corporations_slug;
CREATE UNIQUE INDEX corporations_slug
  ON corporations (slug, dataset_entity) WHERE TRIM(slug) <> '';

DROP INDEX IF EXISTS corporations_cusip;
CREATE UNIQUE INDEX corporations_cusip
  ON corporations (cusip, dataset_entity) WHERE TRIM(cusip) <> '';

DROP INDEX IF EXISTS corporations_isin;
CREATE UNIQUE INDEX corporations_isin
  ON corporations (isin, dataset_entity) WHERE TRIM(isin) <> '';

DROP INDEX IF EXISTS corporations_symbol;
CREATE UNIQUE INDEX corporations_symbol
  ON corporations(exchange, code, dataset_entity);
`
