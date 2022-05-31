package testingutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterFunc(t *testing.T) {
	in := []int{
		1,
		-3,
		0,
		5,
		-1,
	}

	expected := []int{
		1,
		0,
		5,
		-1,
	}

	actual := FilterFunc(in, func(v int) bool {
		return v > -3
	})

	assert.Equal(t, expected, actual)
}
