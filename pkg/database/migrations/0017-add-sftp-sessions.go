package migrations

// language=SQL
const AddSftpSessionsTable = `
CREATE TABLE sftp_sessions (
  sftp_session_id  	UUID PRIMARY KEY,
  session_token		TEXT UNIQUE,
  user_entity		UUID REFERENCES users (user_id),
  created_at    	TIMESTAMPTZ NOT NULL,
  updated_at    	TIMESTAMPTZ NOT NULL,
  archived_at   	TIMESTAMPTZ
);
`
