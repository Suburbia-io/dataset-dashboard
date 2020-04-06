# Tagging app endpoint

### Endpoint
`POST https://dashboard.01.suburbia.io/admin/api/fingerprintTagUpsertCSV`

### Mandatory Query Params
`datasetID`: the Dataset ID you want to insert tags into

`tagAppID`: the ID of your Tagging App

### Authentication Header

`Authorization` : `Bearer {your api key}`

### Request Body (csv format)
`fingerprint,tag_type,tag,confidence`

one per row, header row mandatory. Content type: text/plain

### Example

```
curl --location --request POST 'https://dashboard.01.suburbia.io/admin/api/fingerprintTagUpsertCSV?datasetID=005da436-0ab6-5545-982f-4637780256f0&tagAppID=ae53e18d-33f5-4701-9e61-a5feb10b10d6' \
--header 'Content-Type: text/plain' \
--header 'Authorization: Bearer qzj6orDc7kojRGKWQ6RW7tGlwcjTe5EX' \
--data-raw 'fingerprint,tag_type,tag,confidence
	    3423423423234,brand,coca-cola,0.9'
```
