package model

import "time"

const (
	EngineStatusSuccess          = "success"
	EngineStatusFailed           = "failed"
	EngineStatusFetching         = "fetching"
	EngineStatusFailedAtCreation = "failedAtCreation"
)

type EngineStatus struct {
	FetchStatus      string           `json:"fetchStatus"`
	FetchedAt        time.Time        `json:"fetchedAt"`
	ProviderStatuses ProviderStatuses `json:"providerStatuses"`
}

func newEngineStatus(providerStatuses ProviderStatuses) EngineStatus {
	var engineStatus EngineStatus
	engineStatus.ProviderStatuses = providerStatuses
	engineStatus.setEngineStatus()
	engineStatus.FetchedAt = time.Now()
	return engineStatus
}

func (es *EngineStatus) setEngineStatus() {
	if es.ProviderStatuses != nil {
		es.FetchStatus = es.ProviderStatuses.GetAggregatedStatus()
	}
}

func GetEngineStatus(providerStatuses ProviderStatuses, resourceEvents ResourceEvents) EngineStatus {
	mapProviderTypeResourceEvents := make(map[string]ResourceEvents)
	for _, resourceEvent := range resourceEvents {
		mapProviderTypeResourceEvents[resourceEvent.Provider] = append(mapProviderTypeResourceEvents[resourceEvent.Provider], resourceEvent)
	}

	for providerType, resourceEvents := range mapProviderTypeResourceEvents {
		providerStatuses = append(providerStatuses, NewProviderStatus(providerType, resourceEvents, nil))
	}

	return newEngineStatus(providerStatuses)

}
