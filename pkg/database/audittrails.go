package database

import (
	"encoding/json"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"
)

type AuditTrail = struct {
	tables.AuditTrail
	User *User `json:"user"`
}

type AuditTrailListArgs struct {
	Search string `json:"search" schema:"search"`
	Limit  int    `json:"limit" schema:"limit"`
	Offset int    `json:"offset" schema:"offset"`
}

func (db *DBAL) AuditTrailList(args AuditTrailListArgs) (auditTrails []AuditTrail, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.AuditTrails.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.AuditTrails.Table())
	builder.Write(`WHERE TRUE`)

	if args.Search != "" {
		search := "%" + args.Search + "%"
		builder.Write(`AND (type ILIKE $1 OR related_table ILIKE $1)`, search)
	}

	builder.Write(`ORDER BY created_at DESC`)
	builder.Write(`LIMIT $1 OFFSET $2`, args.Limit, args.Offset)

	query, queryArgs := builder.MustBuild()
	auditTrailList, err := tables.AuditTrails.List(db, query, queryArgs...)
	if err != nil {
		return auditTrails, err
	}

	for _, auditTrailItem := range auditTrailList {
		var extendedAuditTrail = AuditTrail{
			auditTrailItem,
			nil,
		}

		if auditTrailItem.ByUser != nil {
			user, err := db.UserGet(*auditTrailItem.ByUser)
			if err == nil {
				extendedAuditTrail.User = &user
			}
		}

		auditTrails = append(auditTrails, extendedAuditTrail)
	}

	return auditTrails, nil
}

func (db *DBAL) AuditTrailByUserInsertAsync(session Session, table string, relatedID string, eventType string, payload interface{}) {
	go func() {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return
		}

		_ = tables.AuditTrails.Insert(db, &tables.AuditTrail{
			ByUser:       &session.UserID,
			RelatedTable: table,
			RelatedID:    relatedID,
			Type:         eventType,
			Payload:      string(jsonPayload),
			CreatedAt:    time.Now(),
		})
	}()
}

// TODO: make sftpsession & web session use the same struct & table so we can remove this duplicate function
func (db *DBAL) AuditTrailBySftpUserInsertAsync(session SftpSession, table string, relatedID string, eventType string, payload interface{}) {
	go func() {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return
		}

		_ = tables.AuditTrails.Insert(db, &tables.AuditTrail{
			ByUser:       &session.UserID,
			RelatedTable: table,
			RelatedID:    relatedID,
			Type:         eventType,
			Payload:      string(jsonPayload),
			CreatedAt:    time.Now(),
		})
	}()
}

func (db *DBAL) AuditTrailBySystemInsertAsync(table string, relatedID string, eventType string, payload interface{}) {
	go func() {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return
		}

		_ = tables.AuditTrails.Insert(db, &tables.AuditTrail{
			ByUser:       nil,
			RelatedTable: table,
			RelatedID:    relatedID,
			Type:         eventType,
			Payload:      string(jsonPayload),
			CreatedAt:    time.Now(),
		})
	}()
}
