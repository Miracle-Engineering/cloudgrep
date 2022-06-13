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
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	cfg "github.com/run-x/cloudgrep/pkg/config"
)

// Set this env var to force enable the integration tests (will fail tests if creds aren't available)
const testEnvVar = "CLOUD_INTEGRATION_TESTS"
const ciEnvVar = "CI"

const (
	accountIntegrationDev  = "316817240772"
	accountIntegrationProd = "438881294876"
)

// Only run the integration tests on these specially preparred accounts
var integrationAwsAccounts = []string{accountIntegrationDev, accountIntegrationProd}

// Cache the checks for credentials so it doesn't run for every test
var credCheck credChecker

type integrationTestContext struct {
	p   *Provider
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

	creds := credCheck.HasAWSCreds(t, ctx.p.config)
	_, hasIntegrationEnvVar := os.LookupEnv(testEnvVar)
	_, hasCiEnvVar := os.LookupEnv(ciEnvVar)

	hasEnvVar := hasIntegrationEnvVar || hasCiEnvVar

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
	c.Regions = []string{
		testingutil.TestRegion,
	}

	providers, err := NewProviders(ctx.ctx, c, ctx.log)
	if err != nil {
		t.Fatalf("unable to instantiate new providers: %v", err)
	}

	if len(providers) != 1 {
		t.Fatal("currently only have support for single provider in tests")
	}

	provider := providers[0].(Provider)

	ctx.p = &provider
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
