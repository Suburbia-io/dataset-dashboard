package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Alerter interface {
	AlertError(err error, format string, args ...interface{})
}

var alerter Alerter = NewLogAlerter()

func SetGlobalAlerter(a Alerter) {
	alerter = a
}

func AlertError(err error, format string, args ...interface{}) {
	alerter.AlertError(err, format, args...)
}

// -----------------------------------------------------------------------------
type LogAlerter struct{}

func NewLogAlerter() *LogAlerter {
	return &LogAlerter{}
}

func (la *LogAlerter) AlertError(err error, s string, args ...interface{}) {
	log.Printf("ALERT: [%v] %s", err, fmt.Sprintf(s, args...))
}

// -----------------------------------------------------------------------------

type SlackAlerter struct {
	prefix     string
	webhookURL string
}

func NewSlackAlerter(prefix, webhookURL string) *SlackAlerter {
	hostname, err := os.Hostname()
	if err == nil {
		prefix = prefix + "@" + hostname + " "
	}
	return &SlackAlerter{prefix, webhookURL}
}

func (a *SlackAlerter) AlertError(err error, s string, args ...interface{}) {
	s = fmt.Sprintf(s, args...)
	if err != nil {
		s = "[" + err.Error() + "] " + s
	}
	s = a.prefix + s
	log.Println("ALERT: " + s)

	type Msg struct {
		Text string `json:"text"`
	}

	buf, err := json.Marshal(Msg{s})
	if err != nil {
		log.Printf("ALERT: Failed to marshal alert: %v", err)
		return
	}

	resp, err := http.Post(
		a.webhookURL, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("ALERT: Failed to send alert: %v", err)
		return
	}

	resp.Body.Close()
}
