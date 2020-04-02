package migrations

// language=SQL
const AddUpdatedAtToFingerprints = `
ALTER TABLE fingerprints
	ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
`
