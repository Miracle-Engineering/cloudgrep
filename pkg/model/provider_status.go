package model

import "time"

const (
	ProviderStatusFetching = "fetching"
	ProviderStatusFailed   = "failed"
	ProviderStatusSuccess  = "success"
)

type ProviderStatus struct {
	ProviderType   string         `json:"providerType" gorm:"primaryKey"`
	FetchStatus    string         `json:"fetchStatus"`
	ResourceEvents ResourceEvents `json:"resourceEvents" gorm:"-"`
	ErrorMessage   string         `json:"errorMessage"`
	FetchedAt      time.Time      `json:"fetchedAt"`
}

type ProviderStatuses []ProviderStatus

func NewProviderStatus(providerType string, resourceEvents ResourceEvents, err error) ProviderStatus {
	var providerStatus ProviderStatus
	providerStatus.ProviderType = providerType
	providerStatus.FetchedAt = time.Now()
	if err != nil {
		providerStatus.ErrorMessage = err.Error()
		providerStatus.FetchStatus = ProviderStatusFailed
	} else {
		providerStatus.ResourceEvents = resourceEvents
		providerStatus.setProviderStatus()
	}
	return providerStatus
}

func (ps *ProviderStatus) HasError() bool {
	return ps.ErrorMessage != ""
}

func (ps *ProviderStatus) setProviderStatus() {
	if ps.ResourceEvents != nil {
		ps.FetchStatus = ps.ResourceEvents.GetAggregatedStatus()
	}
}

func (providerStatuses ProviderStatuses) GetAggregatedStatus() string {
	fetchStatus := ProviderStatusSuccess
	for _, providerStatus := range providerStatuses {
		if providerStatus.FetchStatus == ProviderStatusFailed {
			fetchStatus = ProviderStatusFailed
			break
		} else if providerStatus.FetchStatus == ProviderStatusFetching {
			fetchStatus = ProviderStatusFetching
		}
	}
	return fetchStatus
}
