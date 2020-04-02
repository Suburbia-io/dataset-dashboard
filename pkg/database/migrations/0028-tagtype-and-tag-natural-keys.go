package migrations

// language=SQL
const TagTypeAndTagNaturalKeys = `
CREATE UNIQUE INDEX ON tag_types(dataset_id, tag_type);
CREATE UNIQUE INDEX ON tags(dataset_id, tag_type_id, tag);
`
