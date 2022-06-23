package regions

import (
	"context"
	_ "embed"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestIsValid_global(t *testing.T) {
	assert.True(t, IsValid("global"))
}

func TestIsValid_invalid(t *testing.T) {
	assert.False(t, IsValid("foo"))
}

func TestIsValid_valid(t *testing.T) {
	assert.True(t, IsValid("us-east-1"))
}

func TestSetConfigRegion_nilConfig(t *testing.T) {
	assert.PanicsWithValue(t, "unexpected nil cfg", func() {
		SetConfigRegion(nil, nil)
	})
}

func TestSetConfigRegion_alreadySet(t *testing.T) {
	cfg := aws.Config{
		Region: "foobar",
	}

	region := mustRegion(t, "us-west-1")

	SetConfigRegion(&cfg, []Region{region})
	assert.Equal(t, "foobar", cfg.Region)
}

func TestSetConfigRegion_empty(t *testing.T) {
	cfg := aws.Config{}

	SetConfigRegion(&cfg, nil)
	assert.Equal(t, "us-east-1", cfg.Region)
}

func TestSetConfigRegion_onlyGlobal(t *testing.T) {
	cfg := aws.Config{}

	region := mustRegion(t, "global")

	SetConfigRegion(&cfg, []Region{region})
	assert.Equal(t, "us-east-1", cfg.Region)
}

func TestSetConfigRegion_multi(t *testing.T) {
	cfg := aws.Config{}

	regions := []Region{
		mustRegion(t, "global"),
		mustRegion(t, "us-west-1"),
		mustRegion(t, "us-west-2"),
	}

	SetConfigRegion(&cfg, regions)
	assert.Equal(t, "us-west-1", cfg.Region)
}

func TestSelectRegions_configured(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}
	configured := []string{"us-east-1", "us-west-1"}

	expected := slices.Clone(configured)

	regions, err := SelectRegions(ctx, configured, cfg)
	assert.NoError(t, err)
	ids := regionIds(regions)
	assert.Equal(t, expected, ids)
}

func TestSelectRegions_configuredInvalid(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}
	configured := []string{"foobar"}

	regions, err := SelectRegions(ctx, configured, cfg)
	assert.ErrorContains(t, err, "unable to use configured regions: invalid AWS region: foobar")
	assert.Empty(t, regions)
}

func TestSelectRegions_config(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{Region: "us-west-1"}

	regions, err := SelectRegions(ctx, nil, cfg)
	assert.NoError(t, err)
	ids := regionIds(regions)
	assert.ElementsMatch(t, []string{"global", "us-west-1"}, ids)
}

func TestSelectRegions_configInvalid(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{Region: "foobar"}

	regions, err := SelectRegions(ctx, nil, cfg)
	assert.ErrorContains(t, err, "invalid AWS region: foobar")
	assert.Empty(t, regions)
}

func TestSelectRegions_promptCanceled(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	regions, err := SelectRegions(ctx, nil, cfg)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, regions)
}

func TestSelectRegions_allMultiple(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}
	configured := []string{"us-east-1", "all"}

	regions, err := SelectRegions(ctx, configured, cfg)
	assert.ErrorContains(t, err, "can only use 'all' as a region if it is the only configured region")
	assert.Empty(t, regions)
}

//go:embed testdata/describe_regions_resp.xml
var describeRegionsResponse []byte

//go:embed testdata/get_caller_identity_resp.xml
var getCallerIdentityResponse []byte

func TestSelectRegions_all(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}
	configured := []string{"all"}

	httpClient := sequencedHttpClient{
		clients: []aws.HTTPClient{
			&mockedHttpClient{
				contentType: "text/xml;charset=UTF-8",
				body:        getCallerIdentityResponse,
			},
			&mockedHttpClient{
				contentType: "text/xml;charset=UTF-8",
				body:        describeRegionsResponse,
			},
		},
	}

	cfg.HTTPClient = &httpClient
	cfg.Credentials = mockedCredentialsProvider{}

	expected := []string{"us-east-1", "eu-west-1", "global"}

	regions, err := SelectRegions(ctx, configured, cfg)
	assert.NoError(t, err)
	ids := regionIds(regions)
	assert.ElementsMatch(t, expected, ids)
}
