package migrations

// language=SQL
const AddDatasetManageFlag = `
ALTER TABLE datasets
  ADD COLUMN managable BOOL DEFAULT FALSE;
`
