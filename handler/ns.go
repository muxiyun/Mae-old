package handler

import (
	//"github.com/muxiyun/Mae/model"
	"k8s.io/api/core/v1"
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/k8sclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/muxiyun/Mae/pkg/errno"
)

func CreateNS(ctx iris.Context){
	ns_name:= ctx.Params().Get("ns")
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name:ns_name}}
	_, err := k8sclient.ClientSet.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		SendResponse(ctx,errno.New(errno.ErrCreateNamespace,err),nil)
		return
	}
	SendResponse(ctx,nil,iris.Map{"message":"create ns "+ns_name})
}


func DeleteNS(ctx iris.Context){
	ns_name:= ctx.Params().Get("ns")
	 err := k8sclient.ClientSet.CoreV1().Namespaces().
	 	Delete(ns_name,meta_v1.NewDeleteOptions(10))
	if err != nil {
		SendResponse(ctx,errno.New(errno.ErrDeleteNamespace,err),nil)
		return
	}
	SendResponse(ctx,nil,iris.Map{"message":"delete ns "+ns_name})
}


func ListNS(ctx iris.Context){
	ns, err := k8sclient.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err!=nil{
		SendResponse(ctx,errno.New(errno.ErrGetNamespace,err),nil)
		return
	}

	//removee default kube-public kube-system if not a admin
	if ctx.Values().GetString("current_user_role")!="admin"{
		ns.Items=ns.Items[3:]
	}

	 SendResponse(ctx,nil,ns)
}



