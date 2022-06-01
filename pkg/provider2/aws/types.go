package aws

import (
	"github.com/run-x/cloudgrep/pkg/provider2/types"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

type mapper struct {
	IdField   string
	TagField  resourceconverter.TagField
	FetchFunc types.FetchFunc
	IsGlobal  bool
}

func (p *Provider) getTypeMapping() map[string]mapper {
	defaultTagField := resourceconverter.TagField{
		Key:   "Key",
		Name:  "Tags",
		Value: "Value",
	}
	return map[string]mapper{
		"ec2.Instance": {
			IdField:   "InstanceId",
			TagField:  defaultTagField,
			FetchFunc: p.FetchEC2Instances,
		},
		"ec2.Volume": {
			IdField:   "VolumeId",
			TagField:  defaultTagField,
			FetchFunc: p.FetchEBSVolumes,
		},
		"elb.LoadBalancer": {
			IdField:   "LoadBalancerArn",
			FetchFunc: p.FetchLoadBalancers,
		},
		"s3.Bucket": {
			IdField:   "Name",
			FetchFunc: p.FetchS3Buckets,
		},
		"rds.DBInstance": {
			IdField:   "DBInstanceIdentifier",
			FetchFunc: p.FetchRDSInstances,
			TagField: resourceconverter.TagField{
				Key:   "Key",
				Name:  "TagList",
				Value: "Value",
			},
		},
		"rds.DBCluster": {
			IdField:   "DBClusterIdentifier",
			FetchFunc: p.FetchRDSClusters,
			TagField: resourceconverter.TagField{
				Key:   "Key",
				Name:  "TagList",
				Value: "Value",
			},
		},
		"lambda.Function": {
			IdField:   "FunctionArn",
			FetchFunc: p.FetchLambdaFunctions,
		},
	}
}
