package testingutil

import (
	"strconv"
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

func TestUnique(t *testing.T) {
	in := []int{
		1, 0, 0, 1, 2, 0, 5, -1, -3, 2,
	}

	expected := []int{
		1, 0, 2, 5, -1, -3,
	}

	actual := Unique(in)
	assert.Equal(t, expected, actual)
}

func TestSliceConvertFunc(t *testing.T) {
	in := []int{
		0, 1, 3, -1,
	}

	expected := []string{
		"0", "1", "3", "-1",
	}

	actual := SliceConvertFunc(in, strconv.Itoa)
	assert.Equal(t, expected, actual)
}
