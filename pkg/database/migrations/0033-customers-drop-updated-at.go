package migrations

// language=SQL
const CustomersDropUpdatedAt = `
ALTER TABLE customers
  DROP COLUMN updated_at;
`
