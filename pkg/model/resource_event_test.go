package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceEventNewResourceEvent(t *testing.T) {
	err := errors.New("failed to fetch")
	type args struct {
		resourceType string
		isFetching   bool
		err          error
	}
	tests := []struct {
		name string
		args args
		want ResourceEvent
	}{
		{
			name: "NewResourceEventFetching",
			args: args{
				resourceType: "ec2.Volume",
				isFetching:   true,
				err:          nil,
			},
			want: ResourceEvent{
				ResourceType: "ec2.Volume",
				FetchStatus:  ResourceEventStatusFetching,
			},
		},
		{
			name: "NewResourceEventSuccess",
			args: args{
				resourceType: "ec2.Volume",
				isFetching:   false,
				err:          nil,
			},
			want: ResourceEvent{
				ResourceType: "ec2.Volume",
				FetchStatus:  ResourceEventStatusSuccess,
			},
		},
		{
			name: "NewResourceEventFailure",
			args: args{
				resourceType: "ec2.Volume",
				isFetching:   false,
				err:          err,
			},
			want: ResourceEvent{
				ResourceType: "ec2.Volume",
				FetchStatus:  ResourceEventStatusFailed,
				ErrorMessage: err.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertEqualsResourceEvent(t, tt.want, NewResourceEvent(tt.args.resourceType, tt.args.isFetching, tt.args.err))
		})
	}
}

func TestResourceEventsGetStatus(t *testing.T) {
	type args struct {
		resourceEvents ResourceEvents
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetStatusSuccess",
			args: args{
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", false, nil),
					NewResourceEvent("ec2.Instance", false, nil),
					NewResourceEvent("Lambda.Function", false, nil),
				},
			},
			want: EngineStatusSuccess,
		},
		{
			name: "GetStatusFailed",
			args: args{
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", false, nil),
					NewResourceEvent("ec2.Instance", true, nil),
					NewResourceEvent("Lambda.Function", false, errors.New("failed to fetch")),
				},
			},
			want: EngineStatusFailed,
		},
		{
			name: "GetStatusFetching",
			args: args{
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", true, nil),
					NewResourceEvent("ec2.Instance", false, nil),
					NewResourceEvent("Lambda.Function", false, nil),
				},
			},
			want: EngineStatusFetching,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.resourceEvents.GetStatus()
			assert.Equal(t, tt.want, got)
		})
	}
}
