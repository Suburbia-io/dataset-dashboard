package migrations

// language=SQL
const AddDatasetRelationships = `
	ALTER TABLE corporations 
	ADD COLUMN dataset_entity UUID REFERENCES datasets(dataset_id);

	ALTER TABLE brands 
	ADD COLUMN dataset_entity UUID REFERENCES datasets(dataset_id);

	ALTER TABLE brand_categories
	ADD COLUMN dataset_entity UUID REFERENCES datasets(dataset_id);

	ALTER TABLE ixrules
	ADD COLUMN dataset_entity UUID REFERENCES datasets(dataset_id);

	ALTER TABLE corp_mappings
	ADD COLUMN dataset_entity UUID REFERENCES datasets(dataset_id);
`
