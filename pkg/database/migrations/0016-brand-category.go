package migrations

// language=SQL
const AddCategoryIDToBrands = `
ALTER TABLE brands 
  ADD COLUMN category_id UUID REFERENCES brand_categories(category_id);
`
