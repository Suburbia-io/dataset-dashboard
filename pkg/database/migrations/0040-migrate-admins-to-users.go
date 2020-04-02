package migrations

// language=SQL
const MigrateAdminsToUsers = `
-- We need to recreate user_view to add these columns.
DROP VIEW user_view;
ALTER TABLE users
	ADD COLUMN is_role_customer_user BOOLEAN NOT NULL DEFAULT true,
  ADD COLUMN is_role_admin BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN is_role_super_admin BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN is_role_labeler BOOLEAN NOT NULL DEFAULT false;
-- This is copied from migration 0029-user-table-cleanup.go.
CREATE OR REPLACE VIEW user_view AS SELECT
  users.*,
  customers.name AS customer_name
FROM users LEFT JOIN customers ON users.customer_id=customers.customer_id;

ALTER TABLE admin_session_tokens DROP CONSTRAINT admin_session_tokens_admin_id_fkey;
ALTER TABLE events DROP CONSTRAINT events_by_admin_fkey;
ALTER TABLE events DROP CONSTRAINT events_on_admin_fkey;
ALTER TABLE audit_trails DROP CONSTRAINT events_by_admin_fkey;

UPDATE audit_trails as at
  SET by_admin = u.user_id
FROM admins as a
  INNER JOIN users AS u ON u.email = a.email
WHERE a.admin_id = at.by_admin;

UPDATE events as e
  SET by_admin = u.user_id
FROM admins as a
  INNER JOIN users AS u ON u.email = a.email
WHERE a.admin_id = e.by_admin;

UPDATE events as e
  SET on_admin = u.user_id
FROM admins as a
  INNER JOIN users AS u ON u.email = a.email
WHERE a.admin_id = e.on_admin;

INSERT INTO users (user_id,
                   name,
                   email,
                   hash,
                   is_role_admin,
                   is_role_super_admin,
                   created_at,
                   archived_at,
                   last_active_at,
                   api_key,
                   customer_id,
                   sftp_username,
                   sftp_pub_key)
SELECT admin_id,
       name,
       email,
       hash,
       true as is_role_admin,
       super as is_role_super_admin,
       created_at,
       archived_at,
       last_active_at,
       api_key,
       '005d9f1f-c5db-1864-a4a0-aa6fe77b76b2' as customer_id,
       REPLACE(email, '@', '.') as sftp_username,
       '' as sftp_pub_key
FROM admins ON CONFLICT (email)
	DO UPDATE SET
		name=EXCLUDED.name,
		email=EXCLUDED.email,
		hash=EXCLUDED.hash,
		is_role_admin=EXCLUDED.is_role_admin,
		is_role_super_admin=EXCLUDED.is_role_super_admin,
		created_at=EXCLUDED.created_at,
		archived_at=EXCLUDED.archived_at,
		last_active_at=EXCLUDED.last_active_at,
		api_key=EXCLUDED.api_key,
		customer_id=EXCLUDED.customer_id,
		sftp_username=EXCLUDED.sftp_username;

ALTER TABLE audit_trails
  ADD CONSTRAINT audit_trails_by_admin_fkey FOREIGN KEY (by_admin) REFERENCES users(user_id);
ALTER TABLE events
  ADD CONSTRAINT events_by_admin_fkey FOREIGN KEY (by_admin) REFERENCES users(user_id);
ALTER TABLE events
  ADD CONSTRAINT events_on_admin_fkey FOREIGN KEY (on_admin) REFERENCES users(user_id);

UPDATE audit_trails
  SET by_user = at2.by_admin
FROM audit_trails AS at2
  WHERE audit_trails.audit_trail_id = at2.audit_trail_id
    AND audit_trails.by_admin IS NOT NULL;
ALTER TABLE audit_trails
  DROP COLUMN by_admin;

ALTER TABLE auth_tokens
	RENAME COLUMN entity_id TO user_id;
ALTER TABLE auth_tokens
	RENAME COLUMN salt TO browser_token;
ALTER TABLE auth_tokens
  ADD CONSTRAINT auth_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id);

UPDATE users
  SET sftp_username = REPLACE(users.email, '@', '.')
FROM users as u2
  WHERE users.user_id = u2.user_id AND u2.sftp_username like '%@%';

ALTER TABLE user_session_tokens
  RENAME TO sessions;

-- Do this in another migration.
-- DROP TABLE admin_session_tokens;
-- DROP TABLE admins;

`
