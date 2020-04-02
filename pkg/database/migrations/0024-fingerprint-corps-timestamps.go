package migrations

// language=SQL
const AddTimestampsToFingerprintCorps = `
ALTER TABLE fingerprint_corps
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN archived_at TIMESTAMPTZ
`
