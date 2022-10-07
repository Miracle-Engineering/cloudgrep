package amplitude

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/juandiegopalomino/cloudgrep/pkg/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendAmplitudeEvent(t *testing.T) {

	client := MemoryClient{}
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		eventType       EventType
		eventProperties map[string]string
		wantEventType   string
	}{
		{
			eventType:     EventLoad,
			wantEventType: "LOAD",
		},
		{
			eventType:       EventCloudConnection,
			eventProperties: map[string]string{"CLOUD_ID": "0987654321"},
			wantEventType:   "CLOUD_CONNECTION",
		},
	}

	for _, tc := range tests {

		client.SendEvent(logger, tc.eventType, tc.eventProperties)

		event, err := client.LastEvent()
		require.NoError(t, err)

		require.Equal(t, tc.wantEventType, event["event_type"])
		require.Equal(t, version.Version, event["app_version"])
		require.Equal(t, "$remote", event["ip"])
		require.Equal(t, "backend", event["platform"])
		require.Equal(t, tc.eventProperties, event["event_properties"])

		//special care for user_id/email which might not always be set
		if userId, found := event["user_id"]; found {
			assert.Truef(t, strings.Contains(fmt.Sprint(userId), "@"), "%v is not a valid email", userId)
		}

		//insert_id - uuid is different
		_, err = uuid.Parse(fmt.Sprint(event["insert_id"]))
		require.NoError(t, err)

		//device id is the machine id
		machineId, err := machineid.ID()
		require.NoError(t, err)
		require.Equal(t, machineId, fmt.Sprint(event["device_id"]))

		//session id is a timestamp
		sessionId := event["session_id"].(int64)
		require.Less(t, sessionId, time.Now().UnixNano())

		require.Contains(t, []string{"linux", "mac", "windows"}, event["os_name"])
	}
}

func TestSendAmplitudeEventErrors(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	client := amplitudeClient{}

	t.Run("TestSendAmplitudeEventDevVersion", func(t *testing.T) {
		err := client.sendEventHTTP(logger, amplitudeUrl, EventLoad, nil)
		assert.ErrorContains(t, err, "dev application, not sending events to amplitude")
	})

	version.Version = "test"
	revertVersion := func() {
		version.Version = "dev"
	}
	defer revertVersion()

	invalidEvent := EventType(99)
	t.Run("TestSendAmplitudeEventInvalidEvent", func(t *testing.T) {
		err := client.sendEventHTTP(logger, amplitudeUrl, invalidEvent, nil)
		require.ErrorContains(t, err, "failed to generate amplitude upload body: invalid event type")
	})

	t.Run("TestSendAmplitudeEventInvalidUri", func(t *testing.T) {
		err := client.sendEventHTTP(logger, "localhost:8080/", EventLoad, nil)
		require.ErrorContains(t, err, "failed to get amplitude response: Post \"localhost:8080/\": unsupported protocol scheme \"localhost\"")
	})

}
