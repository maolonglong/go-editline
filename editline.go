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

// Package editline provides Go bindings for [editline].
//
// [editline]: https://github.com/troglobit/editline
package editline

/*
#cgo pkg-config: libeditline

#include <stdio.h>
#include <stdlib.h>

#include <editline.h>

extern char* _complete_func(char* token, int* match);
extern int _list_possib_func(char* token, char*** av);

static void _set_complete_func(void) { rl_set_complete_func(_complete_func); }
static void _set_list_possib_func(void) { rl_set_list_possib_func(_list_possib_func); }
*/
import "C"

import (
	"io"
	"runtime"
	"unsafe"
)

var (
	globalCompleteFunc CompleteFunc = func(_ string) (string, bool) {
		return "", false
	}
	globalListPossibFunc ListPossibFunc = func(_ string) []string {
		return nil
	}
)

type (
	// CompleteFunc equals rl_complete_func_t.
	CompleteFunc func(token string) (s string, matched bool)

	// ListPossibFunc equals rl_list_possib_func_t.
	ListPossibFunc func(token string) []string
)

func init() {
	C._set_complete_func()
	C._set_list_possib_func()
}

// SetNoEcho sets whether or not to echo the input in the terminal.
func SetNoEcho(b bool) {
	if b {
		C.el_no_echo = C.int(1)
	} else {
		C.el_no_echo = C.int(0)
	}
}

// SetNoHist sets whether or not to enable auto-save and access to history.
func SetNoHist(b bool) {
	if b {
		C.el_no_hist = C.int(1)
	} else {
		C.el_no_hist = C.int(0)
	}
}

// SetHistSize sets the maximum size of the history.
func SetHistSize(i int) {
	C.el_hist_size = C.int(i)
}

// Initialize initializes the editline library.
//
// There is no need to call it manually as it will be initialized automatically.
func Initialize() {
	C.rl_initialize()
}

// Uninitialize frees all internal memory.
func Uninitialize() {
	C.rl_uninitialize()
}

// ReadLine displays the given prompt on stdout, waits for user input on stdin
// and then returns a line of text with the trailing newline removed.
//
// Each line returned is automatically saved in the internal history list,
// unless it happens to be equal to the previous line.
func ReadLine(prompt string) (string, error) {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	cline := C.readline(cprompt)
	if cline == nil {
		return "", io.EOF
	}
	defer C.free(unsafe.Pointer(cline))
	return C.GoString(cline), nil
}

// AddHistroy adds a line to the editline history.
func AddHistroy(line string) {
	cline := C.CString(line)
	defer C.free(unsafe.Pointer(cline))
	C.add_history(cline)
}

// ReadHistroy reads the editline history from a file.
func ReadHistroy(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_, err := C.read_history(cfilename)
	return err
}

// WriteHistroy writes the editline history to a file.
func WriteHistroy(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_, err := C.write_history(cfilename)
	return err
}

// SetCompleteFunc sets the complete function that will be called during tab completion.
func SetCompleteFunc(f CompleteFunc) {
	globalCompleteFunc = f
}

// SetListPossibFunc sets the list possibilities function that will be called during tab completion.
func SetListPossibFunc(f ListPossibFunc) {
	globalListPossibFunc = f
}

//export _complete_func
func _complete_func(token *C.char, match *C.int) *C.char {
	s, matched := globalCompleteFunc(C.GoString(token))
	if matched {
		*match = C.int(1)
	}
	return C.CString(s)
}

//export _list_possib_func
func _list_possib_func(token *C.char, av ***C.char) C.int {
	ss := globalListPossibFunc(C.GoString(token))
	if len(ss) == 0 {
		return C.int(0)
	}
	*av = stringSlice2cstrArray(ss)
	return C.int(len(ss))
}

func stringSlice2cstrArray(ss []string) **C.char {
	n := len(ss)
	ptr := (**C.char)(C.malloc(C.size_t(n * int(unsafe.Sizeof((*C.char)(nil))))))
	a := ([]*C.char)(unsafe.Slice(ptr, n))
	for i := 0; i < n; i++ {
		a[i] = C.CString(ss[i])
	}
	return ptr
}
