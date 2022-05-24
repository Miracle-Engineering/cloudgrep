package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	testCases := []struct {
		name           string
		expectedOutput string
		args           []string
	}{
		{"AllGood", "version.BuildInfo{Version:\"dev\", GitCommit:\"\", BuildTime:\"\", GoVersion:\"testing\"}\n", []string{"version"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd := NewRootCmd(buf)
			rootCmd.SetArgs(tc.args)
			err := rootCmd.Execute()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, buf.String())

		})
	}
}
