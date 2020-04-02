package migrations

// language=SQL
const CustomerUserLogin = `
DROP VIEW user_view;
ALTER TABLE users
	ADD COLUMN login_token TEXT UNIQUE,
	ADD COLUMN login_token_expires_at TIMESTAMPTZ;
-- This is copied from migration 0029-user-table-cleanup.go.
CREATE OR REPLACE VIEW user_view AS SELECT
  users.*,
  customers.name AS customer_name
FROM users LEFT JOIN customers ON users.customer_id=customers.customer_id;
`
