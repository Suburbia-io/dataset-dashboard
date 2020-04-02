package database

import (
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/qb"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
)

type Session struct {
	tables.Session
	User User
}

func addRelatedUser(db *DBAL, session *Session) (err error) {
	session.User, err = db.UserGet(session.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBAL) SessionCreate(userID string, lifetime time.Duration) (session Session, err error) {
	session.UserID = userID
	session.Token = crypto.RandAlphaNum(32)
	session.ExpiresAt = time.Now().Add(lifetime)
	err = tables.Sessions.Insert(db, &session.Session)
	if err != nil {
		return session, err
	}
	err = addRelatedUser(db, &session)
	return session, err
}

func (db *DBAL) SessionGet(token string) (session Session, err error) {
	builder := qb.Builder{}
	builder.Write(`SELECT`)
	builder.Write(tables.Sessions.SelectCols())
	builder.Write(`FROM`)
	builder.Write(tables.Sessions.Table())
	builder.Write(`WHERE token = $1`, token)
	builder.Write(`AND expires_at > NOW()`)

	query, queryArgs := builder.MustBuild()
	sessions, err := tables.Sessions.List(db, query, queryArgs...)
	if err != nil || len(sessions) == 0 {
		return session, err
	}

	session.Session = sessions[0]
	err = addRelatedUser(db, &session)
	return session, err
}

func (db *DBAL) SessionDelete(token string) (err error) {
	return tables.Sessions.Delete(db, token)
}
