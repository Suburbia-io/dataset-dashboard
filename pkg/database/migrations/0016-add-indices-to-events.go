package migrations

// language=SQL
const AlterEventsAddOnCustomerOnDatasetIndices = `

CREATE INDEX events_on_dataset_index
  ON events (on_dataset,name);

CREATE INDEX events_on_customer_index
  ON events (on_customer,name);

`
