package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeEngineStatus(t *testing.T) {
	statusFetching := MakeEngineStatusFetching()
	assert.Equal(t, "", statusFetching.ErrorMessage)
	assert.Equal(t, EngineStatusFetching, statusFetching.Status)

	statusSuccess := MakeEngineStatusSuccess()
	assert.Equal(t, "", statusSuccess.ErrorMessage)
	assert.Equal(t, EngineStatusSuccess, statusSuccess.Status)

	statusFailedError := errors.New("failed")
	statusFailed := MakeEngineStatusFailed(statusFailedError)
	assert.Equal(t, statusFailedError.Error(), statusFailed.ErrorMessage)
	assert.Equal(t, EngineStatusFailed, statusFailed.Status)
}
