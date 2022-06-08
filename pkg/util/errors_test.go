package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestEnableErrorStackTrace(t *testing.T) {
	EnableErrorStackTrace()
	require.True(t, showErrorStackTrace)
}

func TestGetStackTrace(t *testing.T) {
	t.Run("NoStackTrace", func(t *testing.T) {
		err := errors.New("hi")
		require.Equal(t, "", GetStackTrace(err))
	})
	t.Run("StackTrace", func(t *testing.T) {

		err := stackTraceError{
			Err:        errors.New("hi"),
			StackTrace: "blah",
		}
		require.Equal(t, "blah", GetStackTrace(err))
	})
	t.Run("WrappedStackTrace", func(t *testing.T) {

		err := fmt.Errorf("hi %w", stackTraceError{
			Err:        errors.New("hi"),
			StackTrace: "blah",
		})
		require.Equal(t, "blah", GetStackTrace(err))
	})
}

func TestAddStackStrace(t *testing.T) {
	var mu sync.Mutex
	t.Run("NoStackTrace", func(t *testing.T) {
		mu.Lock()
		defer mu.Unlock()
		showErrorStackTrace = false
		err := errors.New("hi")
		require.Equal(t, err, AddStackStrace(err))
	})
	t.Run("WithStackTrace", func(t *testing.T) {
		mu.Lock()
		defer mu.Unlock()
		showErrorStackTrace = true
		err := errors.New("hi")
		newerr := AddStackStrace(err)
		require.NotEqual(t, err, newerr)
		var sTE stackTraceError
		require.True(t, errors.As(newerr, &sTE))
		require.Equal(t, err, sTE.Err)
		require.NotEqual(t, "", sTE.StackTrace)
	})

}

func TestPrintStacktrace(t *testing.T) {

	t.Run("NoStackTrace", func(t *testing.T) {
		err := errors.New("hi")
		var b bytes.Buffer
		PrintStacktrace(err, &b)
		require.Equal(t, 0, b.Len())
	})
	t.Run("WithStackTrace", func(t *testing.T) {
		err := stackTraceError{
			Err:        errors.New("hi"),
			StackTrace: "blah",
		}
		var b bytes.Buffer
		PrintStacktrace(err, &b)
		require.Equal(t, "Stack Trace:\nblah\nhi", b.String())

	})
}
