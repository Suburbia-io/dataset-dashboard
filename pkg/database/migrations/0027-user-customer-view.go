package migrations

// language=SQL
const CreateUserCustomerView = `
CREATE OR REPLACE VIEW user_view AS SELECT
  users.*,
  customers.name AS customer_name
FROM users LEFT JOIN customers ON users.customer_entity=customers.customer_id;`
