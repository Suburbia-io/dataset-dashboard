package database

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

var fpTagSelectViewTmpl = template.Must(template.New("").Parse(`
DROP VIEW IF EXISTS fingerprint_view_{{.TableID}};

CREATE VIEW fingerprint_view_{{.TableID}} AS SELECT
  fps.fingerprint,
  fps.raw_text,
  fps.annotations,
  fps.count,
  tag_apps.tag_app_id

 {{range .TagTypes}},

  -- {{.TagType}}: app
  {{.TagType}}_tags.tag_id AS {{.TagType}}_tag_id,
  {{.TagType}}_tag_names.tag AS {{.TagType}}_tag,
  {{.TagType}}_tags.confidence AS {{.TagType}}_confidence,
  {{.TagType}}_tags.updated_at AS {{.TagType}}_updated_at,

  -- {{.TagType}}: consensus
  {{.TagType}}_cons_tags.tag_id AS {{.TagType}}_cons_tag_id,
  {{.TagType}}_cons_tag_names.tag AS {{.TagType}}_cons_tag,
  {{.TagType}}_cons_tags.confidence AS {{.TagType}}_cons_confidence,
  {{.TagType}}_cons_tags.updated_at AS {{.TagType}}_cons_updated_at

  {{end}}

FROM fingerprints AS fps

LEFT JOIN tag_apps ON TRUE

{{range .TagTypes}}
-- {{.TagType}}: app
LEFT JOIN tag_app_tags AS {{.TagType}}_tags ON
  {{.TagType}}_tags.dataset_id='{{$.DatasetID}}' AND
  {{.TagType}}_tags.tag_type_id='{{.TagTypeID}}' AND
  {{.TagType}}_tags.tag_app_id=tag_apps.tag_app_id AND
  {{.TagType}}_tags.fingerprint=fps.fingerprint

LEFT JOIN tags AS {{.TagType}}_tag_names ON
  {{.TagType}}_tag_names.dataset_id='{{$.DatasetID}}' AND
  {{.TagType}}_tag_names.tag_type_id='{{.TagTypeID}}' AND
  {{.TagType}}_tag_names.tag_id={{.TagType}}_tags.tag_id

-- {{.TagType}}: consensus
LEFT JOIN consensus_tags AS {{.TagType}}_cons_tags ON
  {{.TagType}}_cons_tags.dataset_id='{{$.DatasetID}}' AND
  {{.TagType}}_cons_tags.fingerprint=fps.fingerprint AND
  {{.TagType}}_cons_tags.tag_type_id='{{.TagTypeID}}'

LEFT JOIN tags AS {{.TagType}}_cons_tag_names ON
  {{.TagType}}_cons_tag_names.dataset_id='{{$.DatasetID}}' AND
  {{.TagType}}_cons_tag_names.tag_type_id='{{.TagTypeID}}' AND
  {{.TagType}}_cons_tag_names.tag_id={{.TagType}}_cons_tags.tag_id

{{end}}

WHERE fps.dataset_id='{{.DatasetID}}';`))

type fpTagSelectViewCtx struct {
	DatasetID string
	TableID   string
	TagTypes  []TagType
}

// Must update view when:
//
// (1) Dataset upsert
// (2) TagType upsert, deleted
func (db *DBAL) createFPTagSelectView(tx DBi, datasetID string) error {
	tableID := strings.ReplaceAll(datasetID, "-", "_")

	// If dataset is archived or not managable, drop the view.
	create := true
	err := tx.QueryRow(
		`SELECT archived_at IS NULL AND manageable FROM datasets `+
			`WHERE dataset_id=$1`, datasetID).Scan(&create)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to determine if fingerprint view is required: %w", err).
			Alert()
	}

	if !create {
		_, err := tx.Exec(`DROP VIEW IF EXISTS fingerprint_view_` + tableID)
		if err != nil {
			return errors.Unexpected.
				Wrap("Failed to drop view: %w", err).
				Alert()
		}
		return nil
	}

	// The dataset is not archived and manageable => create view.
	tts, err := db.tagTypeList(tx, TagTypeListArgs{DatasetID: datasetID})
	if err != nil {
		return err
	}

	ctx := fpTagSelectViewCtx{
		DatasetID: datasetID,
		TableID:   tableID,
		TagTypes:  tts,
	}

	buf := &bytes.Buffer{}
	if err := fpTagSelectViewTmpl.Execute(buf, ctx); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create view template: %w", err).
			Alert()
	}

	if _, err := tx.Exec(buf.String()); err != nil {
		return errors.Unexpected.
			Wrap("Failed to create view: %w", err).
			Alert()
	}

	return nil
}

// ----------------------------------------------------------------------------

type FPTagRow struct {
	Fingerprint             string       `json:"fingerprint"`
	RawText                 string       `json:"rawText"`
	Annotations             string       `json:"annotations"`
	Count                   int          `json:"count"`
	TagIDs                  []*string    `json:"tagIDs"`
	Tags                    []*string    `json:"tags"`
	Confidences             []*float64   `json:"confidences"`
	UpdatedAts              []*time.Time `json:"updatedAts"`
	UpdatedAtsFormatted     []string     `json:"updatedAtsFormatted"`
	ConsTagIDs              []*string    `json:"consTagIDs"`
	ConsTags                []*string    `json:"consTags"`
	ConsConfidences         []*float64   `json:"consConfidences"`
	ConsUpdatedAts          []*time.Time `json:"consUpdatedAts"`
	ConsUpdatedAtsFormatted []string     `json:"consUpdatedAtsFormatted"`
}

type FPTagViewQuery struct {
	DatasetID string `json:"datasetID"` // Required.
	TagAppID  string `json:"tagAppID"`  // Required.

	Limit  int `json:"limit"`
	Offset int `json:"offset"`

	FingerprintIncludes []string `json:"fingerprintIncludes"`
	FingerprintExcludes []string `json:"fingerprintExcludes"`

	RawTextIncludes []string `json:"rawTextIncludes"`
	RawTextExcludes []string `json:"rawTextExcludes"`

	// OrderBy should be one of
	//
	// count (default)
	// fingerprint
	// <tag>_confidence
	// <tag>_updated_at
	// <tag>_cons_confidence
	// <tag>_cons_updated_at
	OrderBy  string `json:"orderBy"`
	OrderAsc bool   `json:"orderAsc"`

	TagIncludes map[string][]string `json:"tagIncludes"`
	TagExcludes map[string][]string `json:"tagExcludes"`

	ConsTagIncludes map[string][]string `json:"consTagIncludes"`
	ConsTagExcludes map[string][]string `json:"consTagExcludes"`

	CountThreshold int `json:"countThreshold"`

	// TODO: Possibility: user/consensus confidence min/max

	tagTypes []TagType
	ttMap    map[string]TagType
}

func (q *FPTagViewQuery) TagTypes(db DBi) ([]TagType, error) {
	if err := q.init(db); err != nil {
		return nil, err
	}
	return q.tagTypes, nil
}

func (q *FPTagViewQuery) init(db DBi) error {
	// Only initialize once.
	if q.tagTypes != nil {
		return nil
	}

	// Get list of tag types.
	b := qb.Builder{}
	b.Write(`SELECT`)
	b.Write(tables.TagTypes.SelectCols())
	b.Write(`FROM`)
	b.Write(tables.TagTypes.Table())
	b.Write(`WHERE dataset_id=$1`, q.DatasetID)
	b.Write(`ORDER BY tag_type ASC`)
	ttListQuery, ttListArgs := b.MustBuild()
	tts, err := tables.TagTypes.List(db, ttListQuery, ttListArgs...)
	if err != nil {
		return err
	}
	q.tagTypes = tts

	q.ttMap = map[string]TagType{}
	for _, tt := range tts {
		q.ttMap[tt.TagType] = tt
	}

	allowedOrderBy := map[string]bool{
		"count":       true,
		"fingerprint": true,
	}

	allowedTags := map[string]bool{}

	for _, tt := range tts {
		allowedOrderBy[tt.TagType+"_confidence"] = true
		allowedOrderBy[tt.TagType+"_updated_at"] = true
		allowedOrderBy[tt.TagType+"_cons_confidence"] = true
		allowedOrderBy[tt.TagType+"_cons_updated_at"] = true
		allowedTags[tt.TagType] = true
	}

	if q.OrderBy == "" {
		q.OrderBy = "count"
	}
	if _, ok := allowedOrderBy[q.OrderBy]; !ok {
		return errors.DBUnknownColumn.Wrap("Column was: %s", q.OrderBy)
	}

	maps := []map[string][]string{
		q.TagIncludes,
		q.TagExcludes,
		q.ConsTagIncludes,
		q.ConsTagExcludes,
	}
	for _, m := range maps {
		for tag := range m {
			if _, ok := allowedTags[tag]; !ok {
				return errors.DBUnknownColumn.Wrap("Column was: %s", tag)
			}
		}
	}

	return nil
}

func (q FPTagViewQuery) viewName() string {
	return "fingerprint_view_" + strings.ReplaceAll(q.DatasetID, "-", "_")
}

func (q FPTagViewQuery) selectCols() string {
	l := []string{
		"fingerprint", "raw_text", "annotations", "count",
	}

	for _, tt := range q.tagTypes {
		l = append(l,
			tt.TagType+"_tag_id",
			tt.TagType+"_tag",
			tt.TagType+"_confidence",
			tt.TagType+"_updated_at",
			tt.TagType+"_cons_tag_id",
			tt.TagType+"_cons_tag",
			tt.TagType+"_cons_confidence",
			tt.TagType+"_cons_updated_at",
		)
	}

	return strings.Join(l, ",")
}

func (q FPTagViewQuery) scan(
	rows interface {
		Scan(...interface{}) error
	},
) (row FPTagRow, err error) {
	N := len(q.tagTypes)
	scanArgs := make([]interface{}, 0, 4+8*N)

	row.TagIDs = make([]*string, N)
	row.Tags = make([]*string, N)
	row.Confidences = make([]*float64, N)
	row.UpdatedAts = make([]*time.Time, N)
	row.UpdatedAtsFormatted = make([]string, N)
	row.ConsTagIDs = make([]*string, N)
	row.ConsTags = make([]*string, N)
	row.ConsConfidences = make([]*float64, N)
	row.ConsUpdatedAts = make([]*time.Time, N)
	row.ConsUpdatedAtsFormatted = make([]string, N)

	scanArgs = append(scanArgs,
		&row.Fingerprint,
		&row.RawText,
		&row.Annotations,
		&row.Count)

	for i := range q.tagTypes {
		scanArgs = append(scanArgs,
			&row.TagIDs[i],
			&row.Tags[i],
			&row.Confidences[i],
			&row.UpdatedAts[i],
			&row.ConsTagIDs[i],
			&row.ConsTags[i],
			&row.ConsConfidences[i],
			&row.ConsUpdatedAts[i])
	}

	err = rows.Scan(scanArgs...)
	if err != nil {
		return row, errors.Unexpected.
			Wrap("Failed to scan FPTagRow: %w", err).
			Alert()
	}

	// Formatted dates.
	row.UpdatedAtsFormatted = make([]string, len(row.UpdatedAts))
	for i, date := range row.UpdatedAts {
		if date != nil {
			row.UpdatedAtsFormatted[i] = date.Format("2006-01-02 15:04")
		}
	}

	row.ConsUpdatedAtsFormatted = make([]string, len(row.ConsUpdatedAts))
	for i, date := range row.ConsUpdatedAts {
		if date != nil {
			row.ConsUpdatedAtsFormatted[i] = date.Format("2006-01-02 15:04")
		}
	}

	return row, nil
}

// Returns builder containing:
//
// `WHERE <conditions>`
func (q FPTagViewQuery) buildWhere() (string, []interface{}) {
	b := qb.Builder{}

	b.Write("WHERE tag_app_id=$1", q.TagAppID)

	// filter by minimum amount of rows in the real dataset, aka count
	b.Write("AND count >= $1", q.CountThreshold)

	// Sorting the list by confidence is slow when the NULLs are filtered out.
	if strings.HasSuffix(q.OrderBy, "_cons_confidence") {
		b.Write(fmt.Sprintf("AND %s IS NOT NULL", q.OrderBy))
	}

	for _, s := range q.FingerprintIncludes {
		if s != "" {
			b.Write("AND fingerprint LIKE $1", s)
		}
	}

	for _, s := range q.FingerprintExcludes {
		if s != "" {
			b.Write("AND fingerprint NOT LIKE $1", s)
		}
	}

	for _, s := range q.RawTextIncludes {
		if s != "" {
			b.Write("AND (raw_text LIKE $1 OR annotations LIKE $1)", s)
		}
	}

	for _, s := range q.RawTextExcludes {
		if s != "" {
			b.Write("AND NOT (raw_text LIKE $1 OR annotations LIKE $1)", s)
		}
	}

	for tag, includes := range q.TagIncludes {
		for _, s := range includes {
			switch s {
			case "NULL":
				// This seems nuts, but it's several orders of magnitude faster than
				// the simple approach (like 50ms vs 40 seconds).
				// -jdl 2020-01
				b.Write(`AND NOT EXISTS(`)
				b.Write(`SELECT tag_id`)
				b.Write(`FROM tag_app_tags`)
				b.Write(`WHERE dataset_id=$1`, q.DatasetID)
				b.Write(`AND tag_type_id=$1`, q.ttMap[tag].TagTypeID)
				b.Write(`AND tag_app_id=$1`, q.TagAppID)
				b.Write(`AND fingerprint=` + q.viewName() + `.fingerprint`)
				b.Write(`)`)
			case "":
				// Nothing.
			default:
				b.Write("AND "+tag+"_tag LIKE $1", s)
			}
		}
	}

	for tag, excludes := range q.TagExcludes {
		colName := tag + "_tag"
		for _, s := range excludes {
			switch s {
			case "NULL":
				b.Write("AND " + colName + "_id IS NOT NULL")
			case "":
				// Nothing.
			default:
				b.Write("AND "+colName+" NOT LIKE $1", s)
			}
		}
	}

	for tag, includes := range q.ConsTagIncludes {
		for _, s := range includes {
			switch s {
			case "NULL":
				b.Write(`AND NOT EXISTS(`)
				b.Write(`SELECT tag_id`)
				b.Write(`FROM consensus_tags`)
				b.Write(`WHERE dataset_id=$1`, q.DatasetID)
				b.Write(`AND tag_type_id=$1`, q.ttMap[tag].TagTypeID)
				b.Write(`AND fingerprint=` + q.viewName() + `.fingerprint`)
				b.Write(`)`)
			case "":
				// Nothing.
			default:
				b.Write("AND "+tag+"_cons_tag LIKE $1", s)
			}
		}
	}

	for tag, excludes := range q.ConsTagExcludes {
		colName := tag + "_cons_tag"
		for _, s := range excludes {
			switch s {
			case "NULL":
				b.Write("AND " + colName + "_id IS NOT NULL")
			case "":
				// Nothing.
			default:
				b.Write("AND "+colName+" NOT LIKE $1", s)
			}
		}
	}

	return b.MustBuild()
}

// ----------------------------------------------------------------------------

// Get a FPTagRow for a single fingerprint.
func (q *FPTagViewQuery) get(
	db DBi,
	fingerprint string,
) (
	row FPTagRow,
	err error,
) {
	q.Limit = 1
	q.Offset = 0

	if err = q.init(db); err != nil {
		return row, err
	}

	b := qb.Builder{}
	b.Write(`SELECT`)
	b.Write(q.selectCols())
	b.Write(`FROM`)
	b.Write(q.viewName())
	b.Write(`WHERE fingerprint=$1`, fingerprint)
	b.Write(`AND tag_app_id=$1`, q.TagAppID)

	query, qArgs := b.MustBuild()
	return q.scan(db.QueryRow(query, qArgs...))
}

// ----------------------------------------------------------------------------

// Iterate over matching rows.
func (q *FPTagViewQuery) iterate(
	db DBi,
	onRow func(row FPTagRow) error,
) error {
	if err := q.init(db); err != nil {
		return err
	}

	if q.Offset < 0 {
		return errors.IllegalOffset
	}

	whereClause, whereArgs := q.buildWhere()

	// Create list query.
	b := qb.Builder{}
	b.Write(`SELECT`)
	b.Write(q.selectCols())
	b.Write(`FROM`)
	b.Write(q.viewName())
	b.Write(whereClause, whereArgs...)

	b.Write(`ORDER BY`)
	b.Write(q.OrderBy)

	if q.OrderAsc {
		b.Write("ASC")
	} else {
		b.Write("DESC")
	}

	if q.Limit > 0 {
		b.Write(fmt.Sprintf("LIMIT %d OFFSET %d", q.Limit+1, q.Offset))
	}

	listQuery, listArgs := b.MustBuild()

	rows, err := db.Query(listQuery, listArgs...)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to list fingerprints: %w", err).
			Alert()
	}
	defer rows.Close()

	for rows.Next() {
		row, err := q.scan(rows)
		if err != nil {
			return err
		}
		if err := onRow(row); err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------------

func (db *DBAL) FPTagViewGet(
	datasetID string,
	tagAppID string,
	fingerprint string,
) (
	FPTagRow,
	error,
) {
	q := &FPTagViewQuery{
		DatasetID: datasetID,
		TagAppID:  tagAppID,
	}

	return q.get(db, fingerprint)
}

type FPTagViewListResp struct {
	// DatasetID and TagAppID determine which tags are returned.
	DatasetID string `json:"datasetID"`
	TagAppID  string `json:"tagAppID"`

	// TagTypes is the list of valid tag types for this dataset.
	TagTypes []TagType `json:"tagTypes"`

	// Rows contains the result of the list query.
	Rows []FPTagRow `json:"rows"`

	More bool `json:"more"`

	// TotalCount contains the total number of rows in the dataset represented by
	// the fingerprints in the returned list.
	TotalCount int `json:"totalCount"`
}

func (db *DBAL) FPTagViewList(
	q *FPTagViewQuery,
) (
	resp FPTagViewListResp,
	err error,
) {
	if q.Limit < 1 {
		return resp, errors.IllegalLimit
	}
	if q.Offset < 0 {
		return resp, errors.IllegalOffset
	}

	resp.DatasetID = q.DatasetID
	resp.TagAppID = q.TagAppID
	resp.Rows = []FPTagRow{}

	q.Limit += 1
	err = q.iterate(db, func(row FPTagRow) error {
		resp.TotalCount += row.Count
		resp.Rows = append(resp.Rows, row)
		return nil
	})
	if err != nil {
		return resp, err
	}

	resp.TagTypes = q.tagTypes

	if len(resp.Rows) == q.Limit {
		resp.Rows = resp.Rows[:q.Limit-1]
		resp.More = true
	}

	return resp, nil
}
