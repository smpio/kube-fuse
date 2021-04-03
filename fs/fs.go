package fs

import "github.com/smpio/kube-fuse/k8s"

func ListDir(path string) []string {
	list, err := k8s.ListObjects("v1", "namespaces", "")

	if err != nil {
		return []string{}
	}

	names := make([]string, len(list.Items)+1)
	names[0] = "_"
	for idx, pod := range list.Items {
		names[idx+1] = pod.Name
	}

	return names
}
