package migrations

// language=SQL
const UsersDropSuspendedAt = `
-- We need to recreate user_view to drop suspended_at.
DROP VIEW user_view;
ALTER TABLE users
  DROP COLUMN suspended_at;
-- This is copied from migration 0029-user-table-cleanup.go.
CREATE OR REPLACE VIEW user_view AS SELECT
  users.*,
  customers.name AS customer_name
FROM users LEFT JOIN customers ON users.customer_id=customers.customer_id;
`
