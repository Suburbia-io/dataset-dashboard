package migrations

const DatasetsConsistencyUpdates = `
ALTER TABLE datasets DROP COLUMN updated_at;
ALTER TABLE datasets ALTER COLUMN managable SET NOT NULL;`
