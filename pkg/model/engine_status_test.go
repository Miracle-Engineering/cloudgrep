package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeEngineStatus(t *testing.T) {
	mockResourceName := "engine"
	statusFetching := NewEngineStatus(EngineStatusSuccess, mockResourceName, nil)
	assert.Equal(t, "", statusFetching.ErrorMessage)
	assert.Equal(t, EngineStatusSuccess, statusFetching.Status)
	assert.Equal(t, mockResourceName, statusFetching.ResourceType)

	statusFailedError := errors.New("failed")
	statusFailed := NewEngineStatus(EngineStatusFailed, mockResourceName, statusFailedError)
	assert.Equal(t, statusFailedError.Error(), statusFailed.ErrorMessage)
	assert.Equal(t, EngineStatusFailed, statusFailed.Status)
	assert.Equal(t, mockResourceName, statusFetching.ResourceType)
}
