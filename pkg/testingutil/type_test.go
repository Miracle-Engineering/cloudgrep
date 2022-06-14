package testingutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTypeStrType struct{}

func TestTypeStr(t *testing.T) {
	assert.Equal(t, "github.com/run-x/cloudgrep/pkg/testingutil.testTypeStrType", TypeStr(testTypeStrType{}))
}
