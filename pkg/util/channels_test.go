package util

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendAllFromSlice(t *testing.T) {
	vals := []string{"foo", "bar"}
	c := make(chan string, 2)
	done := make(chan struct{})

	var err error

	go func() {
		err = SendAllFromSlice(context.Background(), c, vals)
		close(done)
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("SendAllFromSlice timed out")
	case <-done:
	}

	assert.Nil(t, err, "expected no error")
	close(c)

	assert.Len(t, c, 2)

	var v string
	var ok bool

	v, ok = <-c
	assert.True(t, ok, "expected to receive value")
	assert.Equal(t, "foo", v)

	v, ok = <-c
	assert.True(t, ok, "expected to receive value")
	assert.Equal(t, "bar", v)

	_, ok = <-c
	assert.False(t, ok, "channel to be closed")
}

func TestSendAllFromSliceCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	vals := []string{"foo", "bar"}
	c := make(chan string)
	done := make(chan struct{})

	var err error

	go func() {
		err = SendAllFromSlice(ctx, c, vals)
		close(done)
	}()

	var v string
	var ok bool

	v, ok = <-c
	assert.True(t, ok, "expected to receive value")
	assert.Equal(t, "foo", v)

	cancel()
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("SendAllFromSlice took too long to exit after cancel")
	case <-done:
	}

	assert.ErrorIs(t, err, context.Canceled)
}
