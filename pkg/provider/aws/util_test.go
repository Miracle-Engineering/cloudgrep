package aws

import (
	"bytes"
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Set this env var to force enable the integration tests (will fail tests if creds aren't available)
const testEnvVar = "CLOUD_INTEGRATION_TESTS"

// Only run the integration tests on these specially preparred accounts
var integrationAwsAccounts = []string{"316817240772", "438881294876"}

// Cache the checks for credentials so it doesn't run for every test
var credCheck credChecker

type integrationTestContext struct {
	p         *AWSProvider
	log       *zap.Logger
	logBuffer *bytes.Buffer
	ctx       context.Context
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
	_, hasEnvVar := os.LookupEnv(testEnvVar)

	if hasEnvVar && !creds {
		t.Fatalf("cannot run integration tests without creds")
	}

	if !creds {
		t.Skip("no active creds for the integration testing account")
	}
}

func setupIntegrationProvider(t testing.TB, ctx *integrationTestContext) {
	t.Helper()

	var err error
	provider := &AWSProvider{}
	provider.logger = ctx.log
	provider.config, err = config.LoadDefaultConfig(ctx.ctx, func(lo *config.LoadOptions) error {
		lo.Region = testingutil.TestRegion
		return nil
	})
	if err != nil {
		t.Fatalf("cannot load config: %v", err)
	}

	provider.initClients()
	provider.mapper, err = mapper.New(embedConfig, ctx.log, reflect.ValueOf(provider))
	if err != nil {
		t.Fatalf("cannot instantiate mapper: %v", err)
	}

	ctx.p = provider
}

func setupIntegrationLogs(t testing.TB, ctx *integrationTestContext) {
	t.Helper()
	buf, ws := logBuffer()
	ctx.logBuffer = buf

	log, err := zap.NewDevelopment(zap.ErrorOutput(ws))
	if err != nil {
		t.Fatalf("cannot create zap logger: %v", err)
	}
	ctx.log = log
}

func logBuffer() (*bytes.Buffer, zapcore.WriteSyncer) {
	buf := &bytes.Buffer{}
	return buf, zapcore.Lock(zapcore.AddSync(buf))
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
