package fs

import (
	"context"

	"github.com/smpio/kube-fuse/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListDir(path string) []string {
	pods, err := k8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []string{}
	}

	names := make([]string, len(pods.Items))
	for idx, pod := range pods.Items {
		names[idx] = pod.Name
	}

	return names
}
