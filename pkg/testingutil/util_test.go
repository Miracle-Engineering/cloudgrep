package testingutil

import (
	"fmt"
	"testing"
)

type FakeTB struct {
	testing.TB
	IsHelper bool
	IsFail   bool
	Logs     []string
}

func Fake(t testing.TB) *FakeTB {
	return &FakeTB{
		TB: t,
	}
}

func (t *FakeTB) Helper() {
	t.IsHelper = true
}

func (t *FakeTB) Fatal(args ...any) {
	panic(fmt.Sprint(args...))
}

func (t *FakeTB) Fatalf(format string, args ...any) {
	panic(fmt.Errorf(format, args...))
}

func (t *FakeTB) Errorf(format string, args ...any) {
	t.IsFail = true
	t.Logs = append(t.Logs, fmt.Sprintf(format, args...))
}
