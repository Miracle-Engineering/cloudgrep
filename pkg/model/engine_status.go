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
	ResourceName string    `json:"resource_name" gorm:"primaryKey"`
	ErrorMessage string    `json:"error_message"`
	Status       string    `json:"status"`
	FetchedAt    time.Time `json:"fetched_at"`
}

func NewEngineStatus(status string, resource string, err error) EngineStatus {
	var engineStatus EngineStatus
	engineStatus.Status = status
	engineStatus.ResourceName = resource
	engineStatus.FetchedAt = time.Now()
	if err != nil {
		engineStatus.ErrorMessage = err.Error()
	}
	return engineStatus
}
