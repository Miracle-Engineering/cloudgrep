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
	ResourceType string    `json:"resourceType" gorm:"primaryKey"`
	ErrorMessage string    `json:"errorMessage"`
	Status       string    `json:"status"`
	FetchedAt    time.Time `json:"fetchedAt"`
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
