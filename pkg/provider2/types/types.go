package types

import (
	"context"
	"github.com/run-x/cloudgrep/pkg/model"
)

type FetchFunc func(context.Context, chan<- model.Resource) error
type Provider interface {
	String() string
	FetchFunctions() map[string]FetchFunc
}
