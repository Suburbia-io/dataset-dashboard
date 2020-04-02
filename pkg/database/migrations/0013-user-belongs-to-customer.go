package migrations

// language=SQL
const UserBelongsToCustomer = `
ALTER TABLE users
  ADD COLUMN customer_entity UUID REFERENCES customers(customer_id);
`
