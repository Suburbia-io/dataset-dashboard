package migrations

// language=SQL
const CreateBrandCategoryTable = `
CREATE TABLE brand_categories(
  category_id UUID NOT NULL PRIMARY KEY,
  parent_id   UUID,
  name        TEXT NOT NULL,
  archived_at TIMESTAMPTZ,

  FOREIGN KEY (parent_id) REFERENCES brand_categories(category_id)
);

CREATE UNIQUE INDEX brand_categories_nk ON brand_categories(parent_id, name)
  WHERE parent_id IS NOT NULL;

CREATE UNIQUE INDEX brand_categories_nk2 ON brand_categories(name)
  WHERE parent_id IS NULL;
`
