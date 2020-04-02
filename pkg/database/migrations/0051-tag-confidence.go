package migrations

const FingerprintTagConfidence = `

ALTER TABLE fingerprint_tags
  ADD COLUMN user_confidence REAL NOT NULL DEFAULT 0.8;

-- Aggregate confidence.
ALTER TABLE fingerprint_tags
  ADD COLUMN confidence REAL NOT NULL DEFAULT 0.8;

DROP VIEW fingerprint_tag_view;

CREATE VIEW fingerprint_tag_view AS SELECT
  fingerprint_tags.*,
  tags.tag AS tag,
  users.email AS user_email,
  users.name AS user_name
FROM fingerprint_tags
JOIN tags ON
  fingerprint_tags.dataset_id=tags.dataset_id AND
  fingerprint_tags.tag_type_id=tags.tag_type_id AND
  fingerprint_tags.tag_id=tags.tag_id
LEFT JOIN users ON
  fingerprint_tags.user_id=users.user_id;

ALTER TABLE fingerprint_tag_historicals
  ADD COLUMN user_confidence REAL NOT NULL DEFAULT 0.8;

ALTER TABLE fingerprint_tag_historicals
  ADD COLUMN confidence REAL NOT NULL DEFAULT 0.8;

DROP VIEW fingerprint_tag_historical_view;

CREATE VIEW fingerprint_tag_historical_view AS SELECT
  fingerprint_tag_historicals.*,
  tags.tag AS tag,
  users.email AS user_email,
  users.name AS user_name
FROM fingerprint_tag_historicals LEFT JOIN tags ON
  fingerprint_tag_historicals.dataset_id=tags.dataset_id AND
  fingerprint_tag_historicals.tag_type_id=tags.tag_type_id AND
  fingerprint_tag_historicals.tag_id=tags.tag_id
LEFT JOIN users ON
  fingerprint_tag_historicals.user_id=users.user_id;
`
