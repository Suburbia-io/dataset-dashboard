package migrations

// language=SQL
const DropAdminTables = `
DROP TABLE admin_session_tokens;
DROP TABLE admins;
`
