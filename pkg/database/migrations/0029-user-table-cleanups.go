package migrations

// language=SQL
const UserTableCleanups = `

-- Rename customer_entity to customer_id.
ALTER TABLE users
  RENAME COLUMN customer_entity TO customer_id;
ALTER TABLE users
	RENAME CONSTRAINT users_customer_entity_fkey TO users_customer_id_fkey;
ALTER TABLE user_view
	RENAME COLUMN customer_entity TO customer_id;

-- Drop some indexes.
DROP INDEX users_lower_name;
DROP INDEX users_lower_email;
DROP INDEX users_activity_index;
DROP INDEX users_created_index;
CREATE UNIQUE INDEX users_email_key
 	ON users (email);

-- Changes to sftp username.
ALTER TABLE users
	ALTER COLUMN sftp_username SET NOT NULL;

-- Change sftp pub key to TEXT.
-- We need to drop and recreate user_view to set new data type on sftp_pub_key.
DROP VIEW user_view;
ALTER TABLE users
	ALTER COLUMN sftp_pub_key SET DATA TYPE TEXT;
-- Copied from migration 0027-user-customer-view.go.
CREATE OR REPLACE VIEW user_view AS SELECT
  users.*,
  customers.name AS customer_name
FROM users LEFT JOIN customers ON users.customer_id=customers.customer_id;

-- Change to multi-column primary key
ALTER TABLE users
	ALTER COLUMN customer_id SET NOT NULL;
ALTER TABLE users
  DROP CONSTRAINT users_pkey CASCADE;
ALTER TABLE users
  ADD PRIMARY KEY (customer_id, user_id);
CREATE UNIQUE INDEX users_user_id_key
  ON users (user_id);
ALTER TABLE events
  ADD CONSTRAINT events_by_user_fkey FOREIGN KEY (by_user) REFERENCES users(user_id);
ALTER TABLE events
  ADD CONSTRAINT events_on_user_fkey FOREIGN KEY (on_user) REFERENCES users(user_id);
ALTER TABLE sftp_sessions
  ADD CONSTRAINT sftp_sessions_user_entity_fkey FOREIGN KEY (user_entity) REFERENCES users(user_id);
ALTER TABLE user_session_tokens
  ADD CONSTRAINT user_session_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id);
`
