package migrations

// language=SQL
const UsersFixSFTPPubKeyEncoding = `
UPDATE users
  SET sftp_pub_key = encode(u2.sftp_pub_key::bytea, 'escape')
FROM users as u2
  WHERE users.user_id = u2.user_id AND u2.sftp_pub_key != '';
`
