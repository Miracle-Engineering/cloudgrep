package model

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	EventStatusFetching string = "fetching"
	EventStatusFailed   string = "failed"
	EventStatusSuccess  string = "success"
	EventStatusLoaded   string = "loaded"
)

const (
	EventTypeEngine   string = "engine"
	EventTypeProvider string = "provider"
	EventTypeResource string = "resource"
)

type Event struct {
	Id           int64     `json:"-" gorm:"primaryKey;autoIncrement"`
	RunId        string    `json:"runId"`
	Type         string    `json:"eventType"`
	Status       string    `json:"status"`
	ProviderName string    `json:"providerName"`
	ResourceType string    `json:"resourceType"`
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

func NewEngineEventLoaded() Event {
	return Event{
		Type:   EventTypeEngine,
		Status: EventStatusLoaded,
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

func (es Events) getAggregatedStatus() string {
	status := EventStatusSuccess
	if es == nil {
		return status
	}
	for _, event := range es {
		if event.Status == EventStatusFetching {
			status = EventStatusFetching
		} else if event.Status == EventStatusFailed {
			status = EventStatusFailed
			break
		}
	}
	return status
}

func (es Events) getAggregatedErrors() string {
	var errors []string
	if es == nil {
		return ""
	}
	for _, event := range es {
		if event.Status == EventStatusFailed {
			errors = append(errors, event.Error)
		}
	}
	return strings.Join(errors, "\n")
}

func (es Events) getAggregatedUpdatedAt() time.Time {
	if es == nil {
		return time.Now()
	}
	updatedAt := time.Time{}
	for _, event := range es {
		if updatedAt.Before(event.UpdatedAt) {
			updatedAt = event.UpdatedAt
		}
	}
	return updatedAt
}

func (es Events) AggregateResourceEvents() Events {
	if len(es) == 0 {
		return Events(nil)
	}
	mapEvents := make(map[string]Events)
	for _, e := range es {
		mapEvents[e.ProviderName] = append(mapEvents[e.ProviderName], e)
	}
	var pes Events
	for providerName, events := range mapEvents {
		pe := Event{RunId: events[0].RunId, ProviderName: providerName, Type: EventTypeProvider}
		pe.RunId = events[0].RunId
		pe.ChildEvents = events
		pe.Status = pe.ChildEvents.getAggregatedStatus()
		pe.Error = pe.ChildEvents.getAggregatedErrors()
		pe.UpdatedAt = pe.ChildEvents.getAggregatedUpdatedAt()
		pes = append(pes, pe)
	}
	return pes
}

func (e *Event) AddChildEvents(events Events) {
	if e == nil {
		return
	}
	e.ChildEvents = events
	e.Status = e.ChildEvents.getAggregatedStatus()
	e.Error = e.ChildEvents.getAggregatedErrors()
	e.UpdatedAt = e.ChildEvents.getAggregatedUpdatedAt()
}
