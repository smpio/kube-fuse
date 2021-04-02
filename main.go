package main

// #cgo pkg-config: fuse
// #include "helpers.h"
// int main2(int argc, char *argv[]);
import "C"
import (
	"os"
)

func ListDir(path string) []string {
	return []string{"hello", "world"}
}

//export ReadDir
func ReadDir(path *C.char) **C.char {
	entries := ListDir(C.GoString(path))

	result := make([]*C.char, len(entries))
	for idx, entry := range entries {
		result[idx] = C.CString(entry)
	}

	return C.copy_str_ptr_array(&result[0], C.ulong(len(entries)))
}

func main() {
	argc := C.int(len(os.Args))
	argv := make([]*C.char, len(os.Args))
	for idx, arg := range os.Args {
		argv[idx] = C.CString(arg)
	}

	code := int(C.main2(argc, &argv[0]))
	os.Exit(code)
}
