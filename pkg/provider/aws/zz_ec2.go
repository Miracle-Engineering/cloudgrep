package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_ec2(mapping map[string]mapper) {
	mapping["ec2.Instance"] = mapper{
		FetchFunc: p.fetch_ec2_Instance,
		IdField:   "InstanceId",
		IsGlobal:  false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetch_ec2_Instance(ctx context.Context, output chan<- model.Resource) error {
	client := ec2.NewFromConfig(p.config)

	input := &ec2.DescribeInstancesInput{}

	paginator := ec2.NewDescribeInstancesPaginator(client, input)

	resourceConverter := p.converterFor("ec2.Instance")
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EC2 Instances: %w", err)
		}

		for _, item_0 := range page.Reservations {
			if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, item_0.Instances); err != nil {
				return err
			}
		}
	}

	return nil
}
