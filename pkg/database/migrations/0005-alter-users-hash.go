package migrations

// language=SQL
const AlterUserTableAddHash = `
ALTER TABLE users
  ADD COLUMN hash TEXT NOT NULL;
`
