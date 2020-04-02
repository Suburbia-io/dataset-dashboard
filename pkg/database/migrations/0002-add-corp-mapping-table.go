package migrations

// language=SQL
const CreateCorpMappingTable = `
CREATE TABLE corp_mappings (
  mapping_id    UUID PRIMARY KEY,
  brand_entity  UUID REFERENCES brands (brand_id),
  corp_entity   UUID REFERENCES corporations (corporation_id),
  corp_parent   UUID NOT NULL REFERENCES corporations (corporation_id),
  location      TEXT NOT NULL DEFAULT '',
  start_date    TIMESTAMPTZ NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL,
  updated_at    TIMESTAMPTZ NOT NULL,
  archived_at   TIMESTAMPTZ
);

ALTER TABLE corp_mappings
  ADD CONSTRAINT corp_mapping_brand_corp_only_one CHECK ((brand_entity IS NULL) <> (corp_entity IS NULL));

CREATE INDEX corp_mapping_brand_time_loc ON corp_mappings (brand_entity, start_date, location);
CREATE INDEX corp_mapping_parent_time_loc ON corp_mappings (corp_parent, start_date, location);
CREATE INDEX corp_mapping_corp_time_loc ON corp_mappings (corp_entity, start_date, location);

ALTER TABLE events
  ADD COLUMN on_corpmapping UUID REFERENCES corp_mappings(mapping_id);

CREATE INDEX events_on_corpmapping_index
  ON events (on_corpmapping,name);
`
