package main

import (
	"fmt"
	"os"

	"github.com/run-x/cloudgrep/pkg/cli"
	"github.com/run-x/cloudgrep/pkg/util"
)

func main() {

	eventProperties := map[string]interface{}{
		"app": "cloudgrep",
	}

	util.SendEvent(util.BASE_EVENT, eventProperties, nil)

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
