package migrations

// language=SQL
const DeleteCPGLocations = `
DELETE FROM locations WHERE dataset_id = '005da436-0ab6-5545-982f-4637780256f0';
`
