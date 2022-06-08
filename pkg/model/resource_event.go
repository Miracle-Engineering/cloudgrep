package model

import (
	"time"
)

const (
	ResourceEventStatusFetching = "fetching"
	ResourceEventStatusSuccess  = "success"
	ResourceEventStatusFailed   = "failed"
)

type ResourceEvent struct {
	Id           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ResourceType string    `json:"resourceType"`
	FetchStatus  string    `json:"fetchStatus"`
	ErrorMessage string    `json:"errorMessage"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ResourceEvents []ResourceEvent

func NewResourceEvent(resourceType string, isFetching bool, err error) ResourceEvent {
	var resourceEvent ResourceEvent
	resourceEvent.ResourceType = resourceType
	if isFetching {
		resourceEvent.FetchStatus = ResourceEventStatusFetching
	} else {
		if err != nil {
			resourceEvent.FetchStatus = ResourceEventStatusFailed
			resourceEvent.ErrorMessage = err.Error()
		} else {
			resourceEvent.FetchStatus = ResourceEventStatusSuccess
		}
	}
	return resourceEvent
}

func (resourceEvents ResourceEvents) GetStatus() string {
	fetchStatus := EngineStatusSuccess
	for _, resourceEvent := range resourceEvents {
		if resourceEvent.FetchStatus == ResourceEventStatusFailed {
			fetchStatus = EngineStatusFailed
			break
		} else if resourceEvent.FetchStatus == ResourceEventStatusFetching {
			fetchStatus = EngineStatusFetching
		}
	}
	return fetchStatus
}
