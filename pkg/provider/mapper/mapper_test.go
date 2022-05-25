package mapper

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"strings"
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
		"method WrongReturnType has invalid signature; expecting one of [func(context.Context) ([]T, error), func(context.Context, chan<- T) error]",
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

func TestFetchResourcesAsync(t *testing.T) {
	typeID := typeStr(TestInstance{})

	config := buildMapperConfig()
	provider := TestProvider{}
	mapper, err := new(config, *zaptest.NewLogger(t), reflect.ValueOf(provider))
	if err != nil {
		t.Fatalf("cannot create mapper: %v", err)
	}

	if !isFetchMethodAsync(mapper.Mappings[typeID].Method.Type()) {
		t.Fatalf("method %s for type %s not async", mapper.Mappings[typeID].Method, typeID)
	}

	resources, err := mapper.fetchResourcesAsync(context.Background(), mapper.Mappings[typeID], "foo")
	if err != nil {
		t.Fatalf("error calling fetchResourcesAsync: %v", err)
	}

	assert.Len(t, resources, 2)
	assert.Equal(t, resources[0].Id, "i-1")
	assert.Equal(t, resources[1].Id, "i-2")
}

func TestFetchResourcesAsyncCanceled(t *testing.T) {
	typeID := typeStr(TestInstance{})

	config := buildMapperConfig()
	provider := TestProvider{}
	mapper, err := new(config, *zaptest.NewLogger(t), reflect.ValueOf(provider))
	if err != nil {
		t.Fatalf("cannot create mapper: %v", err)
	}

	if !isFetchMethodAsync(mapper.Mappings[typeID].Method.Type()) {
		t.Fatalf("method %s for type %s not async", mapper.Mappings[typeID].Method, typeID)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cancel()
	resources, err := mapper.fetchResourcesAsync(ctx, mapper.Mappings[typeID], "foo")
	if assert.Error(t, err) {
		assert.Equal(t, context.Canceled, err)
	}

	assert.Len(t, resources, 0)
}

func TestIsFetchMethodAsync(t *testing.T) {
	tests := []struct {
		name     string
		f        any
		expected bool
	}{
		{
			name:     "good",
			f:        func(_ context.Context, _ chan<- any) (a error) { return },
			expected: true,
		},
		{
			name: "not a func",
			f:    1,
		},
		{
			name: "single param",
			f:    func(_ context.Context) (a error) { return },
		},
		{
			name: "first param not context",
			f:    func(_ int, _ chan<- any) (a error) { return },
		},
		{
			name: "second param not chan",
			f:    func(_ context.Context, _ any) (a error) { return },
		},
		{
			name: "chan dir is recv",
			f:    func(_ context.Context, _ <-chan any) (a error) { return },
		},
		{
			name: "chan dir is both",
			f:    func(_ context.Context, _ chan any) (a error) { return },
		},
		{
			name: "return not error",
			f:    func(_ context.Context, _ chan<- any) (a any) { return },
		},
		{
			name: "multiple returns",
			f:    func(_ context.Context, _ chan<- any) (a error, b error) { return },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			typ := reflect.TypeOf(test.f)
			actual := isFetchMethodAsync(typ)
			if actual != test.expected {
				var result string
				if actual {
					result = "func is async fetch func"
				} else {
					result = "func not async fetch func"
				}
				t.Fatalf("%s; signature: %s", result, funcSignature(typ))
			}
		})
	}
}

func TestFindTagMethod(t *testing.T) {
	provider := reflect.ValueOf(TestProvider{})

	assert.PanicsWithError(
		t,
		"method FetchTestInstancesAsync has invalid signature; expecting func(context.Context, T) (model.Tags, error)",
		func() {
			findTagMethod(provider, "FetchTestInstancesAsync")
		},
	)

	tagMethod := findTagMethod(provider, "FetchTestLoadBalancerTags")
	if tagMethod.IsZero() {
		t.Error("findTagMethod returned zero for method FetchTestLoadBalancerTags")
	}
}

func TestIsFetchTagSync(t *testing.T) {
	tests := []struct {
		name     string
		f        any
		expected bool
	}{
		{
			name:     "good",
			f:        func(_ context.Context, _ any) (a model.Tags, b error) { return },
			expected: true,
		},
		{
			name: "not a func",
			f:    1,
		},
		{
			name: "single param",
			f:    func(_ context.Context) (a model.Tags, b error) { return },
		},
		{
			// This one is okay because `any` type accepts assignment from context.Context
			name:     "first param any",
			f:        func(_ any, _ any) (a model.Tags, b error) { return },
			expected: true,
		},
		{
			name: "first param not context",
			f:    func(_ int, _ any) (a model.Tags, b error) { return },
		},
		{
			name: "first return not model.Tags",
			f:    func(_ context.Context, _ any) (a any, b error) { return },
		},
		{
			name: "second return not error",
			f:    func(_ context.Context, _ any) (a model.Tags, b any) { return },
		},
		{
			name: "one return",
			f:    func(_ context.Context, _ any) (a model.Tags) { return },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			typ := reflect.TypeOf(test.f)
			actual := isFetchTagSync(typ)
			if actual != test.expected {
				var result string
				if actual {
					result = "func is tag func"
				} else {
					result = "func not tag func"
				}
				t.Fatalf("%s; signature: %s", result, funcSignature(typ))
			}
		})
	}
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

func (TestProvider) FetchTestInstancesAsync(ctx context.Context, output chan<- TestInstance) error {
	resources := []TestInstance{
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
	}

	for _, resource := range resources {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case output <- resource:
		}
	}

	return nil
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

func buildMapperConfig() Config {
	config := Config{
		TagField: TagField{
			Name:  "SomeTags",
			Key:   "Name",
			Value: "Val",
		},
		Mappings: []Mapping{
			{
				Type:         typeStr(TestInstance{}),
				ResourceType: "test.InstanceAsync",
				IdField:      "Id",
				Impl:         "FetchTestInstancesAsync",
			},
			{
				Type:         typeStr(TestLoadBalancer{}),
				ResourceType: "mapper.LoadBalancer",
				IdField:      "Arn",
				Impl:         "FetchTestLoadBalancers",
				TagImpl:      "FetchTestLoadBalancerTags",
			},
		},
	}

	return config
}

func TestTypeStr(t *testing.T) {
	assert.Equal(t, "github.com/run-x/cloudgrep/pkg/provider/mapper.TestInstance", typeStr(TestInstance{}))
}

func typeStr(v any) string {
	t := reflect.TypeOf(v)
	return fmt.Sprintf("%v/%v", path.Dir(t.PkgPath()), t.String())
}

func funcSignature(t reflect.Type) string {
	// Source: https://stackoverflow.com/a/54129236
	if t.Kind() != reflect.Func {
		return "<not a function>"
	}

	buf := strings.Builder{}
	buf.WriteString("func (")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(t.In(i).String())
	}
	buf.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}
