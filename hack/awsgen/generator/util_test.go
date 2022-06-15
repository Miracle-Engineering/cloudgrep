package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinenumbers(t *testing.T) {
	input := "a\nb\nc\n"
	expectedParts := []string{
		"/* 1 */ a",
		"/* 2 */ b",
		"/* 3 */ c",
		"/* 4 */ ",
	}
	expected := strings.Join(expectedParts, "\n") + "\n"

	actual := linenumbers(input)
	assert.Equal(t, expected, actual)
}

func TestAwsServicePackage(t *testing.T) {
	expected := "github.com/aws/aws-sdk-go-v2/service/foo/bar/spam"
	actual := awsServicePackage("foo", "bar", "spam")

	assert.Equal(t, expected, actual)
}
