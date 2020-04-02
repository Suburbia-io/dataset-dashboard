package migrations

// language=SQL
const DataFingerprintingHistory = `

ALTER TABLE fingerprint_tags
  ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE fingerprint_tags
  ADD COLUMN user_id UUID REFERENCES users(user_id);

DROP VIEW fingerprint_tag_view;

CREATE VIEW fingerprint_tag_view AS SELECT
  fingerprint_tags.*,
  tags.tag AS tag
FROM fingerprint_tags JOIN tags ON
  fingerprint_tags.dataset_id=tags.dataset_id AND
  fingerprint_tags.tag_type_id=tags.tag_type_id AND
  fingerprint_tags.tag_id=tags.tag_id;

CREATE TABLE fingerprint_tag_historicals(
  dataset_id  UUID NOT NULL,
  fingerprint TEXT NOT NULL,
  tag_type_id UUID NOT NULL,
  version     INT  NOT NULL,
  tag_id      UUID,
  user_id     UUID REFERENCES users(user_id),
  created_at  TIMESTAMPTZ NOT NULL,

  PRIMARY KEY(dataset_id, fingerprint, tag_type_id, version),

  FOREIGN KEY(dataset_id, tag_type_id, tag_id)
    REFERENCES tags(dataset_id, tag_type_id, tag_id),

  FOREIGN KEY(dataset_id, fingerprint)
    REFERENCES fingerprints(dataset_id, fingerprint)
);

CREATE VIEW fingerprint_tag_historical_view AS SELECT
  fingerprint_tag_historicals.*,
  tags.tag AS tag
FROM fingerprint_tag_historicals LEFT JOIN tags ON
  fingerprint_tag_historicals.dataset_id=tags.dataset_id AND
  fingerprint_tag_historicals.tag_type_id=tags.tag_type_id AND
  fingerprint_tag_historicals.tag_id=tags.tag_id;
`
