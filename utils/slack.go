package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	url     = "https://hooks.slack.com/services"
	errChan = "/T01DBF9B7U1/B024F7JSMN3/N0ulEVHAvfMgtKK381hHJ2x2"
)

type body struct {
	Text string `json:"text"`
}

func NotifySlack(msg string, env string) {
	client := &http.Client{}
	text := "[" + env + "] - Error Watcher - " + msg
	data := &body{Text: text}

	b, _ := json.Marshal(&data)

	if _, err := client.Post(url+errChan, jsonContentType, bytes.NewReader(b)); err != nil {
		ErrorReport.LogAndPrintError(err)
	}
}
