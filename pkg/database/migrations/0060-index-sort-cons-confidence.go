package migrations

// language=SQL
const IndexSortBrandConsConfidence = `
CREATE INDEX ON consensus_tags(dataset_id,tag_type_id,confidence);
`
