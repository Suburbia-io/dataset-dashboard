package migrations

// language=SQL
const AlterBrandsTableAddPublished = `
ALTER TABLE brands
  ADD COLUMN published_at TIMESTAMPTZ DEFAULT NOW();

CREATE INDEX brands_published
  ON brands (published_at);
`
