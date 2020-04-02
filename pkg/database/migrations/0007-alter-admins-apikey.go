package migrations

// language=SQL
const AlterAdminTableAddApiKey = `
ALTER TABLE admins 
    ADD COLUMN api_key TEXT UNIQUE;
`
