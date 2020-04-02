package migrations

// language=SQL
const InitialTagAppCreation = `

-- Speed up fingerprint table interface.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX ON fingerprints(dataset_id, count);
CREATE INDEX ON fingerprints USING gin(raw_text gin_trgm_ops);
CREATE INDEX ON fingerprints USING gin(annotations gin_trgm_ops);

CREATE TABLE tag_apps(
  tag_app_id  UUID NOT NULL PRIMARY KEY,
  name        TEXT NOT NULL,
  weight      REAL NOT NULL,
  archived_at TIMESTAMPTZ
);

CREATE TABLE tag_app_tags(
  dataset_id  UUID NOT NULL,
  fingerprint TEXT NOT NULL,
  tag_app_id  UUID NOT NULL REFERENCES tag_apps(tag_app_id),
  tag_type_id UUID NOT NULL,
  tag_id      UUID NOT NULL,
  confidence  REAL NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL,

  user_id     UUID NOT NULL REFERENCES users(user_id),

  PRIMARY KEY (dataset_id, fingerprint, tag_type_id, tag_app_id),

  FOREIGN KEY (dataset_id, tag_type_id, tag_id)
    REFERENCES tags(dataset_id, tag_type_id, tag_id),

  FOREIGN KEY (dataset_id, fingerprint)
    REFERENCES fingerprints(dataset_id, fingerprint)
);

CREATE TABLE tag_app_historical_tags(
  dataset_id  UUID NOT NULL,
  fingerprint TEXT NOT NULL,
  tag_app_id  UUID NOT NULL,
  tag_type_id UUID NOT NULL,
  tag_id      UUID NOT NULL,
  confidence  REAL NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL,
  user_id     UUID NOT NULL,

  PRIMARY KEY(dataset_id, fingerprint, tag_type_id, tag_app_id, updated_at)
);

CREATE VIEW tag_app_historical_tag_view AS
SELECT
  h.*,
  tag_apps.name AS tag_app_name,
  users.name AS user_name,
  users.email AS user_email,
  tags.tag AS tag
FROM tag_app_historical_tags AS h
JOIN tag_apps ON tag_apps.tag_app_id=h.tag_app_id
JOIN users ON users.user_id=h.user_id
JOIN tags ON
  tags.dataset_id=h.dataset_id AND
  tags.tag_type_id=h.tag_type_id AND
  tags.tag_id=h.tag_id;

CREATE TABLE consensus_tags(
  dataset_id  UUID NOT NULL,
  fingerprint TEXT NOT NULL,
  tag_type_id UUID NOT NULL,
  tag_id      UUID NOT NULL,
  confidence  REAL NOT NULL,
	source_count INTEGER NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL,

  PRIMARY KEY (dataset_id, fingerprint, tag_type_id),

  FOREIGN KEY (dataset_id, tag_type_id, tag_id)
    REFERENCES tags(dataset_id, tag_type_id, tag_id),

  FOREIGN KEY (dataset_id, fingerprint)
    REFERENCES fingerprints(dataset_id, fingerprint)
);

-- Create a tag app for the human tagging pool.

INSERT INTO tag_apps(
  tag_app_id,
  name,
  weight
) VALUES (
  '00000000-0000-0000-0000-100000000000',
  'Human Tag Pool',
  1
);

-- Populate new tables.

INSERT INTO tag_app_tags(
  dataset_id,
  fingerprint,
  tag_app_id,
  tag_type_id,
  tag_id,
  confidence,
  updated_at,
  user_id
) SELECT
  dataset_id,
  fingerprint,
  '00000000-0000-0000-0000-100000000000',
  tag_type_id,
  tag_id,
  confidence,
  NOW(),
  user_id
FROM fingerprint_tags;

INSERT INTO consensus_tags(
  dataset_id,
  fingerprint,
  tag_type_id,
  tag_id,
  confidence,
	source_count,
  updated_at
) SELECT
  dataset_id,
  fingerprint,
  tag_type_id,
  tag_id,
  confidence,
	1,
  NOW()
FROM fingerprint_tags;
`
