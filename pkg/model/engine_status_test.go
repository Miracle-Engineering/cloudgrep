package model

import (
	"testing"
)

func TestEngineStatusNewEngineStatus(t *testing.T) {
	type args struct {
		resourceEvents   ResourceEvents
		failedAtCreation bool
	}
	tests := []struct {
		name string
		args args
		want EngineStatus
	}{
		{
			name: "failedAtCreation",
			args: args{
				resourceEvents:   ResourceEvents(nil),
				failedAtCreation: true,
			},
			want: EngineStatus{
				FetchStatus:    EngineStatusFailedAtCreation,
				ResourceEvents: ResourceEvents(nil),
			},
		},
		{
			name: "success",
			args: args{
				resourceEvents: ResourceEvents{
					{
						FetchStatus:  ResourceEventStatusSuccess,
						ResourceType: "ec2.Volume",
						ErrorMessage: "",
					},
					{
						FetchStatus:  ResourceEventStatusSuccess,
						ResourceType: "ec2.Instance",
						ErrorMessage: "",
					},
				},
				failedAtCreation: false,
			},
			want: EngineStatus{
				FetchStatus: EngineStatusSuccess,
				ResourceEvents: ResourceEvents{
					{
						FetchStatus:  ResourceEventStatusSuccess,
						ResourceType: "ec2.Volume",
						ErrorMessage: "",
					},
					{
						FetchStatus:  ResourceEventStatusSuccess,
						ResourceType: "ec2.Instance",
						ErrorMessage: "",
					},
				},
			},
		},
		{
			name: "failed",
			args: args{
				resourceEvents: ResourceEvents{
					{
						FetchStatus:  ResourceEventStatusSuccess,
						ResourceType: "ec2.Volume",
						ErrorMessage: "",
					},
					{
						FetchStatus:  ResourceEventStatusFailed,
						ResourceType: "ec2.Instance",
						ErrorMessage: "Unable to fetch EC2.Instances",
					},
				},
				failedAtCreation: false,
			},
			want: EngineStatus{
				FetchStatus: EngineStatusFailed,
				ResourceEvents: ResourceEvents{
					{
						FetchStatus:  ResourceEventStatusSuccess,
						ResourceType: "ec2.Volume",
						ErrorMessage: "",
					},
					{
						FetchStatus:  ResourceEventStatusFailed,
						ResourceType: "ec2.Instance",
						ErrorMessage: "Unable to fetch EC2.Instances",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertEqualsEngineStatus(t, tt.want, NewEngineStatus(tt.args.resourceEvents, tt.args.failedAtCreation))
		})
	}
}
