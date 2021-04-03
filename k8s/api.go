package k8s

import (
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

var Clientset *kubernetes.Clientset

func Init(kubeconfig *string) error {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}

	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}
