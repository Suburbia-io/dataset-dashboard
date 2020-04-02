package migrations

// language=SQL
const AlterAuthTokensAddSalt = `
ALTER TABLE auth_tokens 
    ADD COLUMN salt TEXT;
`
