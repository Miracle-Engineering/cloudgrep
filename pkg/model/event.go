package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id           int64     `json:"-" gorm:"primaryKey;autoIncrement"`
	RunId        string    `json:"runId"`
	Type         string    `json:"eventType"`
	Status       string    `json:"status"`
	ProviderName string    `json:"providerName,omitempty"`
	ResourceType string    `json:"resourceType,omitempty"`
	Error        string    `json:"error"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	ChildEvents  Events    `json:"childEvents" gorm:"-"`
}

type Events []Event

func NewEngineEventStart() Event {
	return Event{
		RunId:  uuid.New().String(),
		Type:   EventTypeEngine,
		Status: EventStatusFetching,
	}
}

func NewEngineEventEnd(err error) Event {
	status := EventStatusSuccess
	errr := ""
	if err != nil {
		status = EventStatusFailed
		errr = err.Error()
	}
	return Event{
		Type:   EventTypeEngine,
		Status: status,
		Error:  errr,
	}
}

func NewProviderEventStart(providerName string) Event {
	return Event{
		Type:         EventTypeProvider,
		Status:       EventStatusFetching,
		ProviderName: providerName,
	}
}

func NewProviderEventEnd(providerName string, err error) Event {
	status := EventStatusSuccess
	errr := ""
	if err != nil {
		status = EventStatusFailed
		errr = err.Error()
	}
	return Event{
		Type:         EventTypeProvider,
		Status:       status,
		ProviderName: providerName,
		Error:        errr,
	}
}

func NewResourceEventStart(providerName string, resourceType string) Event {
	return Event{
		Type:         EventTypeResource,
		Status:       EventStatusFetching,
		ProviderName: providerName,
		ResourceType: resourceType,
	}
}

func NewResourceEventEnd(providerName string, resourceType string, err error) Event {
	status := EventStatusSuccess
	errr := ""
	if err != nil {
		status = EventStatusFailed
		errr = err.Error()
	}
	return Event{
		Type:         EventTypeResource,
		Status:       status,
		ProviderName: providerName,
		ResourceType: resourceType,
		Error:        errr,
	}
}
