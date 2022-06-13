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
		provider     string
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
				provider:     "AWS Account: 000000000000, Region: Mock",
				isFetching:   true,
				err:          nil,
			},
			want: ResourceEvent{
				ResourceType: "ec2.Volume",
				Provider:     "AWS Account: 000000000000, Region: Mock",
				FetchStatus:  ResourceEventStatusFetching,
			},
		},
		{
			name: "NewResourceEventSuccess",
			args: args{
				resourceType: "ec2.Volume",
				provider:     "AWS Account: 000000000000, Region: Mock",
				isFetching:   false,
				err:          nil,
			},
			want: ResourceEvent{
				ResourceType: "ec2.Volume",
				Provider:     "AWS Account: 000000000000, Region: Mock",
				FetchStatus:  ResourceEventStatusSuccess,
			},
		},
		{
			name: "NewResourceEventFailure",
			args: args{
				resourceType: "ec2.Volume",
				provider:     "AWS Account: 000000000000, Region: Mock",
				isFetching:   false,
				err:          err,
			},
			want: ResourceEvent{
				ResourceType: "ec2.Volume",
				Provider:     "AWS Account: 000000000000, Region: Mock",
				FetchStatus:  ResourceEventStatusFailed,
				ErrorMessage: err.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertEqualsResourceEvent(t, tt.want, NewResourceEvent(tt.args.resourceType, tt.args.provider, tt.args.isFetching, tt.args.err))
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
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
				},
			},
			want: ProviderStatusSuccess,
		},
		{
			name: "GetStatusFailed",
			args: args{
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", true, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, errors.New("failed to fetch")),
				},
			},
			want: ProviderStatusFailed,
		},
		{
			name: "GetStatusFetching",
			args: args{
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", true, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
				},
			},
			want: ProviderStatusFetching,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.resourceEvents.GetAggregatedStatus()
			assert.Equal(t, tt.want, got)
		})
	}
}
