package k8sclient

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Config     *rest.Config
	ClientSet  *kubernetes.Clientset
	Restclient *rest.RESTClient
)

func init() {

	var err error
	Config, err = clientcmd.BuildConfigFromFlags("", "conf/admin.kubeconfig")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("k8s client config success")

	//get clientset
	ClientSet, err = kubernetes.NewForConfig(Config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("get clientset success")

	groupversion := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	Config.GroupVersion = &groupversion
	Config.APIPath = "/api"
	Config.ContentType = runtime.ContentTypeJSON
	Config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	//get restclient
	Restclient, err = rest.RESTClientFor(Config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("get restclient success")

}
