package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProviderStatus(t *testing.T) {
	type args struct {
		providerType   string
		resourceEvents ResourceEvents
		err            error
	}
	tests := []struct {
		name string
		args args
		want ProviderStatus
	}{
		{
			name: "ProviderStatusSuccess",
			args: args{
				providerType: "aws",
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
				},
				err: nil,
			},
			want: ProviderStatus{
				ProviderType: "aws",
				FetchStatus:  ProviderStatusSuccess,
				ErrorMessage: "",
				ResourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
				},
			},
		},
		{
			name: "ProviderStatusFailed",
			args: args{
				providerType: "aws",
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", true, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, errors.New("failed to fetch")),
				},
				err: nil,
			},
			want: ProviderStatus{
				ProviderType: "aws",
				FetchStatus:  ProviderStatusFailed,
				ErrorMessage: "",
				ResourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", true, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, errors.New("failed to fetch")),
				},
			},
		},
		{
			name: "ProviderStatusFetching",
			args: args{
				providerType: "aws",
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", true, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
				},
				err: nil,
			},
			want: ProviderStatus{
				ProviderType: "aws",
				FetchStatus:  ProviderStatusFetching,
				ErrorMessage: "",
				ResourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", true, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewProviderStatus(tt.args.providerType, tt.args.resourceEvents, tt.args.err)
			AssertEqualsProviderStatus(t, got, tt.want)
		})
	}
}

func TestProviderStatusesGetAggregatedStatus(t *testing.T) {
	tests := []struct {
		name             string
		providerStatuses ProviderStatuses
		want             string
	}{
		{
			name: "ProviderStatusSuccess",
			providerStatuses: ProviderStatuses{
				ProviderStatus{ProviderType: "aws1", FetchStatus: ProviderStatusSuccess},
				ProviderStatus{ProviderType: "aws2", FetchStatus: ProviderStatusSuccess},
				ProviderStatus{ProviderType: "aws3", FetchStatus: ProviderStatusSuccess},
			},
			want: ProviderStatusSuccess,
		},
		{
			name: "ProviderStatusFailed",
			providerStatuses: ProviderStatuses{
				ProviderStatus{ProviderType: "aws1", FetchStatus: ProviderStatusSuccess},
				ProviderStatus{ProviderType: "aws2", FetchStatus: ProviderStatusFetching},
				ProviderStatus{ProviderType: "aws3", FetchStatus: ProviderStatusFailed},
			},
			want: ProviderStatusFailed,
		},
		{
			name: "ProviderStatusFetching",
			providerStatuses: ProviderStatuses{
				ProviderStatus{ProviderType: "aws1", FetchStatus: ProviderStatusFetching},
				ProviderStatus{ProviderType: "aws2", FetchStatus: ProviderStatusSuccess},
				ProviderStatus{ProviderType: "aws3", FetchStatus: ProviderStatusSuccess},
			},
			want: ProviderStatusFetching,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.providerStatuses.GetAggregatedStatus())
		})
	}
}
