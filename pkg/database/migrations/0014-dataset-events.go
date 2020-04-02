package migrations

// language=SQL
const AlterEventsAddOnDataset = `
ALTER TABLE events
ADD COLUMN on_dataset UUID REFERENCES datasets(dataset_id);
`
