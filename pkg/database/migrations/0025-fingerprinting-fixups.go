package migrations

const DataFingerprintingUpdates = `

-- Drop unused tables.

DROP TABLE dataset_tag_types CASCADE;
DROP TABLE fingerprints CASCADE;
DROP TABLE fingerprint_tags CASCADE;

CREATE TABLE tag_types(
  dataset_id  UUID NOT NULL REFERENCES datasets(dataset_id),
  tag_type_id UUID NOT NULL UNIQUE,
  tag_type    TEXT NOT NULL,
  description TEXT NOT NULL,

  PRIMARY KEY(dataset_id, tag_type_id)
);

CREATE TABLE tags(
  dataset_id   UUID NOT NULL,
  tag_type_id  UUID NOT NULL,
  tag_id       UUID NOT NULL UNIQUE,
  tag          TEXT NOT NULL,
  description  TEXT NOT NULL,

  PRIMARY KEY (dataset_id, tag_type_id, tag_id),

  FOREIGN KEY (dataset_id, tag_type_id)
    REFERENCES tag_types(dataset_id, tag_type_id)
);

CREATE TABLE fingerprints(
  dataset_id     UUID   NOT NULL REFERENCES datasets(dataset_id),
  fingerprint    TEXT   NOT NULL, -- Unique by definition.
  raw_text       TEXT   NOT NULL DEFAULT '',
  annotations    TEXT   NOT NULL DEFAULT '',
  count          BIGINT NOT NULL DEFAULT 0,

  PRIMARY KEY(dataset_id, fingerprint)
);

CREATE TABLE fingerprint_tags(
  dataset_id         UUID NOT NULL,
  fingerprint        TEXT NOT NULL,
  tag_type_id        UUID NOT NULL,
  tag_id             UUID NOT NULL,

  PRIMARY KEY(dataset_id, fingerprint, tag_type_id, tag_id),

  FOREIGN KEY(dataset_id, tag_type_id, tag_id)
    REFERENCES tags(dataset_id, tag_type_id, tag_id),

  FOREIGN KEY(dataset_id, fingerprint)
    REFERENCES fingerprints(dataset_id, fingerprint)
);`
