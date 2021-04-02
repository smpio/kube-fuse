package main

// #cgo pkg-config: fuse
// int main2(int argc, char *argv[]);
import "C"
import (
	"os"
)

func main() {
	argc := C.int(len(os.Args))
	argv := make([]*C.char, len(os.Args))
	for idx, arg := range os.Args {
		argv[idx] = C.CString(arg)
	}

	code := int(C.main2(argc, &argv[0]))
	os.Exit(code)
}
