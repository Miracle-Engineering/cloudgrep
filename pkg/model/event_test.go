package model

import (
	_ "embed"
	"errors"

	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEvent(t *testing.T, ee, ae Event) {
	assert.Equal(t, ee.Type, ae.Type)
	assert.Equal(t, ee.Status, ae.Status)
	assert.Equal(t, ee.ProviderName, ae.ProviderName)
	assert.Equal(t, ee.ResourceType, ae.ResourceType)
	assert.Equal(t, ee.Error, ae.Error)
	if ee.ChildEvents != nil {
		assert.Equal(t, len(ee.ChildEvents), len(ae.ChildEvents))
		for _, e := range ee.ChildEvents {
			for _, a := range ae.ChildEvents {
				switch e.Type {
				case EventTypeEngine:
					if e.Type == a.Type {
						assertEvent(t, e, a)
					}
				case EventTypeProvider:
					if e.Type == a.Type && e.ProviderName == a.ProviderName {
						assertEvent(t, e, a)
					}
				case EventTypeResource:
					if e.Type == a.Type && e.ProviderName == a.ProviderName && e.ResourceType == a.ResourceType {
						assertEvent(t, e, a)
					}
				}
			}
		}
	}
}

func TestNewEngineEventStart(t *testing.T) {
	t.Run("NewEngineEventStart", func(t *testing.T) {
		e := NewEngineEventStart()
		assert.Equal(t, EventTypeEngine, e.Type)
		assert.Equal(t, EventStatusFetching, e.Status)
		assert.NotEmpty(t, e.RunId)
	},
	)

}

func TestNewEngineEventEnd(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "NewEngineEventEndNoError",
			args: args{
				err: nil,
			},
			want: Event{
				Type:   EventTypeEngine,
				Status: EventStatusSuccess,
			},
		},
		{
			name: "NewProviderEventEndWithError",
			args: args{
				err: errors.New("mock error"),
			},
			want: Event{
				Type:   EventTypeEngine,
				Status: EventStatusFailed,
				Error:  "mock error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngineEventEnd(tt.args.err)
			assertEvent(t, tt.want, e)
		},
		)
	}
}

func TestNewProviderEventStart(t *testing.T) {
	t.Run("NewProviderEventStart", func(t *testing.T) {
		mockProvider := "mockProvider"
		e := NewProviderEventStart(mockProvider)
		assert.Equal(t, EventTypeProvider, e.Type)
		assert.Equal(t, EventStatusFetching, e.Status)
		assert.Equal(t, mockProvider, e.ProviderName)
		assert.Empty(t, e.RunId)
	},
	)
}

func TestNewProviderEventEnd(t *testing.T) {
	type args struct {
		providerName string
		err          error
	}

	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "NewProviderEventEndNoError",
			args: args{
				providerName: "mockProvider",
				err:          nil,
			},
			want: Event{
				Type:         EventTypeProvider,
				Status:       EventStatusSuccess,
				ProviderName: "mockProvider",
			},
		},
		{
			name: "NewProviderEventEndWithError",
			args: args{
				providerName: "mockProvider",
				err:          errors.New("mock error"),
			},
			want: Event{
				Type:         EventTypeProvider,
				Status:       EventStatusFailed,
				ProviderName: "mockProvider",
				Error:        "mock error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewProviderEventEnd(tt.args.providerName, tt.args.err)
			assertEvent(t, tt.want, e)
		},
		)
	}
}

func TestNewResourceEventStart(t *testing.T) {
	t.Run("NewResourceEventStart", func(t *testing.T) {
		mockProvider := "mockProvider"
		mockResourceType := "mockResourceType"
		e := NewResourceEventStart(mockProvider, mockResourceType)
		assert.Equal(t, EventTypeResource, e.Type)
		assert.Equal(t, EventStatusFetching, e.Status)
		assert.Equal(t, mockProvider, e.ProviderName)
		assert.Equal(t, mockResourceType, e.ResourceType)
		assert.Empty(t, e.RunId)
	},
	)
}

func TestNewResourceEventEnd(t *testing.T) {
	type args struct {
		providerName string
		resourceType string
		err          error
	}

	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "NewResourceEventEndNoError",
			args: args{
				providerName: "mockProvider",
				resourceType: "mockResourceType",
				err:          nil,
			},
			want: Event{
				Type:         EventTypeResource,
				Status:       EventStatusSuccess,
				ProviderName: "mockProvider",
				ResourceType: "mockResourceType",
			},
		},
		{
			name: "NewResourceEventEndWithError",
			args: args{
				providerName: "mockProvider",
				resourceType: "mockResourceType",
				err:          errors.New("mock error"),
			},
			want: Event{
				Type:         EventTypeResource,
				Status:       EventStatusFailed,
				ProviderName: "mockProvider",
				ResourceType: "mockResourceType",
				Error:        "mock error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewResourceEventEnd(tt.args.providerName, tt.args.resourceType, tt.args.err)
			assertEvent(t, tt.want, e)
		},
		)
	}
}
