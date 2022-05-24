package script

// Make sure these modules aren't removed from go.mod by `go mod tidy`
import (
	_ "github.com/stretchr/testify"

	// Using pkg/config so linting doesn't complain about the `main` package
	_ "github.com/vektra/mockery/v2/pkg/config"
)
