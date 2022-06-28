package types

import (
	"context"

	"github.com/run-x/cloudgrep/pkg/model"
)

type FetchFunc func(context.Context, chan<- model.Resource) error
type Provider interface {
	//Id identifies the current cloud provider. For AWS, it returns the AWS Account ID.
	Id() string
	String() string
	FetchFunctions() map[string]FetchFunc
}
