package migrations

const FingerprintTagView = `
DROP TABLE fingerprint_tags;

CREATE TABLE fingerprint_tags(
  dataset_id         UUID NOT NULL,
  fingerprint        TEXT NOT NULL,
  tag_type_id        UUID NOT NULL,
  tag_id             UUID NOT NULL,

  PRIMARY KEY(dataset_id, fingerprint, tag_type_id),

  FOREIGN KEY(dataset_id, tag_type_id, tag_id)
    REFERENCES tags(dataset_id, tag_type_id, tag_id),

  FOREIGN KEY(dataset_id, fingerprint)
    REFERENCES fingerprints(dataset_id, fingerprint)
);

CREATE OR REPLACE VIEW fingerprint_tag_view AS SELECT
  fingerprint_tags.*,
  tags.tag AS tag
FROM fingerprint_tags JOIN tags ON
  fingerprint_tags.dataset_id=tags.dataset_id AND
  fingerprint_tags.tag_type_id=tags.tag_type_id AND
  fingerprint_tags.tag_id=tags.tag_id;`
