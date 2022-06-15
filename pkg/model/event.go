package model

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	EventStatusFetching string = "fetching"
	EventStatusFailed          = "failed"
	EventStatusSuccess         = "success"
)

const (
	EventTypeEngine   string = "engine"
	EventTypeProvider        = "provider"
	EventTypeResource        = "resource"
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

func NewEvent(eventType string, providerName string, resourceType string) Event {
	event := Event{}
	event.Type = eventType
	event.Status = EventStatusFetching
	switch eventType {
	case EventTypeEngine:
		event.RunId = uuid.New().String()
		break
	case EventTypeProvider:
		event.ProviderName = providerName
		break
	case EventTypeResource:
		event.ProviderName = providerName
		event.ResourceType = resourceType
		break
	}
	return event
}

func (e *Event) UpdateError(err error) {
	if err != nil {
		e.Status = EventStatusFailed
		e.Error = err.Error()
	} else {
		e.Status = EventStatusSuccess
	}
}

func (e *Event) HasError() bool {
	return e.Error != ""
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
			//fmt.Sprintf("%s\n%s", errors, event.Error)
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
	if es == nil || len(es) == 0 {
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
