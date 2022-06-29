package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunks(t *testing.T) {
	testCases := []struct {
		input     []string
		chunkSize int
		output    [][]string
	}{
		{[]string{}, 3, [][]string{}},
		{[]string{"1", "2", "3", "4", "5", "6"}, 3, [][]string{{"1", "2", "3"}, {"4", "5", "6"}}},
		{[]string{"1", "2", "3", "4", "5", "6"}, 9, [][]string{{"1", "2", "3", "4", "5", "6"}}},
		{[]string{"1", "2", "3", "4", "5", "6"}, 6, [][]string{{"1", "2", "3", "4", "5", "6"}}},
		{[]string{"1", "2", "3", "4", "5", "6"}, 5, [][]string{{"1", "2", "3", "4", "5"}, {"6"}}},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.output, Chunks(tc.input, tc.chunkSize))
	}
}
