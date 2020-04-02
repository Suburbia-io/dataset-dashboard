package migrations

// language=SQL
const AlterEventsAddAuthMethod = `
ALTER TABLE events 
    ADD COLUMN auth_method TEXT;
`
