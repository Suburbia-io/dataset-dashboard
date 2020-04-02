package migrations

// language=SQL
const AlterUserTableMakeApiKeyUnique = `
ALTER TABLE users 
    ADD UNIQUE (api_key);
`
