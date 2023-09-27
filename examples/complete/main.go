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
	"strings"

	"github.com/maolonglong/go-editline"
)

var list = []string{"foo ", "bar ", "bsd ", "cli ", "ls ", "cd ", "malloc ", "tee "}

func main() {
	defer editline.Uninitialize()

	editline.SetCompleteFunc(func(token string) (s string, matched bool) {
		if token == "" {
			return "", false
		}
		var count int
		idx := -1
		for i, x := range list {
			if strings.HasPrefix(x, token) {
				count++
				idx = i
			}
		}
		if count == 1 {
			return list[idx][len(token):], true
		}
		return "", false
	})

	editline.SetListPossibFunc(func(token string) []string {
		var ss []string
		for _, x := range list {
			if strings.HasPrefix(x, token) {
				ss = append(ss, x)
			}
		}
		return ss
	})

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
}
