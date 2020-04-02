package migrations

// language=SQL
const ChangeCorpsConstraint = `

DROP INDEX corporations_symbol;

CREATE UNIQUE INDEX corporations_symbol
ON corporations (exchange, code, dataset_id)
WHERE (btrim(exchange) <> ''::text) AND (btrim(code) <> ''::text);
`
