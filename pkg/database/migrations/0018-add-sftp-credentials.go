package migrations

// language=SQL
const AlterUserTableAddSftpCreds = `
ALTER TABLE users 
    ADD COLUMN sftp_username TEXT UNIQUE,
	ADD COLUMN sftp_pub_key BYTEA;
`
