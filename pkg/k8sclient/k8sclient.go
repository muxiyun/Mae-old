package k8sclient

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"github.com/spf13/viper"
	//"fmt"
)

var ClientSet *kubernetes.Clientset

func init() {
	//kubeconfig:=viper.GetString("kubeconfig")

	config, err := clientcmd.BuildConfigFromFlags("", "conf/admin.kubeconfig")
	if err != nil {
		panic(err.Error())
	}

	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

}
