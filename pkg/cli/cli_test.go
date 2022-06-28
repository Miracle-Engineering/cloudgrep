package cli

import (
	"context"
	"testing"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/util/amplitude"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestAmplitudeSent(t *testing.T) {

	ctx := context.Background()
	cfg := config.Config{}
	logger := zaptest.NewLogger(t)
	amplitudeClient := amplitude.UseMemoryClient()

	err := Run(ctx, cfg, logger)
	require.ErrorContains(t, err, "unknown datastore type")

	//the cli sends one amplitude event per run (even if error)
	require.Equal(t, 1, amplitudeClient.Size())
	event, err := amplitudeClient.LastEvent()
	require.NoError(t, err)
	require.Equal(t, "LOAD", event["event_type"])
}
