package database

import (
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validators"
)

type SftpSession struct {
	SftpSessionID string     `json:"sftpSessionID" db:"sftp_session_id"`
	SessionToken  string     `json:"sessionToken" db:"session_token"`
	User          User       `json:"user"`
	UserID        string     `json:"userID" db:"user_entity"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	ArchivedAt    *time.Time `json:"archivedAt" db:"archived_at"`
}

func (db *DBAL) SftpSessionCreate(sessionToken string, user User) (session SftpSession, err error) {
	session = SftpSession{
		SftpSessionID: crypto.NewUUID(),
		SessionToken:  sessionToken,
		UserID:        user.UserID,
		ArchivedAt:    nil,
		UpdatedAt:     time.Now().Round(time.Microsecond),
		CreatedAt:     time.Now().Round(time.Microsecond),
	}

	if err = validators.UUIDDeleteMe(user.UserID); err != nil {
		return session, err
	}

	return session, db.insert("sftp_sessions", session)
}

func (db *DBAL) SftpSessionGetByToken(sessionToken string) (sftpSession SftpSession, err error) {
	return sftpSessionGetByToken(db, sessionToken)
}

func sftpSessionGetByToken(dbTX *DBAL, sessionToken string) (sftpSession SftpSession, err error) {
	err = dbTX.getBy(&sftpSession, "sftp_sessions", map[string]interface{}{"session_token": sessionToken})
	user, err := dbTX.UserGet(sftpSession.UserID)
	sftpSession.User = user
	return sftpSession, err
}
