package k8sclient

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)




var (
	ClientSet *kubernetes.Clientset
	Restclient  *rest.RESTClient
	Config *rest.Config
)

//var Cluster K8sCluster
//
//type K8sCluster struct{
//	ClientSet *kubernetes.Clientset
//	Restclient  *rest.RESTClient
//	Config *rest.Config
//}

func init() {
	//kubeconfig:=viper.GetString("kubeconfig")

 	var err error
	// get config
	Config, err := clientcmd.BuildConfigFromFlags("", "conf/admin.kubeconfig")
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
		Group:   "k8s.io",
		Version: "v1",
	}
	Config.GroupVersion = &groupversion
	Config.APIPath = "/apis"
	Config.ContentType = runtime.ContentTypeJSON
	Config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	Config.WrapTransport=nil

	//get restclient
	Restclient, err = rest.RESTClientFor(Config)
	if err!=nil{
		panic(err.Error())
	}
	fmt.Println("get restclient success")

}
