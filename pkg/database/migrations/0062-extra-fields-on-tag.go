package migrations

// language=SQL
const ExtraFieldsOnTag = `
ALTER TABLE tags ADD COLUMN internal_notes TEXT NOT NULL DEFAULT '';
ALTER TABLE tags ADD COLUMN is_included BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE tags ADD COLUMN grade INT NOT NULL DEFAULT 0;
ALTER TABLE tags ADD COLUMN num_fingerprints INT NOT NULL DEFAULT 0;
ALTER TABLE tags ADD COLUMN num_line_items INT NOT NULL DEFAULT 0;
`
