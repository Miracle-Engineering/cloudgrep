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
	Provider     string    `json:"-"`
	FetchStatus  string    `json:"fetchStatus"`
	ErrorMessage string    `json:"errorMessage"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ResourceEvents []ResourceEvent

func NewResourceEvent(resourceType string, provider string, isFetching bool, err error) ResourceEvent {
	var resourceEvent ResourceEvent
	resourceEvent.ResourceType = resourceType
	resourceEvent.Provider = provider
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

func (resourceEvents ResourceEvents) GetAggregatedStatus() string {
	fetchStatus := ProviderStatusSuccess
	for _, resourceEvent := range resourceEvents {
		if resourceEvent.FetchStatus == ResourceEventStatusFailed {
			fetchStatus = ProviderStatusFailed
			break
		} else if resourceEvent.FetchStatus == ResourceEventStatusFetching {
			fetchStatus = ProviderStatusFetching
		}
	}
	return fetchStatus
}
