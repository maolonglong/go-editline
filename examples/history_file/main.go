// Copyright 2023 Shaolong Chen <shaolong.chen@outlook.it>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"

	"github.com/maolonglong/go-editline"
)

const _historyFilename = ".cli-history"

func main() {
	defer editline.Uninitialize()

	_ = editline.ReadHistroy(_historyFilename)

	for {
		line, err := editline.ReadLine("> ")
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		fmt.Println(line)
	}

	_ = editline.WriteHistroy(_historyFilename)
}
