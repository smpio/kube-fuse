package main

// #cgo pkg-config: fuse
// #include "helpers.h"
// int main2(int argc, char *argv[]);
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
	//return []string{}

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
	args := flag.Args()

	var err error
	api, err = NewAPI(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	api.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	argc := C.int(len(args) + 1)
	argv := make([]*C.char, len(args)+1)
	argv[0] = C.CString(os.Args[0])
	for idx, arg := range args {
		argv[idx+1] = C.CString(arg)
	}

	code := int(C.main2(argc, &argv[0]))
	os.Exit(code)
}
