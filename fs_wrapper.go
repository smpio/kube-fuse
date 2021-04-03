package main

// #include "helpers.h"
import "C"
import "github.com/smpio/kube-fuse/fs"

//export ReadDir
func ReadDir(path *C.char) **C.char {
	entries := fs.ListDir(C.GoString(path))

	result := make([]*C.char, len(entries)+1)
	for idx, entry := range entries {
		result[idx] = C.CString(entry)
	}

	return C.copy_str_ptr_array(&result[0], C.ulong(len(entries)))
}
