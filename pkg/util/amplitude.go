package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/run-x/cloudgrep/pkg/version"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/run-x/cloudgrep/pkg/config"

	"github.com/google/uuid"
	"github.com/matishsiao/goInfo"
)

func GetConfigEmail() string {
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
	amplitudeToken = "751db5fc75ff34f08a83381f4d54ead6"
	amplitudeUrl   = "https://api2.amplitude.com/2/httpapi"
	application    = "cloudgrep"
)

const (
	BaseEvent = "BASE_EVENT"
)

var sessionId = SessionId()
var userId = GetConfigEmail()
var validEvents = []string{BaseEvent}

func isValidEvent(eventType string) bool {
	for _, validEvent := range validEvents {
		if validEvent == eventType {
			return true
		}
	}
	return false
}

func GenerateAmplitudeEvent(eventType string, eventProperties map[string]interface{}) (map[string]interface{}, error) {
	if eventProperties == nil {
		eventProperties = make(map[string]interface{})
	}

	eventProperties["application"] = application

	systemInfo, err := goInfo.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}

	id := uuid.New().String()
	event := map[string]interface{}{
		"user_id":          userId,
		"device_id":        systemInfo.Hostname,
		"event_type":       eventType,
		"event_properties": eventProperties,
		"app_version":      version.Version,
		"platform":         systemInfo.Platform,
		"insert_id":        id,
		"session_id":       sessionId,
	}

	return event, nil
}

func SendAmplitudeEvent(ctx context.Context, cfg config.Config, eventType string, eventProperties map[string]interface{}) (int, error) {
	if version.IsDev() {
		return 1, fmt.Errorf("dev application, not sending events to amplitude") //dev application, not sending events to amplitude
	}

	if !isValidEvent(eventType) {
		return 1, fmt.Errorf("invalid event type: %s, not sending events to amplitude\n", eventType) // not sending invalid events
	}

	amplitudeEvent, err := GenerateAmplitudeEvent(eventType, eventProperties)
	if err != nil {
		return 1, fmt.Errorf("failed to generate amplitude event: %w", err) // don't send event to amplitude
	}

	amplitudeBody, err := json.Marshal(map[string]interface{}{"api_key": amplitudeToken, "events": []interface{}{amplitudeEvent}})
	if err != nil {
		return 1, fmt.Errorf("failed to marshal amplitude event: %w", err) // don't send event to amplitude
	}

	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("POST", amplitudeUrl, bytes.NewBuffer(amplitudeBody))
	if err != nil {
		return 1, fmt.Errorf("failed to create amplitude request: %w", err) // don't send event to amplitude
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	response, err := client.Do(req)
	if err != nil {
		return 1, fmt.Errorf("failed amplitude response: %w", err) //failed to get amplitude response
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 1, fmt.Errorf("failed to read amplitude response body: %w", err)
		}
		cfg.Logging.Logger.Sugar().Debug("amplitude response status code: %d", response.StatusCode)
		cfg.Logging.Logger.Sugar().Debug("amplitude response body: %s", string(responseBody))
	}
	return 0, nil
}

func SendEvent(ctx context.Context, cfg config.Config, eventType string, eventProperties map[string]interface{}) {
	go func() {
		_, err := SendAmplitudeEvent(ctx, cfg, eventType, eventProperties)
		if err != nil {
			cfg.Logging.Logger.Sugar().Debug(err)
		}
	}()
}
