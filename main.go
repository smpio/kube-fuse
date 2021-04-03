package main

// #cgo pkg-config: fuse
// #include "helpers.h"
// int c_main(int argc, char *argv[]);
import "C"
import (
	"context"
	"flag"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/util/homedir"
)

var api *API

func ListDir(path string) []string {
	pods, err := api.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []string{}
	}

	names := make([]string, len(pods.Items))
	for idx, pod := range pods.Items {
		names[idx] = pod.Name
	}

	return names
}

//export ReadDir
func ReadDir(path *C.char) **C.char {
	entries := ListDir(C.GoString(path))

	result := make([]*C.char, len(entries)+1)
	for idx, entry := range entries {
		result[idx] = C.CString(entry)
	}

	return C.copy_str_ptr_array(&result[0], C.ulong(len(entries)))
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var err error
	api, err = NewAPI(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	fuseArgs := []string{os.Args[0]}            // program name
	fuseArgs = append(fuseArgs, flag.Args()...) // positional params (or everything after --)
	fuseArgs = append(fuseArgs, "-f")           // always run in foreground
	os.Exit(fuseMain(fuseArgs))
}

func fuseMain(args []string) int {
	argc := C.int(len(args))
	argv := make([]*C.char, len(args))
	for idx, arg := range args {
		argv[idx] = C.CString(arg)
	}

	return int(C.c_main(argc, &argv[0]))
}
