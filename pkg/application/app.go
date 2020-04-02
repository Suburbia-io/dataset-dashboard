package application

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Suburbia-io/dashboard/pkg/helpers/slice"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/mailer"
	"github.com/Suburbia-io/dashboard/pkg/views"
)

type App struct {
	Config Config
	DBAL   *database.DBAL
	Views  *views.Views
	Mailer mailer.Mailer
	Router *mux.Router
}

type apiResponse struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data"`
}

type Service func(app *App) (err error)

type Config struct {
	Env       string // What possible values jerk?
	StaticDir string
	DataDir   string

	DB     database.Config
	Views  views.Config
	Mailer mailer.Config

	AdminAuthTokenLifetimeSec time.Duration
	SessionLifetimeSec        time.Duration
	SessionCookieDomain       string
	SessionCookieSecure       bool

	AlertURL string

	Hostnames []string
	HTTPPort  string

	// Used for email links and showing users the sftp url.
	ServerHostname string

	SftpDir        string
	SftpListenAddr string
	SftpPrivateKey string

	SuperAdminEmail    string
	SuperAdminPassword string

	GoogleMapsApiKey    string
	GeonamesApiUsername string
	GeonamesApiToken    string
	CSRFSecret          string
}

var decoder = schema.NewDecoder()
var encoder = schema.NewEncoder()

func Mount(config Config, services []Service) (app *App, err error) {
	app = &App{
		Config: config,
	}

	for _, s := range services {
		if err = s(app); err != nil {
			return app, err
		}
	}

	// check if the csrf secret is set to a reasonably long string
	if len(config.CSRFSecret) != 32 {
		panic("no 32 character csrf secret has been set in the config file. please set CSRFSecret.")
	}

	app.Router = mux.NewRouter()

	decoder.IgnoreUnknownKeys(true)

	app.registerTmplFuncs()

	return app, err
}

func (app *App) decodeRequest(r *http.Request, args interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		return errors.HttpBadRequestArgs
	}
	return nil
}

func (app *App) respondApi(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	if err != nil {
		if appErr, ok := err.(errors.Error); !ok {
			app.serverErr(w, r, err)
			return
		} else {
			if errors.Unwrap(appErr) != nil {
				data = fmt.Sprintf("%s", errors.Unwrap(appErr))
			} else {
				data = appErr.Code()
			}
		}
	}

	if resp, err := json.Marshal(apiResponse{err == nil, data}); err != nil {
		app.serverErr(w, r, err)
	} else {
		w.Write(resp)
	}
}

func (app *App) serverErr(w http.ResponseWriter, r *http.Request, err error) {
	log.Output(2, fmt.Sprintf("ERROR: %s\n%s", err.Error(), debug.Stack()))

	if r.Header.Get("Content-Type") == "application/json" {
		w.Write([]byte(`{"ok": false, "data": "Unexpected"}`))
	} else {
		w.WriteHeader(500)
		w.Write([]byte("Unexpected"))
	}
}

func (app *App) registerTmplFuncs() {
	// reverse url lookup based on route names
	app.Views.RegisterFunc("ReverseURL", func(name string, pairs ...string) string {
		route := app.Router.Get(name)
		if route == nil {
			return ""
		}
		revUrl, err := route.URL(pairs...)
		if err != nil {
			return ""
		}
		return revUrl.String()
	})

	app.Views.RegisterFunc("ReverseURLFromAuditLogPayload", func(routeName string, jsonString string, keys ...string) string {
		target := make(map[string]string)
		err := json.Unmarshal([]byte(jsonString), &target)
		if err != nil {
			return ""
		}
		var pairs []string

		for key, value := range target {
			uppercasedKey := strings.Title(key)
			if slice.ContainsString(keys, uppercasedKey) {
				pairs = append(pairs, uppercasedKey)
				pairs = append(pairs, value)
			}
		}

		route := app.Router.Get(routeName)
		if route == nil {
			return ""
		}
		revUrl, err := route.URL(pairs...)
		if err != nil {
			return ""
		}
		return revUrl.String()
	})

	app.Views.RegisterFunc("Split", strings.Split)

	app.Views.RegisterFunc("ReplaceAll", strings.ReplaceAll)

	app.Views.RegisterFunc("DerefBool", func(b *bool) bool {
		if b == nil {
			return false
		}
		return *b
	})

	app.Views.RegisterFunc("DerefString", func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	})

	// Returns a query param string with the current query params adding or overriding the supplied key and value. Empty
	// and nil query param values are not included in the result.
	app.Views.RegisterFunc("GetQueryParams",
		func(listArgs interface{}, overrideKey, overrideValue string) template.URL {
			values := make(map[string][]string)
			err := encoder.Encode(listArgs, values)
			if err != nil {
				return ""
			}

			var removeKeys []string
			for key, _ := range values {
				// Deal with the supplied overrides.
				if key == overrideKey {
					if overrideValue == "" {
						removeKeys = append(removeKeys, key)
					} else {
						values[key] = []string{overrideValue}
					}
				}

				// Remove nil and empty query params
				valStr := fmt.Sprintf("%s", values[key])
				if valStr == "[null]" || valStr == "[]" {
					removeKeys = append(removeKeys, key)
				}
			}
			for _, value := range removeKeys {
				delete(values, value)
			}

			// This safe because the listArgs have already been decoded in the function handlers.
			return template.URL(url.Values(values).Encode())
		})

}
