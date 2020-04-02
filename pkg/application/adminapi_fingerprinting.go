package application

import (
	"encoding/csv"
	"net/http"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
	"github.com/Suburbia-io/dashboard/shared"
)

func (app *App) AdminApiFingerprintUpsertCSV(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	datasetID := r.URL.Query().Get("datasetID")
	if err := validate.UUID(datasetID); err != nil {
		return err
	}

	reader := csv.NewReader(r.Body)

	if err := app.DBAL.FingerprintUpsertCSV(datasetID, reader); err != nil {
		return errors.Unexpected.
			Wrap("Failed to upsert fingerprint CSV: %w", err)
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Fingerprints.Table(), s.UserID, "fingerprintUpsertCSV", "")

	app.respondApi(w, r, true, nil)
	return nil
}

func (app *App) AdminApiFingerprintingListCorps(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	var args database.CorporationListArgs
	if err := app.decodeRequest(r, &args); err != nil {
		return err
	}

	corps, err := app.DBAL.CorporationList(args)
	if err != nil {
		app.respondApi(w, r, nil, err)
		return nil
	}
	app.respondApi(w, r, corps, nil)
	return nil
}

func (app *App) AdminApiFingerprintAnnotationsUpsertCSV(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	datasetID := r.URL.Query().Get("datasetID")
	reader := csv.NewReader(r.Body)
	if err := app.DBAL.FingerprintAnnotationsUpsertCSV(datasetID, reader); err != nil {
		return errors.HttpBadRequestArgs
	}
	app.respondApi(w, r, true, nil)
	return nil
}

func (app *App) AdminApiFingerprintTagUpsertCSV(
	w http.ResponseWriter,
	r *http.Request,
	s shared.SharedSession,
) error {
	datasetID := r.URL.Query().Get("datasetID")
	if err := validate.UUID(datasetID); err != nil {
		return err
	}

	tagAppID := r.URL.Query().Get("tagAppID")
	if err := validate.UUID(tagAppID); err != nil {
		return err
	}

	reader := csv.NewReader(r.Body)

	err := app.DBAL.TagAppTagUpsertCSV(
		datasetID,
		tagAppID,
		s.User.UserID,
		reader)
	if err != nil {
		return errors.Unexpected.
			Wrap("Failed to upsert tag app tag CSV: %w", err)
	}

	app.respondApi(w, r, true, nil)
	return nil
}

type FingerprintUpsertTagArgs struct {
	DatasetID    string   `json:"datasetID"`
	TagTypeID    string   `json:"tagTypeID"`
	TagValue     *string  `json:"tagValue"`
	Fingerprints []string `json:"fingerprints"`
}

func (app *App) AdminApiFingerprintUpsertTags(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	var args FingerprintUpsertTagArgs
	err := app.decodeRequest(r, &args)
	if err != nil {
		return err
	}

	// Don't allow users to update the CPG EU nested_category tags.
	// TODO Remove this when nested_categories have been removed.
	if args.TagTypeID == "005e15a3-6149-7d63-4f4b-6f9c7d86db8d" {
		err = errors.UpsertNotAllowed.Wrap("tag nested_category cannot be updated")
		app.respondApi(w, r, []database.FPTagRow{}, err)
		return nil
	}

	// Get the specific tag value.
	var tag database.Tag
	if args.TagValue != nil {
		tag, err = app.DBAL.TagLookup(args.DatasetID, args.TagTypeID, *args.TagValue)
		if err != nil {
			return err
		}
	}

	confidence := 0.9
	if len(args.Fingerprints) > 1 {
		confidence = 0.8
	}

	response := []database.FPTagRow{}
	for _, fp := range args.Fingerprints {
		if args.TagValue == nil {
			err := app.DBAL.TagAppTagDelete(args.DatasetID, fp, args.TagTypeID, database.TAG_APP_HUMAN_POOL)
			if err != nil {
				return err
			}
		} else {
			appTag := database.TagAppTag{
				DatasetID:   args.DatasetID,
				Fingerprint: fp,
				TagTypeID:   args.TagTypeID,
				TagAppID:    database.TAG_APP_HUMAN_POOL,
				TagID:       tag.TagID,
				Confidence:  confidence,
				UserID:      s.User.UserID,
			}

			err = app.DBAL.TagAppTagUpsert(&appTag)
			if err != nil {
				return err
			}
		}

		fpTag, err := app.DBAL.FPTagViewGet(
			args.DatasetID,
			database.TAG_APP_HUMAN_POOL,
			fp)
		if err != nil {
			return err
		}
		response = append(response, fpTag)
	}

	app.respondApi(w, r, response, nil)
	return nil
}

type TagListArgs struct {
	database.TagListArgs
	DatasetID string `json:"datasetID"`
	TagTypeID string `json:"tagTypeID"`
}

func (app *App) AdminApiFingerprintTagSuggestions(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	var args TagListArgs
	if err := app.decodeRequest(r, &args); err != nil {
		return err
	}

	suggestions, err := app.DBAL.TagList(args.DatasetID, args.TagTypeID, args.TagListArgs)
	if err != nil {
		return err
	}

	suggestionList := []string{}
	for _, suggestion := range suggestions {
		suggestionList = append(suggestionList, suggestion.Tag)
	}

	app.respondApi(w, r, suggestionList, nil)
	return nil
}

func (app *App) AdminApiFingerprintListTags(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	var args database.TagTypeListArgs
	if err := app.decodeRequest(r, &args); err != nil {
		return err
	}

	tagList, err := app.DBAL.TagTypeList(args)
	if err != nil {
		return err
	}

	app.respondApi(w, r, tagList, nil)
	return nil
}

func (app *App) AdminApiFingerprintListCorpTypes(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	var args database.CorporationTypeListArgs
	if err := app.decodeRequest(r, &args); err != nil {
		return err
	}

	tagList, err := app.DBAL.CorporationTypeList(args)
	if err != nil {
		return err
	}

	app.respondApi(w, r, tagList, nil)
	return nil
}

type FingerprintInsertCorpMappingArg struct {
	DatasetID           string            `json:"datasetID"`
	CorpTypeID          string            `json:"corpTypeID"`
	CorpID              string            `json:"corpID"`
	Location            string            `json:"location"`
	FromDate            string            `json:"fromDate"`
	FingerprintIncludes string            `json:"fingerprintIncludes"` // Prefix search (useful for EANs).
	FingerprintExcludes string            `json:"fingerprintExcludes"`
	TagsIncludes        map[string]string `json:"tagIncludes"`
	TagsExcludes        map[string]string `json:"tagExcludes"`
}

func (app *App) AdminApiFingerprintList(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	args := database.FPTagViewQuery{}
	if err := app.decodeRequest(r, &args); err != nil {
		return err
	}

	if args.Limit > 200 {
		args.Limit = 200
	}

	resp, err := app.DBAL.FPTagViewList(&args)
	if err != nil {
		return err
	}

	app.respondApi(w, r, resp, nil)
	return nil
}

func (app *App) AdminApiFingerprintListTagApps(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	archived := false
	apps, err := app.DBAL.TagAppList(database.TagAppListArgs{Archived: &archived})
	if err != nil {
		return err
	}
	app.respondApi(w, r, apps, nil)
	return nil
}
