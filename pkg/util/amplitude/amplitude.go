package amplitude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/denisbrodbeck/machineid"
	"github.com/juandiegopalomino/cloudgrep/pkg/version"

	"github.com/google/uuid"
	"github.com/matishsiao/goInfo"
	"github.com/tcnksm/go-gitconfig"
)

const (
	amplitudeApiKey = "2b0167b9ea1dacf8f0dae96326abd879"
	amplitudeUrl    = "https://api2.amplitude.com/2/httpapi"
)

type EventType int

type client interface {
	SendEvent(logger *zap.Logger, eventType EventType, eventProperties map[string]string)
}

type amplitudeClient struct {
}

const (
	Invalid = "invalid"
	//When adding an event here, also add it to the String() function above
	EventLoad EventType = iota
	EventCloudConnection
)

var (
	defaultClient client = amplitudeClient{}
)

func (s EventType) String() string {
	switch s {
	case EventLoad:
		return "LOAD"
	case EventCloudConnection:
		return "CLOUD_CONNECTION"
	}
	return Invalid
}

func (s EventType) isValid() bool {
	return s.String() != Invalid
}

//newUploadRequest creates an amplitude body for upload, the upload always contain one event in our case.
//https://developers.amplitude.com/docs/http-api-v2#uploadrequestbody
func newUploadRequest(eventType EventType, eventProperties map[string]string) (map[string]interface{}, error) {

	if !eventType.isValid() {
		return nil, fmt.Errorf("invalid event type")
	}

	systemInfo, err := goInfo.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}
	//if we can't have a unique device_id, we should not report to amplitude
	machineId, err := machineid.ID()
	if err != nil {
		return nil, fmt.Errorf("failed to get the platform specific machine id: %w", err)
	}

	event := map[string]interface{}{
		//A device-specific identifier, such as the Identifier for Vendor on iOS.
		//anonymous users sharing the same device_id will be merged, so the device_id has to uniquely identify physical devices,
		//and should not be set to your server's machine name.
		"device_id":   machineId,
		"event_type":  eventType.String(),
		"app_version": version.Version,
		"platform":    "backend",
		//should be: Linux, Mac, Windows...
		"os_name": strings.Replace(systemInfo.GoOS, "darwin", "mac", 1),
		//A unique identifier for the event
		"insert_id": uuid.New().String(),
		//The start time of the session in milliseconds since epoch
		"session_id": time.Now().UnixNano(),
		//"$remote" to use the IP address on the upload request
		"ip":               "$remote",
		"event_properties": eventProperties,
	}

	//if email is found, add it
	email, _ := gitconfig.Email()
	if email != "" {
		event["user_id"] = email
	}

	//wrap the event inton a upload bject
	amplitudeUpload := map[string]interface{}{
		"api_key": amplitudeApiKey,
		"events":  []map[string]interface{}{event},
	}
	return amplitudeUpload, nil
}

func (amplitudeClient) sendEventHTTP(logger *zap.Logger, amplitudeUri string, eventType EventType, eventProperties map[string]string) error {

	amplitudeUpload, err := newUploadRequest(eventType, eventProperties)
	if err != nil {
		return fmt.Errorf("failed to generate amplitude upload body: %w", err)
	}

	if version.IsDev() {
		//dev application, not sending events to amplitude
		return fmt.Errorf("dev application, not sending events to amplitude")
	}

	amplitudeBody, err := json.Marshal(amplitudeUpload)
	if err != nil {
		return fmt.Errorf("failed to marshal amplitude event: %w", err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("POST", amplitudeUri, bytes.NewBuffer(amplitudeBody))
	if err != nil {
		return fmt.Errorf("failed to create amplitude request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get amplitude response: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("failed to read amplitude response body: %w", err)
		}
		logger.Sugar().Debug("amplitude response status code: %d", response.StatusCode)
		logger.Sugar().Debug("amplitude response body: %s", string(responseBody))
	}
	return nil
}

func (c amplitudeClient) SendEvent(logger *zap.Logger, eventType EventType, eventProperties map[string]string) {
	go func() {
		err := c.sendEventHTTP(logger, amplitudeUrl, eventType, eventProperties)
		if err != nil {
			logger.Sugar().Debug(err)
		}
	}()
}

func SendEvent(logger *zap.Logger, eventType EventType, eventProperties map[string]string) {
	defaultClient.SendEvent(logger, eventType, eventProperties)
}
