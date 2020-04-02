package migrations

// language=SQL
const AlterTickersCorporationsTable = `
ALTER TABLE events DROP COLUMN on_corporation;
DROP TABLE IF EXISTS corporations CASCADE;

ALTER TABLE tickers RENAME TO corporations;
ALTER TABLE corporations 
  RENAME COLUMN ticker_id TO corporation_id;
ALTER TABLE corporations ALTER COLUMN exchange DROP NOT NULL;
ALTER TABLE corporations ALTER COLUMN code DROP NOT NULL;
DROP INDEX IF EXISTS tickers_symbol;
CREATE UNIQUE INDEX corporations_symbol 
  ON corporations(exchange, code);
DROP INDEX IF EXISTS tickers_name;
CREATE INDEX corporation_name 
  ON corporations(lower(name));

ALTER TABLE events RENAME COLUMN on_ticker TO on_corporation;
ALTER TABLE events 
  ADD CONSTRAINT events_on_corporations_fkey FOREIGN KEY (on_corporation) REFERENCES corporations(corporation_id);
CREATE INDEX events_on_corporation_index
  ON events (on_corporation,name);
DROP INDEX IF EXISTS events_on_ticker_index;

UPDATE events SET name = replace(name, 'Ticker', 'Corporation');

DROP TABLE IF EXISTS tickers CASCADE;
`
