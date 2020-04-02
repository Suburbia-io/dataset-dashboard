package application

import (
	"bytes"
	"encoding/json"
	errors2 "errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/testdb"
)

var tdb *testdb.Manager
var testConfig Config

func TestMain(m *testing.M) {
	_, err := toml.DecodeFile(os.Getenv("SUBURBIA_DASHBOARD_CONFIG"), &testConfig)
	if err != nil {
		panic(err)
	}

	tdb, err = testdb.NewManager(
		"app_testing", 3,
		testConfig.DB.DBHost,
		testConfig.DB.DBUser,
		testConfig.DB.DBPassword,
		testConfig.DB.DBPort,
		testConfig.DB.DBSSLMode,
	)
	if err != nil {
		panic(err)
	}

	log.SetOutput(ioutil.Discard)

	var status int
	defer func() {
		recover()
		tdb.TearDown()
		os.Exit(status)
	}()
	status = m.Run()
}

func NewTestServer() (app *App, srv *httptest.Server, close func()) {
	app, closeapp := NewTestApp()
	srv = httptest.NewServer(app.Routes())

	return app, srv, func() {
		defer closeapp()
		defer srv.Close()
	}
}

func NewTestApp() (app *App, close func()) {
	conn, close, err := tdb.NewConn()
	if err != nil {
		panic(err)
	}

	app, err = Mount(testConfig, []Service{
		NewTestDBService(conn),
	})
	if err != nil {
		panic(err)
	}

	return app, close
}

// -----------------------------------------------------------------------------
// TestClient
// -----------------------------------------------------------------------------
func CallWeb(url string) (resp *http.Response, err error) {
	return callWeb("", "", url)
}

func CallAdminWeb(tokenVal, url string) (resp *http.Response, err error) {
	return callWeb("suburbia-dashboard-admin", tokenVal, url)
}

func callWeb(tokenName, tokenVal, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resp, err
	}

	if tokenVal != "" {
		req.AddCookie(&http.Cookie{Name: tokenName, Value: tokenVal})
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func CallApi(fn string, args interface{}, reply interface{}) (resp *http.Response, err error) {
	return callAPI("", "", fn, args, reply)
}

func CallAdminApi(tokenVal, fn string, args interface{}, reply interface{}) (resp *http.Response, err error) {
	return callAPI("suburbia-dashboard-admin", tokenVal, fn, args, reply)
}

func callAPI(tokenName, tokenVal, fn string, args interface{}, reply interface{}) (resp *http.Response, err error) {
	b, err := json.Marshal(&args)
	if err != nil {
		return resp, err
	}

	req, err := http.NewRequest("POST", fn, bytes.NewBuffer(b))
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	if tokenVal != "" {
		req.AddCookie(&http.Cookie{Name: tokenName, Value: tokenVal})
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 200 {
		return resp, errors2.New("HttpBadStatusCode")
	}

	if err := unmarshalApiResponse(resp, reply); err != nil {
		return resp, err
	}

	return resp, nil
}

func unmarshalApiResponse(resp *http.Response, reply interface{}) (err error) {
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var wrapper struct {
		OK   bool            `json:"ok"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(resBody, &wrapper); err != nil {
		return err
	}

	if !wrapper.OK {
		var errMsg string
		if err := json.Unmarshal(wrapper.Data, &errMsg); err != nil {
			return err
		} else {
			return errors.NewErr(errMsg)
		}
	}

	if err := json.Unmarshal(wrapper.Data, &reply); err != nil {
		return err
	}

	return nil
}
