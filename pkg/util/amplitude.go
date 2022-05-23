package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/matishsiao/goInfo"
	"github.com/run-x/cloudgrep/pkg/app"
	"github.com/tcnksm/go-gitconfig"
)

func GetGitConfigEmail() string {
	email, err := gitconfig.Email()
	if err != nil {
		return "no_user"
	}
	return email
}

func SessionId() int64 {
	return time.Now().UnixNano()
}

const (
	AMPLITUDE_TOKEN = "751db5fc75ff34f08a83381f4d54ead6"
	BASE_EVENT      = "BASE_EVENT"
	AMPLITUDE_URL   = "https://api2.amplitude.com/2/httpapi"
)

var SESSION_ID = SessionId()
var USER_ID = GetGitConfigEmail()
var VALID_EVENTS = []string{BASE_EVENT}

func isValidEvent(eventType string) bool {
	for _, validEvent := range VALID_EVENTS {
		if validEvent == eventType {
			return true
		}
	}
	return false
}

func sendAmplitudeEvent(eventType string, eventProperties map[string]interface{}, userProperties map[string]interface{}) {
	if app.IsDev() {
		fmt.Println("dev application, not sending events to amplitude")
	}

	if !isValidEvent(eventType) {
		fmt.Printf("invalid event type: %s, not sending events to amplitude\n", eventType)
	}

	if eventProperties == nil {
		eventProperties = make(map[string]interface{})
	}

	if userProperties == nil {
		userProperties = make(map[string]interface{})
	}

	systemInfo, _ := goInfo.GetInfo()

	id := uuid.New().String()
	event := map[string]interface{}{
		"user_id":          USER_ID,
		"device_id":        systemInfo.Hostname,
		"event_type":       eventType,
		"event_properties": eventProperties,
		"user_properties":  userProperties,
		"app_version":      app.Version,
		"platform":         systemInfo.Platform,
		"insert_id":        id,
		"session_id":       SESSION_ID,
	}

	amplitudeBody, _ := json.Marshal(map[string]interface{}{"api_key": AMPLITUDE_TOKEN, "events": []interface{}{event}})

	client := &http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("POST", AMPLITUDE_URL, bytes.NewBuffer(amplitudeBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	} else {
		if response.StatusCode != 200 {
			fmt.Printf("amplitude response: %s\n", response.Status)
			responseBody, _ := ioutil.ReadAll(response.Body)
			fmt.Printf("amplitude response: %s\n", responseBody)
		}
	}
}

func SendEvent(eventType string, eventProperties map[string]interface{}, userProperties map[string]interface{}) {
	go sendAmplitudeEvent(eventType, eventProperties, userProperties)
}
