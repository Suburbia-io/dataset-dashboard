package migrations

// language=SQL
const ManageableTypoFix = `
ALTER TABLE datasets
  RENAME COLUMN managable TO manageable;
`
