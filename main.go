/*
Copyright © 2022 RunX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/run-x/cloudgrep/cmd"
	"github.com/run-x/cloudgrep/pkg/util"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			switch err := err.(type) {
			case util.UserError:
				fmt.Println(err.Error())
			default:
				panic(fmt.Errorf("Failed to start cloudgrep: %v", err))
			}
		}
	}()
	cmd.Execute()
}
