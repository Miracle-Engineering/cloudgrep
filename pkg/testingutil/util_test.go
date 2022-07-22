package testingutil

import (
	"fmt"
	"testing"
)

type FakeTB struct {
	IsHelper bool
	IsFail   bool
	Logs     []string
}

var _ TestingTB = &FakeTB{}

func Fake(t testing.TB) *FakeTB {
	return &FakeTB{}
}

func (t *FakeTB) Helper() {
	t.IsHelper = true
}

func (t *FakeTB) Fatalf(format string, args ...any) {
	t.Errorf(format, args...)
	t.FailNow()
}

func (t *FakeTB) Errorf(format string, args ...any) {
	t.IsFail = true
	t.Logs = append(t.Logs, fmt.Sprintf(format, args...))
}

func (t *FakeTB) FailNow() {
	// We have no way of "aborting", since we don't want to actually fail the test.
	t.IsFail = true
}
