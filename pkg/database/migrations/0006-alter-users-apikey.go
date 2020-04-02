package migrations

// language=SQL
const AlterUserTableAddApiKey = `
ALTER TABLE users 
    ADD COLUMN api_key TEXT;
`
