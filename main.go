package main

// #cgo pkg-config: fuse
// int c_main(int argc, char *argv[]);
import "C"
import (
	"flag"
	"os"
	"path/filepath"

	"github.com/smpio/kube-fuse/fs"
	_ "github.com/smpio/kube-fuse/fs"
	"github.com/smpio/kube-fuse/k8s"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var err error
	err = k8s.Init(kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	err = fs.InitCache()
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
