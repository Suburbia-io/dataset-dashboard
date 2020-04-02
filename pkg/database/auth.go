package database

import (
	"time"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
)

func GenerateAuthToken() string {
	return crypto.RandAlphaNumForHumans(4) +
		"-" + crypto.RandAlphaNumForHumans(4) +
		"-" + crypto.RandAlphaNumForHumans(4) +
		"-" + crypto.RandAlphaNumForHumans(4)
}

func (db *DBAL) AuthPurgeExpiredTokens() (err error) {
	_, err = db.Exec(`DELETE FROM auth_tokens WHERE expires_at < NOW()`)
	if err != nil {
		return errors.UnexpectedError(err, "Failed purging expired auth-token")
	}
	return nil
}

func authTokensCount(db *DBAL) (count int, err error) {
	err = db.QueryRow(`SELECT COUNT(*) FROM auth_tokens`).Scan(&count)
	if err != nil {
		return count, errors.UnexpectedError(err, "Failed counting admin-tokens")
	}
	return count, nil
}

// -----------------------------------------------------------------------------
// Auth tokens for Admins
// -----------------------------------------------------------------------------

func (db *DBAL) AuthCreateAdminToken(userID string, lifetime time.Duration) (emailToken string, browserToken string, err error) {
	defer db.AuthPurgeExpiredTokens()

	emailToken = GenerateAuthToken()
	browserToken = crypto.RandAlphaNum(32)
	expiresAt := time.Now().Add(lifetime)

	stmt := `INSERT INTO auth_tokens(
        token,
        browser_token,
        expires_at,
        user_id
    ) VALUES ($1,$2,$3,(
        SELECT user_id FROM users
        WHERE user_id=$4 AND archived_at IS NULL
    ))`
	_, err = db.Exec(stmt, emailToken, browserToken, expiresAt, userID)

	if err == nil {
		return emailToken, browserToken, nil
	}

	if dbIsNullConstraintErr(err, "user_id") {
		return "", "", errors.DBNotFound
	}

	return "", "", errors.UnexpectedError(err, "Failed creating auth token")
}

func (db *DBAL) AuthAuthenticateAdminToken(token, browserToken, suburbiaInternalCustomerUUID string) (userID string, err error) {
	stmt := `DELETE FROM auth_tokens WHERE token=(
        SELECT token
        FROM auth_tokens t
        JOIN users u ON u.user_id=t.user_id
        WHERE u.archived_at IS NULL
            AND t.token=$1
            AND t.browser_token=$2
            AND u.customer_id=$3
            AND t.expires_at > NOW()
    ) RETURNING user_id`
	err = db.QueryRow(stmt, token, browserToken, suburbiaInternalCustomerUUID).Scan(&userID)
	if err != nil {
		return "", errors.AuthTokenNotFound
	}

	return userID, nil
}

func (db *DBAL) AuthPeekAdminTokens(userID string) (tokens []string, err error) {
	if _, err := db.UserGet(userID); err != nil {
		return tokens, err
	}

	stmt := `SELECT token
        FROM auth_tokens t
        JOIN users u ON u.user_id=t.user_id
        WHERE t.user_id=$1
          AND u.archived_at IS NULL
          AND u.is_role_admin IS TRUE
          AND t.expires_at > NOW()`

	rows, err := db.Query(stmt, userID)
	if err != nil {
		return tokens, errors.UnexpectedError(err, "Failed selecting auth-tokens by userID")
	}
	defer rows.Close()

	for rows.Next() {
		token := ""
		if err := rows.Scan(&token); err != nil {
			return tokens, errors.UnexpectedError(err, "Failed scanning auth-tokens row")
		}

		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return tokens, errors.UnexpectedError(err, "Failed iterating auth-tokens rows")
	}

	return tokens, err
}
