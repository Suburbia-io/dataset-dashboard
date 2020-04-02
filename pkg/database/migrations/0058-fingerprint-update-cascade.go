package migrations

const FingerprintUpdateCascade = `

ALTER TABLE consensus_tags
  DROP CONSTRAINT consensus_tags_dataset_id_fkey1;

ALTER TABLE consensus_tags
  ADD FOREIGN KEY(dataset_id, fingerprint)
  REFERENCES fingerprints(dataset_id, fingerprint)
  ON UPDATE CASCADE;

ALTER TABLE fingerprint_tag_historicals
  DROP CONSTRAINT fingerprint_tag_historicals_dataset_id_fkey1;

ALTER TABLE fingerprint_tags
  DROP CONSTRAINT fingerprint_tags_dataset_id_fkey1;

ALTER TABLE tag_app_tags
  DROP CONSTRAINT tag_app_tags_dataset_id_fkey1;

ALTER TABLE tag_app_tags
  ADD FOREIGN KEY(dataset_id, fingerprint)
  REFERENCES fingerprints(dataset_id, fingerprint)
  ON UPDATE CASCADE;

ALTER TABLE tag_app_historical_tags
  ADD FOREIGN KEY(dataset_id, fingerprint)
  REFERENCES fingerprints(dataset_id, fingerprint)
  ON UPDATE CASCADE;
`
