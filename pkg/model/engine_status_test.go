package model

import (
	"testing"
)

func TestEngineStatusGetEngineStatus(t *testing.T) {
	type args struct {
		providerStatuses ProviderStatuses
		resourceEvents   ResourceEvents
	}
	tests := []struct {
		name string
		args args
		want EngineStatus
	}{
		{
			name: "no data",
			args: args{
				providerStatuses: ProviderStatuses(nil),
				resourceEvents:   ResourceEvents(nil),
			},
			want: EngineStatus{},
		},
		{
			name: "failed1",
			args: args{
				providerStatuses: ProviderStatuses{
					ProviderStatus{
						ProviderType: "provider1",
						FetchStatus:  ProviderStatusFailed,
					},
					ProviderStatus{
						ProviderType: "provider2",
						FetchStatus:  ProviderStatusFailed,
					},
				},
				resourceEvents: ResourceEvents(nil),
			},
			want: EngineStatus{
				FetchStatus: EngineStatusFailed,
				ProviderStatuses: ProviderStatuses{
					ProviderStatus{
						ProviderType: "provider1",
						FetchStatus:  ProviderStatusFailed,
					},
					ProviderStatus{
						ProviderType: "provider2",
						FetchStatus:  ProviderStatusFailed,
					},
				},
			},
		},
		{
			name: "failed2",
			args: args{
				providerStatuses: ProviderStatuses{
					ProviderStatus{
						ProviderType: "provider1",
						FetchStatus:  ProviderStatusFailed,
					},
					ProviderStatus{
						ProviderType: "provider2",
						FetchStatus:  ProviderStatusFailed,
					},
				},
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
				},
			},
			want: EngineStatus{
				FetchStatus: EngineStatusFailed,
				ProviderStatuses: ProviderStatuses{
					ProviderStatus{
						ProviderType: "provider1",
						FetchStatus:  ProviderStatusFailed,
					},
					ProviderStatus{
						ProviderType: "provider2",
						FetchStatus:  ProviderStatusFailed,
					},
					NewProviderStatus(
						"AWS Account: 000000000000, Region: Mock",
						ResourceEvents{
							NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
							NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
							NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
						},
						nil),
				},
			},
		},
		{
			name: "success",
			args: args{
				providerStatuses: ProviderStatuses(nil),
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("mock.resourceType.1", "Provider Account: Provider ID, Region: Mock", false, nil),
					NewResourceEvent("mock.resourceType.2", "Provider Account: Provider ID, Region: Mock", false, nil),
					NewResourceEvent("mock.resourceType.3", "Provider Account: Provider ID, Region: Mock", false, nil),
				},
			},
			want: EngineStatus{
				FetchStatus: EngineStatusSuccess,
				ProviderStatuses: ProviderStatuses{
					NewProviderStatus(
						"AWS Account: 000000000000, Region: Mock",
						ResourceEvents{
							NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
							NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", false, nil),
							NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
						},
						nil),
					NewProviderStatus(
						"Provider Account: Provider ID, Region: Mock",
						ResourceEvents{
							NewResourceEvent("mock.resourceType.1", "Provider Account: Provider ID, Region: Mock", false, nil),
							NewResourceEvent("mock.resourceType.2", "Provider Account: Provider ID, Region: Mock", false, nil),
							NewResourceEvent("mock.resourceType.3", "Provider Account: Provider ID, Region: Mock", false, nil),
						},
						nil),
				},
			},
		},
		{
			name: "fetching",
			args: args{
				providerStatuses: ProviderStatuses(nil),
				resourceEvents: ResourceEvents{
					NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", true, nil),
					NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
					NewResourceEvent("mock.resourceType.1", "Provider Account: Provider ID, Region: Mock", false, nil),
					NewResourceEvent("mock.resourceType.2", "Provider Account: Provider ID, Region: Mock", false, nil),
					NewResourceEvent("mock.resourceType.3", "Provider Account: Provider ID, Region: Mock", false, nil),
				},
			},
			want: EngineStatus{
				FetchStatus: EngineStatusFetching,
				ProviderStatuses: ProviderStatuses{
					NewProviderStatus(
						"AWS Account: 000000000000, Region: Mock",
						ResourceEvents{
							NewResourceEvent("ec2.Volume", "AWS Account: 000000000000, Region: Mock", false, nil),
							NewResourceEvent("ec2.Instance", "AWS Account: 000000000000, Region: Mock", true, nil),
							NewResourceEvent("Lambda.Function", "AWS Account: 000000000000, Region: Mock", false, nil),
						},
						nil),
					NewProviderStatus(
						"Provider Account: Provider ID, Region: Mock",
						ResourceEvents{
							NewResourceEvent("mock.resourceType.1", "Provider Account: Provider ID, Region: Mock", false, nil),
							NewResourceEvent("mock.resourceType.2", "Provider Account: Provider ID, Region: Mock", false, nil),
							NewResourceEvent("mock.resourceType.3", "Provider Account: Provider ID, Region: Mock", false, nil),
						},
						nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertEqualsEngineStatus(t, tt.want, GetEngineStatus(tt.args.providerStatuses, tt.args.resourceEvents))
		})
	}
}
