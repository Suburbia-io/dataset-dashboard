package migrations

// language=SQL
const LocationsTableCleanups = `
-- Drop some stuff we're not using
DROP INDEX locations_checked;
ALTER TABLE locations
  DROP COLUMN checked_at;
ALTER TABLE locations
  DROP COLUMN archived_at;

-- Rename some stuff
ALTER TABLE locations
  RENAME COLUMN loc_id TO location_id;
ALTER TABLE locations
  RENAME COLUMN loc_hash TO location_hash;
ALTER TABLE locations
  RENAME COLUMN loc_string TO location_string;
ALTER TABLE locations
	RENAME CONSTRAINT locations_loc_hash_key TO locations_location_hash_key;

-- Migrate dataset_ids
ALTER TABLE locations
	ADD COLUMN dataset_id UUID;

UPDATE locations AS l
  SET dataset_id = d.dataset_id
FROM datasets AS d
WHERE l.location_string LIKE 'FRA,%'
	AND d.slug = 'perfcos';

UPDATE locations AS l
  SET dataset_id = d.dataset_id
FROM datasets AS d
WHERE l.dataset_id IS NULL
	AND d.slug = 'cpg';

ALTER TABLE locations
	ALTER COLUMN dataset_id SET NOT NULL;

-- Change to multi-column primary key using dataset_id and location_hash
ALTER TABLE locations
	DROP CONSTRAINT locations_pkey CASCADE;

ALTER TABLE locations
  ADD PRIMARY KEY (dataset_id, location_hash);

ALTER TABLE events
	ALTER COLUMN on_location SET DATA TYPE TEXT;

UPDATE events AS e
  SET on_location = l.location_hash
FROM locations AS l
  WHERE e.on_location::text = l.location_id::text AND e.on_location IS NOT NULL;

ALTER TABLE events
  ADD CONSTRAINT events_on_location_fkey FOREIGN KEY (on_location) REFERENCES locations(location_hash);

ALTER TABLE locations
  DROP COLUMN location_id;
`
