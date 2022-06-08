package model

import "time"

const (
	EngineStatusSuccess          = "success"
	EngineStatusFailed           = "failed"
	EngineStatusFetching         = "fetching"
	EngineStatusFailedAtCreation = "failedAtCreation"
)

type EngineStatus struct {
	FetchStatus    string         `json:"fetchStatus"`
	FetchedAt      time.Time      `json:"fetchedAt"`
	ResourceEvents ResourceEvents `json:"resourceEvents"`
}

func NewEngineStatus(resourceEvents ResourceEvents, failedAtCreation bool) EngineStatus {
	if failedAtCreation {
		return EngineStatus{
			FetchStatus: EngineStatusFailedAtCreation,
		}
	}
	return EngineStatus{
		FetchStatus:    resourceEvents.GetStatus(),
		FetchedAt:      time.Now(),
		ResourceEvents: resourceEvents,
	}
}
