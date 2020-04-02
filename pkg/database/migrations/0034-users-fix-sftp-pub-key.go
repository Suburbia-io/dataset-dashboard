package migrations

// language=SQL
const UsersFixSFTPPubKey = `
ALTER TABLE users
	ALTER COLUMN sftp_pub_key SET NOT NULL;
`
