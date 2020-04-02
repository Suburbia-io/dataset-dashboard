package migrations

// language=SQL
const AlterEventsAddOnCustomer = `
ALTER TABLE events
ADD COLUMN on_customer UUID REFERENCES customers(customer_id);
`
