package mapper

import (
	"context"
	"reflect"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNewMapperOk(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	cfg := Config{
		Mappings: []Mapping{
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
		},
	}
	provider := NewTestProvider(cfg)
	mapper, err := New(cfg, *logger, reflect.ValueOf(provider))
	assert.NoError(t, err)

	instances, err := provider.FetchTestInstances(ctx)
	assert.NoError(t, err)
	//test expected resource
	r1 := model.Resource{
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
	_r1, err := mapper.ToResource(instances[0], "us-west-2")
	assert.NoError(t, err)
	util.AssertEqualsResource(t, r1, _r1)
}

func TestNewMapperError(t *testing.T) {
	logger := zaptest.NewLogger(t)
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
	provider := NewTestProvider(cfg)
	assert.PanicsWithError(t,
		"could not find a method called 'Unknown' on 'mapper.TestProvider'",
		func() { New(cfg, *logger, reflect.ValueOf(provider)) },
	)

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
	provider = NewTestProvider(cfg)
	assert.PanicsWithError(t,
		"method WrongReturnType has invalid return type, expecting ([]any, error)",
		func() { New(cfg, *logger, reflect.ValueOf(provider)) },
	)
}

type TestProvider struct {
	Config
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

func NewTestProvider(c Config) TestProvider {
	p := TestProvider{
		Config: c,
	}
	return p
}

func (p TestProvider) GetMapperConfig() Config {
	return p.Config
}

func (p TestProvider) Region() string {
	return "us-east-1"
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

func (TestProvider) WrongReturnType(ctx context.Context) []TestInstance {
	return []TestInstance{}
}
