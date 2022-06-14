package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRa() RecursiveAppend[string] {
	return RecursiveAppend[string]{
		Keys: []string{"foo", "bar", "spam", "ham"},
		Root: "apple",
	}
}

func TestRecursiveAppend_IsLast(t *testing.T) {
	var err error
	ra := makeRa()

	for i := 0; i < 3; i++ {
		assert.False(t, ra.IsLast())
		ra, err = ra.Next()
		require.NoError(t, err)
	}

	assert.True(t, ra.IsLast())
}

func TestRecursiveAppend_IterVar(t *testing.T) {
	var err error
	ra := makeRa()

	assert.Equal(t, "apple", ra.IterVar())

	ra, err = ra.Next()
	require.NoError(t, err)
	assert.Equal(t, "item_0", ra.IterVar())

	ra, err = ra.Next()
	require.NoError(t, err)
	assert.Equal(t, "item_1", ra.IterVar())
}

func TestRecursiveAppend_Current(t *testing.T) {
	var err error
	ra := makeRa()

	assert.Equal(t, "foo", ra.Current())

	ra, err = ra.Next()
	require.NoError(t, err)
	assert.Equal(t, "bar", ra.Current())
}

func TestRecursiveAppend_NextIterVar(t *testing.T) {
	var err error
	ra := makeRa()

	assert.Equal(t, "item_0", ra.NextIterVar())

	ra, err = ra.Next()
	require.NoError(t, err)
	assert.Equal(t, "item_1", ra.NextIterVar())
}

func TestRecursiveAppend_Next(t *testing.T) {
	ra := makeRa()

	require.Equal(t, 0, ra.Idx)

	// Make sure that Next does not mutate current
	newRa, err := ra.Next()
	require.NoError(t, err)
	assert.Equal(t, 1, newRa.Idx)
	assert.Equal(t, 0, ra.Idx)

	for i := 0; i < 2; i++ {
		newRa, err = newRa.Next()
		assert.NoError(t, err)
	}

	newerRa, err := newRa.Next()
	assert.Error(t, err)
	assert.Equal(t, newRa.Idx, newerRa.Idx)
}

func TestRecursiveAppend_SetData(t *testing.T) {
	ra := makeRa()

	ra2, err := ra.Next()
	require.NoError(t, err)

	// SetData only propagates up if parents already have a data
	err = ra2.SetData("foo", "bar")
	assert.NoError(t, err)

	assert.Contains(t, ra2.Data, "foo")
	assert.Equal(t, "bar", ra2.Data["foo"])
	assert.Empty(t, ra.Data)

	ra3, err := ra2.Next()
	require.NoError(t, err)

	err = ra3.SetData("spam", "ham")
	assert.NoError(t, err)

	assert.Contains(t, ra2.Data, "spam")
	assert.Equal(t, "ham", ra2.Data["spam"])
}

func TestRecursiveAppend_WithRoot(t *testing.T) {
	ra := makeRa()
	ra2 := ra.WithRoot("foo")
	assert.Equal(t, "apple", ra.Root)
	assert.Equal(t, "foo", ra2.Root)
}
