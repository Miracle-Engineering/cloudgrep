package types

import (
	"context"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
)

type FetchFunc func(context.Context, chan<- model.Resource) error
type Provider interface {
	//For a cloud provider this would be the account/project ID, set to empty if not relevant
	AccountId() string
	String() string
	FetchFunctions() map[string]FetchFunc
}
