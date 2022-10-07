package aws

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	cfg "github.com/juandiegopalomino/cloudgrep/pkg/config"
	regionutil "github.com/juandiegopalomino/cloudgrep/pkg/provider/aws/regions"
)

// Default region to run tests against if AWS_REGION is not set.
// Always runs integration tests for the "global" region.
const defaultRegion = "us-east-1"
const globalRegion = "global"

const (
	accountIntegrationDev  = "316817240772"
	accountIntegrationProd = "438881294876"
)

// Only run the integration tests on these specially preparred accounts
var integrationAwsAccounts = []string{accountIntegrationDev, accountIntegrationProd}

// Set any of these env vars to a non-empty value to force-enable the integration tests
// (they will fail tests if creds aren't available)
var integrationTestVars = []string{
	"CLOUD_INTEGRATION_TESTS",
	// "CI",
}

// Cache the checks for credentials so it doesn't run for every test
var credCheck credChecker

type integrationTestContext struct {
	p   []Provider
	log *zap.Logger
	ctx context.Context
}

func setupIntegrationTest(t testing.TB) *integrationTestContext {
	t.Helper()

	if testing.Short() {
		t.Skip("integration tests run long")
	}

	ctx := &integrationTestContext{}
	ctx.ctx = context.Background()
	setupIntegrationLogs(t, ctx)
	setupIntegrationProvider(t, ctx)

	checkShouldRunIntegrationTests(t, ctx)

	return ctx
}

func checkShouldRunIntegrationTests(t testing.TB, ctx *integrationTestContext) {
	t.Helper()

	if len(ctx.p) == 0 {
		t.Skip("no providers configured")
		return
	}

	cfg := ctx.p[0].config

	creds := credCheck.HasAWSCreds(t, cfg)

	hasEnvVar := false
	for _, key := range integrationTestVars {
		val := os.Getenv(key)
		if len(val) > 0 {
			hasEnvVar = true
		}
	}

	if hasEnvVar && !creds {
		t.Fatalf("cannot run integration tests without creds")
	}

	if !creds {
		t.Skip("no active creds for the integration testing account")
	}
}

func setupIntegrationProvider(t testing.TB, ctx *integrationTestContext) {
	t.Helper()

	c := cfg.Provider{}
	c.Cloud = "aws"
	c.Regions = []string{regionutil.All}

	providers, err := NewProviders(ctx.ctx, c, ctx.log)
	if err != nil {
		credErrors := []string{
			"invalid AWS credentials",
			"no AWS credentials found",
		}

		for _, s := range credErrors {
			if strings.Contains(err.Error(), s) {
				t.Skipf("no valid AWS credentials present")
			}
		}

		t.Fatalf("unable to instantiate new providers: %v", err)
	}

	var awsProviders []Provider
	for _, provider := range providers {
		awsProvdier := provider.(Provider)
		awsProviders = append(awsProviders, awsProvdier)
	}

	ctx.p = awsProviders
}

func setupIntegrationLogs(t testing.TB, ctx *integrationTestContext) {
	t.Helper()
	ctx.log = zaptest.NewLogger(t)
}

type credChecker struct {
	l        sync.Mutex
	hasCreds bool
	done     bool
}

func (c *credChecker) HasAWSCreds(t testing.TB, cfg aws.Config) bool {
	c.l.Lock()
	defer c.l.Unlock()

	if c.done {
		return c.hasCreds
	}

	c.done = true

	client := sts.NewFromConfig(cfg)
	output, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		var re *awshttp.ResponseError
		if !errors.As(err, &re) {
			if strings.Contains(err.Error(), "failed to retrieve credentials") {
				return false
			}

			t.Fatalf("unknown error calling sts:GetCallerIdentity: %v", err)
		}
		if re.HTTPStatusCode() == 403 {
			// No creds
			return false
		}

		if re.HTTPStatusCode() == 400 {
			// Bad creds
			return false
		}
	}

	for _, id := range integrationAwsAccounts {
		if *output.Account == id {
			c.hasCreds = true
			return true
		}
	}

	return false
}
