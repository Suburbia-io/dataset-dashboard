package application

import (
	"net/http"
	"net/url"

	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/shared"
	"golang.org/x/crypto/ssh"
)

func (app *App) CustomerUserDashboard(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	ctx := struct {
		Datasets          []tables.Dataset
		SFTPHostname      string
		PubKeyFingerprint string
		Session           shared.SharedSession
	}{
		Session: s,
	}

	u, err := url.Parse(app.Config.ServerHostname)
	if err != nil {
		return err
	}
	ctx.SFTPHostname = u.Hostname()

	ctx.Datasets, err = app.DBAL.DatasetListActiveForCustomer(s.User.CustomerID)
	if err != nil {
		return err
	}

	if s.User.SFTPPubKey != "" {
		pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(s.User.SFTPPubKey))
		if err != nil {
			return err
		}
		ctx.PubKeyFingerprint = ssh.FingerprintLegacyMD5(pk)
	}

	return app.Views.ExecuteSST(w, "customer-user-dashboard.gohtml", ctx)
}
