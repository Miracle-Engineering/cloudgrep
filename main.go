package main

import (
	"fmt"

	"github.com/run-x/cloudgrep/pkg/cli"
)

func main() {

	if err := cli.Run(); err != nil {
		fmt.Println(err)
	}
}
