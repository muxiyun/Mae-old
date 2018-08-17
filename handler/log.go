package handler

import (
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/k8sclient"
	"k8s.io/api/core/v1"
)

//get log from a container, this handler need namespace, pod name and container name
func GetLog(ctx iris.Context) {
	ns := ctx.Params().Get("ns")
	pod_name := ctx.Params().Get("pod_name")
	container_name := ctx.Params().Get("container_name")

	current_user_role := ctx.Values().Get("current_user_role")
	if current_user_role == "user" && (ns == "default" || ns == "kube-public" || ns == "kube-system") {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.WriteString("Forbidden")
		return
	}

	// get the log query request
	restclientRequest := k8sclient.ClientSet.CoreV1().Pods(ns).
		GetLogs(pod_name, &v1.PodLogOptions{Container: container_name})

	// do the request and get the result
	result, _ := restclientRequest.Do().Raw()

	SendResponse(ctx, nil, string(result))
}
