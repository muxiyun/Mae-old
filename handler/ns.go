package handler

import (
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/k8sclient"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNS(ctx iris.Context) {
	ns_name := ctx.Params().Get("ns")
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns_name}}
	_, err := k8sclient.ClientSet.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrCreateNamespace, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"message": "create ns " + ns_name})
}

func DeleteNS(ctx iris.Context) {
	ns_name := ctx.Params().Get("ns")
	err := k8sclient.ClientSet.CoreV1().Namespaces().
		Delete(ns_name, meta_v1.NewDeleteOptions(10))
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDeleteNamespace, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"message": "delete ns " + ns_name})
}



type Nsmsg struct {
	Name string `json:"name"`
	Status v1.NamespacePhase `json:"status"`
	CreateTime string `json:"create_time"`
}


func selectMsgFromNsList(nsList *v1.NamespaceList,current_user_role string)([]Nsmsg) {
	var nsMsgs []Nsmsg
	for _, item := range nsList.Items {
		// normal user can not see 'default','kube-system','kube-public' namespace
		if current_user_role == "user" && (item.Name == "default" || item.Name == "kube-system" || item.Name == "kube-public") {
			continue
		}
		var nsMsg Nsmsg
		nsMsg.Name=item.Name
		nsMsg.Status=item.Status.Phase
		nsMsg.CreateTime=item.CreationTimestamp.String()
		nsMsgs=append(nsMsgs,nsMsg)
	}
	return nsMsgs
}


func ListNS(ctx iris.Context) {
	ns, err := k8sclient.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrGetNamespace, err), nil)
		return
	}

	current_user_role:=ctx.Values().GetString("current_user_role")

	SendResponse(ctx, nil, selectMsgFromNsList(ns,current_user_role))
}
