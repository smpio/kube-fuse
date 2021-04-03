package k8s

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ClientConfig *rest.Config
	Clientset    *kubernetes.Clientset

	PreferredResources           []*metav1.APIResourceList
	PreferredNamespacedResources []*metav1.APIResourceList
)

type KubeObject struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

type KubeObjectList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// A list of objects.
	Items []KubeObject `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func (in *KubeObjectList) DeepCopyObject() runtime.Object {
	return nil
}

func Init(kubeconfig *string) error {
	var err error

	ClientConfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}

	Clientset, err = kubernetes.NewForConfig(ClientConfig)
	if err != nil {
		return err
	}

	PreferredResources, err = Clientset.ServerPreferredResources()
	if err != nil {
		return err
	}

	PreferredNamespacedResources, err = Clientset.ServerPreferredResources()
	if err != nil {
		return err
	}

	return nil
}

func ListObjects(groupVersion string, resourceName string, namespace string) (*KubeObjectList, error) {
	gv := strings.SplitN(groupVersion, "/", 2)

	config := *ClientConfig

	if len(gv) < 2 {
		config.GroupVersion = &schema.GroupVersion{
			Group:   "",
			Version: gv[0],
		}
		config.APIPath = "/api"
	} else {
		config.GroupVersion = &schema.GroupVersion{
			Group:   gv[0],
			Version: gv[1],
		}
		config.APIPath = "/apis"
	}

	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	result := KubeObjectList{}

	request := client.Get()
	if namespace != "" {
		request = request.Namespace(namespace)
	}

	err = request.
		Resource(resourceName).
		//VersionedParams(&opts, scheme.ParameterCodec).
		//Timeout(timeout).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
