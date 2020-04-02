package migrations

// language=SQL
const TagJobAndTaskRemove = `
DROP TABLE tag_tasks;
DROP TABLE tag_jobs;
`
