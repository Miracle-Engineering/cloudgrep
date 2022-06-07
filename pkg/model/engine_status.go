package model

import (
	"time"
)

const (
	EngineStatusFetching = "fetching"
	EngineStatusSuccess  = "success"
	EngineStatusFailed   = "failed"
)

type EngineStatus struct {
	ResourceType string    `json:"resource_type" gorm:"primaryKey"`
	ErrorMessage string    `json:"error_message"`
	Status       string    `json:"status"`
	FetchedAt    time.Time `json:"fetched_at"`
}

func NewEngineStatus(status string, resource string, err error) EngineStatus {
	var engineStatus EngineStatus
	engineStatus.Status = status
	engineStatus.ResourceType = resource
	engineStatus.FetchedAt = time.Now()
	if err != nil {
		engineStatus.ErrorMessage = err.Error()
	}
	return engineStatus
}
