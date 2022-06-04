package main

import (
	"os"

	"github.com/run-x/cloudgrep/hack/awsgen/cmd"
)

func main() {
	cmd.Run(os.Args[1:])
}
