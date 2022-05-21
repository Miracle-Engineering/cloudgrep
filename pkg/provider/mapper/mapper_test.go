package mapper

import (
	"context"
	"reflect"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewMapperOk(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	cfg := Config{
		Mappings: []Mapping{
			//most common case - use fields for id, properties, tags
			{
				Type:         "github.com/run-x/cloudgrep/pkg/provider/mapper.TestInstance",
				ResourceType: "mapper.Instance",
				IdField:      "Id",
				Impl:         "FetchTestInstances",
				TagField: TagField{
					Name:  "SomeTags",
					Key:   "Name",
					Value: "Val",
				},
			},
			//this test has a custom method to generate tags
			{
				Type:         "github.com/run-x/cloudgrep/pkg/provider/mapper.TestLoadBalancer",
				ResourceType: "mapper.LoadBalancer",
				IdField:      "Arn",
				Impl:         "FetchTestLoadBalancers",
				TagImpl:      "FetchTestLoadBalancerTags",
			},
		},
	}
	provider := TestProvider{}
	mapper, err := new(cfg, *logger, reflect.ValueOf(provider))
	assert.NoError(t, err)

	instances, err := provider.FetchTestInstances(ctx)
	assert.NoError(t, err)

	//test expected resource
	expectedInstance := model.Resource{
		Id: "i-1", Region: "us-west-2", Type: "mapper.Instance",
		Tags: []model.Tag{
			{Key: "tag1", Value: "val1"},
			{Key: "tag2", Value: "val2"},
		},
		Properties: []model.Property{
			{Name: "Id", Value: "i-1"},
			{Name: "Value", Value: "abc"},
		},
	}
	instance, err := mapper.ToResource(ctx, instances[0], "us-west-2")
	assert.NoError(t, err)
	model.AssertEqualsResource(t, expectedInstance, instance)

	expectedLB := model.Resource{
		Id:     "arn:aws:elasticloadbalancing:us-east-1:0123456789:loadbalancer/net/my-load-balancer/14522ba1bd959dd6",
		Region: "us-west-2",
		Type:   "mapper.LoadBalancer",
		Tags: []model.Tag{
			{Key: "tag1", Value: "val1"},
			{Key: "tag2", Value: "val2"},
		},
		Properties: []model.Property{
			{Name: "Arn", Value: "arn:aws:elasticloadbalancing:us-east-1:0123456789:loadbalancer/net/my-load-balancer/14522ba1bd959dd6"},
			{Name: "DNSName", Value: "my-load-balancer-14522ba1bd959dd6.elb.us-east-1.amazonaws.com"},
		},
	}
	loadBalancers, err := provider.FetchTestLoadBalancers(ctx)
	assert.NoError(t, err)
	resourceLB, err := mapper.ToResource(ctx, loadBalancers[0], "us-west-2")
	assert.NoError(t, err)
	model.AssertEqualsResource(t, expectedLB, resourceLB)
}

func TestNewMapperError(t *testing.T) {
	logger := zaptest.NewLogger(t)

	//Impl references an unknown method
	cfg := Config{
		Mappings: []Mapping{
			{
				Type:         "github.com/run-x/cloudgrep/pkg/provider/mapper.TestInstance",
				ResourceType: "mapper.Instance",
				IdField:      "Id",
				Impl:         "Unknown",
				TagField: TagField{
					Name:  "SomeTags",
					Key:   "Name",
					Value: "Val",
				},
			},
		},
	}
	assert.PanicsWithError(t,
		"could not find a method called 'Unknown' on 'mapper.TestProvider'",
		func() {
			new(cfg, *logger, reflect.ValueOf(TestProvider{})) //nolint
		},
	)

	//Impl references a method with wrong return type
	cfg = Config{
		Mappings: []Mapping{
			{
				Type:         "github.com/run-x/cloudgrep/pkg/provider/mapper.TestInstance",
				ResourceType: "mapper.Instance",
				IdField:      "Id",
				Impl:         "WrongReturnType",
				TagField: TagField{
					Name:  "SomeTags",
					Key:   "Name",
					Value: "Val",
				},
			},
		},
	}
	assert.PanicsWithError(t,
		"method WrongReturnType has invalid return type, expecting ([]any, error)",
		func() {
			new(cfg, *logger, reflect.ValueOf(TestProvider{})) //nolint
		},
	)

	//Impl references an unknown method for tag
	cfg = Config{
		Mappings: []Mapping{
			{
				Type:         "github.com/run-x/cloudgrep/pkg/provider/mapper.TestLoadBalancer",
				ResourceType: "mapper.LoadBalancer",
				IdField:      "Arn",
				Impl:         "FetchTestLoadBalancers",
				TagImpl:      "Unknown",
			},
		},
	}
	assert.PanicsWithError(t,
		"could not find a method called 'Unknown' on 'mapper.TestProvider'",
		func() {
			new(cfg, *logger, reflect.ValueOf(TestProvider{})) //nolint
		},
	)
}

type TestProvider struct {
	// Config
}

type TestInstance struct {
	Id       string
	Value    string
	SomeTags []TestTag
}
type TestTag struct {
	Name string
	Val  string
}

type TestLoadBalancer struct {
	Arn     string
	DNSName string
}

type Return string

func (TestProvider) FetchTestInstances(ctx context.Context) ([]TestInstance, error) {
	return []TestInstance{
		{
			Id:    "i-1",
			Value: "abc",
			SomeTags: []TestTag{
				{Name: "tag1", Val: "val1"},
				{Name: "tag2", Val: "val2"},
			},
		}, {
			Id:    "i-2",
			Value: "edf",
		},
	}, nil
}

func (TestProvider) FetchTestLoadBalancers(ctx context.Context) ([]TestLoadBalancer, error) {
	return []TestLoadBalancer{
		{
			Arn:     "arn:aws:elasticloadbalancing:us-east-1:0123456789:loadbalancer/net/my-load-balancer/14522ba1bd959dd6",
			DNSName: "my-load-balancer-14522ba1bd959dd6.elb.us-east-1.amazonaws.com",
		},
	}, nil
}

func (TestProvider) FetchTestLoadBalancerTags(ctx context.Context, lb TestLoadBalancer) (model.Tags, error) {
	return model.Tags{
		model.Tag{Key: "tag1", Value: "val1"},
		model.Tag{Key: "tag2", Value: "val2"},
	}, nil
}

func (TestProvider) WrongReturnType(ctx context.Context) []TestInstance {
	return []TestInstance{}
}
