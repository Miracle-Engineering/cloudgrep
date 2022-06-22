package regions

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/stretchr/testify/require"
)

func mustRegion(t testing.TB, raw string) Region {
	region, err := regionForRaw(raw)
	require.NoError(t, err)
	return region
}

func regionIds(regions []Region) []string {
	return testingutil.SliceConvertFunc(regions, func(region Region) string {
		return region.ID()
	})
}

type mockedHttpClient struct {
	contentType string
	body        []byte
}

var _ aws.HTTPClient = &mockedHttpClient{}

func (c *mockedHttpClient) Do(req *http.Request) (*http.Response, error) {
	body := bytes.NewReader(c.body)

	headers := make(http.Header)

	if c.contentType != "" {
		headers[http.CanonicalHeaderKey("Content-Type")] = []string{c.contentType}
	}

	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,

		Header:        headers,
		Body:          io.NopCloser(body),
		ContentLength: int64(body.Len()),
		Close:         true,
		Request:       req,
	}

	return resp, nil
}

type mockedCredentialsProvider struct{}

var _ aws.CredentialsProvider = mockedCredentialsProvider{}

func (p mockedCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "foo",
		SecretAccessKey: "bar",
	}, nil
}
