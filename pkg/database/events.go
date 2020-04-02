package database

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

type EventName string

type Event struct {
	EventID   int64     `json:"eventID"`
	Timestamp time.Time `json:"timestamp"`

	// Subject:
	BySystem bool `json:"-"`

	ByAdmin             *string    `json:"byAdmin,omitempty"`
	ByAdminName         *string    `json:"byAdminName,omitempty"`
	ByAdminArchivedAt   *time.Time `json:"byAdminArchivedAt,omitempty"`
	ByAdminLastActiveAt *time.Time `json:"byAdminLastActiveAt,omitempty"`

	ByUser             *string    `json:"byUser,omitempty"`
	ByUserName         *string    `json:"byUserName,omitempty"`
	ByUserArchivedAt   *string    `json:"byUserArchivedAt,omitempty"`
	ByUserLastActiveAt *time.Time `json:"byUserLastActiveAt,omitempty"`

	// Object
	OnAdmin             *string    `json:"onAdmin,omitempty"`
	OnAdminName         *string    `json:"onAdminName,omitempty"`
	OnAdminArchivedAt   *time.Time `json:"onAdminArchivedAt,omitempty"`
	OnAdminLastActiveAt *time.Time `json:"onAdminLastActiveAt,omitempty"`

	OnUser             *string    `json:"onUser,omitempty"`
	OnUserName         *string    `json:"onUserName,omitempty"`
	OnUserArchivedAt   *string    `json:"onUserArchivedAt,omitempty"`
	OnUserLastActiveAt *time.Time `json:"onUserLastActiveAt,omitempty"`

	OnBrand           *string `json:"onBrand,omitempty"`
	OnBrandLabel      *string `json:"onBrandLabel,omitempty"`
	OnBrandArchivedAt *string `json:"onBrandArchivedAt,omitempty"`

	OnIXrule           *string `json:"onIXrule,omitempty"`
	OnIXruleName       *string `json:"onIXruleName,omitempty"`
	OnIXruleArchivedAt *string `json:"onIXruleArchivedAt,omitempty"`

	OnCorporation           *string `json:"onCorporation,omitempty"`
	OnCorporationName       *string `json:"onCorporationName,omitempty"`
	OnCorporationArchivedAt *string `json:"onCorporationArchivedAt,omitempty"`

	OnCorpMapping           *string `json:"onCorpMapping,omitempty"`
	OnCorpMappingName       *string `json:"onCorpMappingName,omitempty"`
	OnCorpMappingArchivedAt *string `json:"onCorpMappingArchivedAt,omitempty"`

	OnDataset           *string `json:"onDataset,omitempty"`
	OnDatasetName       *string `json:"onDatasetName,omitempty"`
	OnDatasetArchivedAt *string `json:"onDatasetArchivedAt,omitempty"`

	OnCustomer           *string `json:"onCustomer,omitempty"`
	OnCustomerName       *string `json:"onCustomerName,omitempty"`
	OnCustomerArchivedAt *string `json:"onCustomerArchivedAt,omitempty"`

	OnLocation *string `json:"onLocation,omitempty"`

	// Verb:
	Name       string        `json:"name"`
	Payload    []interface{} `json:"payload,omitempty"`
	AuthMethod *string       `json:"authMethod,omitempty"`
}

func (db *DBAL) EventCreate(event *Event) (err error) {
	if strings.TrimSpace(event.Name) == "" {
		return errors.EventInvalidName
	}
	event.Timestamp = time.Now().Round(time.Microsecond)

	payloadB, err := json.Marshal(event.Payload)
	if err != nil {
		return errors.UnexpectedError(err, "Failed marshalling events payload")
	}

	err = db.QueryRow(`
	INSERT INTO events (
		timestamp,
		by_system,
		by_admin,
		by_user,
		on_admin,
		on_user,
		on_brand,
		on_ixrule,
		on_corporation,
		on_corpmapping,
		on_location,
		on_dataset,
		on_customer,
		name,
		payload,
        auth_method
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING event_id`,
		event.Timestamp,
		event.BySystem,
		event.ByAdmin,
		event.ByUser,
		event.OnAdmin,
		event.OnUser,
		event.OnBrand,
		event.OnIXrule,
		event.OnCorporation,
		event.OnCorpMapping,
		event.OnLocation,
		event.OnDataset,
		event.OnCustomer,
		event.Name,
		payloadB,
		event.AuthMethod,
	).Scan(&event.EventID)
	if err != nil {
		return errors.UnexpectedError(err, "Failed inserting events")
	}

	return nil
}

type EventListArgs struct {
	Limit         int
	Offset        int
	EventType     string
	ForAdmin      string
	ForUser       string
	OnBrand       string
	OnIXrule      string
	OnCorporation string
	OnCorpMapping string
	OnLocation    string
	OnDataset     string
	OnCustomer    string
}

func (db *DBAL) EventList(args EventListArgs) (events []Event, err error) {
	if args.Limit < 0 {
		return events, errors.IllegalLimit
	}
	if args.Offset < 0 {
		return events, errors.IllegalOffset
	}

	var andWheres []string
	var whereParams []interface{}
	if args.ForAdmin != "" {
		if _, err := db.UserGet(args.ForAdmin); err != nil {
			return events, err
		}
		andWheres = append(andWheres, fmt.Sprintf("(by_admin=$%d OR on_admin=$%d)", len(whereParams)+1, len(whereParams)+1))
		whereParams = append(whereParams, args.ForAdmin)
	}

	if args.ForUser != "" {
		if _, err := db.UserGet(args.ForUser); err != nil {
			return events, err
		}
		andWheres = append(andWheres, fmt.Sprintf("(by_user=$%d OR on_user=$%d)", len(whereParams)+1, len(whereParams)+1))
		whereParams = append(whereParams, args.ForUser)
	}

	if args.EventType != "" {
		andWheres = append(andWheres, fmt.Sprintf("(e.name=$%d)", len(whereParams)+1))
		whereParams = append(whereParams, args.EventType)
	}

	if args.OnCorporation != "" {
		if _, err := db.CorporationGet(args.OnCorporation); err != nil {
			return events, err
		}
		andWheres = append(andWheres, fmt.Sprintf("(on_corporation=$%d)", len(whereParams)+1))
		whereParams = append(whereParams, args.OnCorporation)
	}

	if args.OnLocation != "" {
		if _, err := db.LocationGetByLocationHash(args.OnLocation); err != nil {
			return events, err
		}
		andWheres = append(andWheres, fmt.Sprintf("(on_location=$%d)", len(whereParams)+1))
		whereParams = append(whereParams, args.OnLocation)
	}

	if args.OnDataset != "" {
		if _, err := db.DatasetGet(args.OnDataset); err != nil {
			return events, err
		}
		andWheres = append(andWheres, fmt.Sprintf("(on_dataset=$%d)", len(whereParams)+1))
		whereParams = append(whereParams, args.OnDataset)
	}

	if args.OnCustomer != "" {
		if _, err := db.CustomerGet(args.OnCustomer); err != nil {
			return events, err
		}
		andWheres = append(andWheres, fmt.Sprintf("(on_customer=$%d)", len(whereParams)+1))
		whereParams = append(whereParams, args.OnCustomer)
	}

	whereFilter := ""
	if len(whereParams) > 0 {
		whereFilter = "WHERE " + strings.Join(andWheres, " AND ")
	}

	// language=SQL
	stmt := fmt.Sprintf(`SELECT
				e.event_id,
				timestamp,

				by_system,

				by_admin,
				by_a.name as by_admin_name,
				by_a.archived_at as by_admin_archived_at,
				by_a.last_active_at as by_admin_last_active_at,

				by_user,
				by_u.name as by_user_name,
				by_u.archived_at as by_user_archived_at,
				by_u.last_active_at as by_user_last_active_at,

				on_admin,
				on_a.name as on_admin_name,
				on_a.archived_at as on_admin_archived_at,
				on_a.last_active_at as on_admin_last_active_at,

				on_user,
				on_u.name as on_user_name,
				on_u.archived_at as on_user_archived_at,
				on_u.last_active_at as on_user_last_active_at,

				on_brand,
				on_b.label as on_brand_label,
				on_b.archived_at as on_brand_archived_at,

				on_ixrule,
				on_brand_rule.label,
				on_ix.archived_at as on_ixrule_archived_at,

				on_corporation,
				on_c.name as on_corporation_name,
				on_c.archived_at as on_corporation_archived_at,

				on_corpmapping,
				on_c.name as on_corpmapping_name,
				on_cm.archived_at as on_corpmapping_archived_at,

       			on_location,

       			on_dataset,
       			on_ds.name as on_dataset_name,
       			on_ds.archived_at as on_dataset_archived_at,

				on_customer,
       			on_cu.name as on_customer_name,
				on_cu.archived_at as on_customer_archived_at,

				e.name,
				payload,
       			e.auth_method
			FROM events e
			LEFT JOIN admins by_a ON by_a.admin_id=e.by_admin
			LEFT JOIN users by_u ON by_u.user_id=e.by_user
			LEFT JOIN admins on_a ON on_a.admin_id=e.on_admin
			LEFT JOIN users on_u ON on_u.user_id=e.on_user
			LEFT JOIN brands on_b ON on_b.brand_id=e.on_brand
			LEFT JOIN ixrules on_ix ON on_ix.rule_id=e.on_ixrule
			LEFT JOIN brands on_brand_rule ON (on_ix.tag_group='brand' AND on_ix.tag_id=on_brand_rule.brand_id)
			LEFT JOIN corporations on_c ON on_c.corporation_id=e.on_corporation
			LEFT JOIN corp_mappings on_cm ON on_cm.mapping_id=e.on_corpmapping
			LEFT JOIN locations on_loc ON on_loc.location_id=e.on_location
			LEFT JOIN datasets on_ds ON on_ds.dataset_id=e.on_dataset
			LEFT JOIN customers on_cu ON on_cu.customer_id=e.on_customer
			%s
			ORDER BY timestamp DESC
			LIMIT %d OFFSET %d
	`, whereFilter, args.Limit, args.Offset)

	rows, err := db.Query(stmt, whereParams...)
	if err != nil {
		return events, errors.UnexpectedError(err, "Failed listing events")
	}
	defer rows.Close()

	for rows.Next() {
		event := Event{}
		payload := ""

		if err := rows.Scan(
			&event.EventID,
			&event.Timestamp,

			&event.BySystem,

			&event.ByAdmin,
			&event.ByAdminName,
			&event.ByAdminArchivedAt,
			&event.ByAdminLastActiveAt,

			&event.ByUser,
			&event.ByUserName,
			&event.ByUserArchivedAt,
			&event.ByUserLastActiveAt,

			&event.OnAdmin,
			&event.OnAdminName,
			&event.OnAdminArchivedAt,
			&event.OnAdminLastActiveAt,

			&event.OnUser,
			&event.OnUserName,
			&event.OnUserArchivedAt,
			&event.OnUserLastActiveAt,

			&event.OnBrand,
			&event.OnBrandLabel,
			&event.OnBrandArchivedAt,

			&event.OnIXrule,
			&event.OnIXruleName,
			&event.OnIXruleArchivedAt,

			&event.OnCorporation,
			&event.OnCorporationName,
			&event.OnCorporationArchivedAt,

			&event.OnCorpMapping,
			&event.OnCorpMappingName,
			&event.OnCorpMappingArchivedAt,

			&event.OnLocation,

			&event.OnDataset,
			&event.OnDatasetName,
			&event.OnDatasetArchivedAt,

			&event.OnCustomer,
			&event.OnCustomerName,
			&event.OnCustomerArchivedAt,

			&event.Name,
			&payload,
			&event.AuthMethod,
		); err != nil {
			return events, errors.UnexpectedError(err, "Failed scanning events row")
		}

		if err := json.Unmarshal([]byte(payload), &event.Payload); err != nil {
			return events, errors.UnexpectedError(err, "Failed unmarshalling events payload")
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return events, errors.UnexpectedError(err, "Failed iterating events rows")
	}

	return events, err
}
