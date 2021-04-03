package main

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type API struct {
	Clientset       *kubernetes.Clientset
	DiscoveryClient *discovery.DiscoveryClient
}

func NewAPI(kubeconfig *string) (*API, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	return &API{
		Clientset:       clientset,
		DiscoveryClient: discoveryClient,
	}, nil
}
