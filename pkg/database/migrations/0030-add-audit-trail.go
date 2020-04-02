package migrations

// language=SQL
const AddAuditTrail = `

CREATE TABLE audit_trails(
	audit_trail_id UUID NOT NULL,
  by_admin UUID,
  by_user  UUID,
  type TEXT NOT NULL,
  related_table TEXT NOT NULL,
  related_id UUID NOT NULL,
  payload JSON,
	created_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE audit_trails
  ADD PRIMARY KEY (audit_trail_id);
ALTER TABLE audit_trails
  ADD CONSTRAINT events_by_admin_fkey FOREIGN KEY (by_admin) REFERENCES admins(admin_id);
ALTER TABLE audit_trails
  ADD CONSTRAINT events_by_user_fkey FOREIGN KEY (by_user) REFERENCES users(user_id);


`
