package migrations

// language=SQL
const AddDataFingerprinting = `

	CREATE TABLE dataset_tag_types (
	  tag_type_id	  	UUID PRIMARY KEY,
	  dataset_id		UUID REFERENCES datasets(dataset_id),
	  tag_name			TEXT,
	  created_at    	TIMESTAMPTZ NOT NULL,
	  archived_at   	TIMESTAMPTZ
	);

	CREATE UNIQUE INDEX tag_types_unique ON dataset_tag_types(dataset_id, tag_name);

	CREATE TABLE dataset_corp_types (
	  corp_type_id	  	UUID PRIMARY KEY,
	  dataset_id		UUID REFERENCES datasets(dataset_id),
	  corp_type_name	TEXT,
	  created_at    	TIMESTAMPTZ NOT NULL,
	  archived_at   	TIMESTAMPTZ
	);

	CREATE UNIQUE INDEX corp_types_unique ON dataset_corp_types(dataset_id, corp_type_name);

	CREATE TABLE fingerprints (
		fingerprint_id		UUID PRIMARY KEY,
		dataset_id			UUID REFERENCES datasets(dataset_id),
		fingerprint			TEXT,
		raw_text			TEXT, 
		annotations			TEXT DEFAULT NULL,
		created_at    		TIMESTAMPTZ NOT NULL,
	  	archived_at   		TIMESTAMPTZ
	);

	CREATE INDEX fingerprint 
	  ON fingerprints(lower(fingerprint));

	CREATE UNIQUE INDEX fingerprint_unique
	ON fingerprints(dataset_id, lower(fingerprint));

	CREATE TABLE fingerprint_tags (
	    fingerprint_tag_id 		UUID PRIMARY KEY,
	    fingerprint_id			UUID REFERENCES fingerprints(fingerprint_id),
	    tag_type_id				UUID REFERENCES dataset_tag_types(tag_type_id),
	    tag						TEXT NOT NULL,
	  	created_at    			TIMESTAMPTZ NOT NULL,
	  	archived_at   			TIMESTAMPTZ
	);
	
	CREATE UNIQUE INDEX fingerprint_tag_unique 
	ON fingerprint_tags(fingerprint_id, tag_type_id);

	CREATE TABLE fingerprint_corps (
	    fingerprint_corp_id		UUID PRIMARY KEY,
	    corp_type_id			UUID REFERENCES dataset_corp_types(corp_type_id),
	    fingerprint_id			UUID REFERENCES fingerprints(fingerprint_id),
	    corp_id					UUID REFERENCES corporations(corporation_id)
	);
	
	CREATE UNIQUE INDEX fingerprint_corp_unique 
	ON fingerprint_corps(fingerprint_id, corp_type_id);

`
