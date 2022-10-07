package main

import (
	"os"

	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/cmd"
)

func main() {
	cmd.Run(os.Args[1:])
}
