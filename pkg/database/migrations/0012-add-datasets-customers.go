package migrations

// language=SQL
const CreateDatasetsAndCustomersTable = `

CREATE TABLE datasets (
  dataset_id       	UUID PRIMARY KEY,
  name     			TEXT NOT NULL UNIQUE,
  slug 	    		TEXT NOT NULL UNIQUE,
  created_at        TIMESTAMPTZ NOT NULL,
  updated_at        TIMESTAMPTZ NOT NULL,
  archived_at       TIMESTAMPTZ
);

CREATE TABLE customers (
  customer_id       UUID PRIMARY KEY,
  name     			TEXT NOT NULL UNIQUE,
  created_at        TIMESTAMPTZ NOT NULL,
  updated_at        TIMESTAMPTZ NOT NULL,
  archived_at       TIMESTAMPTZ
);

CREATE TABLE customer_datasets (
  customer_dataset_id	UUID PRIMARY KEY,
  created_at        	TIMESTAMPTZ NOT NULL,
  customer_entity 	UUID REFERENCES customers(customer_id),
  dataset_entity	UUID REFERENCES datasets(dataset_id)
);

CREATE UNIQUE INDEX customer_dataset_entities
  ON customer_datasets (dataset_entity,customer_entity);

`
