//　本来想直接隐藏pod的概念的，但是log查询以及web terminal中都需要
// 指定pod,所以这里提供查询pod的api

package handler

import (
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/k8sclient"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//直接返回pods *v1.PodList会返回太多无用信息，所以这里做了一层过滤，只选择其中有用的信息
//以减少带宽和响应时间。如需增加其他信息，修改PodMessage结构和SelectMessage函数即可
type PodMessage struct {
	PodName    string            `json:"pod_name"`
	Containers []string          `json:"containers"`
	PodStatus  v1.PodPhase       `json:"pod_status"`
	PodLabels  map[string]string `json:"pod_labels"`
}

// select useful message from pod list
func selectMessage(pods *v1.PodList) []PodMessage {
	var podmsgs []PodMessage
	for _, item := range pods.Items {
		var podmsg PodMessage
		var containers []string

		podmsg.PodName = item.Name
		podmsg.PodStatus = item.Status.Phase
		podmsg.PodLabels = item.Labels

		for _, container := range item.Spec.Containers {
			containers = append(containers, container.Name)
		}
		podmsg.Containers = containers

		podmsgs = append(podmsgs, podmsg)
	}
	return podmsgs
}

// get useful pod messages from specified namespace
func GetPod(ctx iris.Context) {
	ns := ctx.Params().Get("ns")
	current_user_role := ctx.Values().Get("current_user_role")

	// unadmin user is not allowed to access kube-system,kube-public,default namespace
	if current_user_role == "user" && (ns == "kube-system" || ns == "kube-public" || ns == "default") {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.WriteString("Forbidden")
		return
	}

	pods, err := k8sclient.ClientSet.CoreV1().Pods(ns).List(meta_v1.ListOptions{})
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrListPods, err), nil)
		return
	}

	SendResponse(ctx, nil, selectMessage(pods))
}
