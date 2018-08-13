package handler

import (
	"bytes"
	"fmt"
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/k8sclient"

	"github.com/muxiyun/Mae/pkg/errno"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func Terminal(ctx iris.Context) {

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	req := k8sclient.Restclient.Post().
		Resource("pods").
		Name("kube-test-deploy-3598112474-175fp").
		Namespace("kube-test").
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Container: "kube-test-ct",
		Command:   []string{"ls -a"},
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	fmt.Println("--------->",k8sclient.Config.WrapTransport)
	exec, err := remotecommand.NewSPDYExecutor(k8sclient.Config, "POST", req.URL())
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrInitExecutor, err), nil)
		return
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    false,
	})

	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrCannotExec, err), nil)
		return
	}

	if execErr.Len() > 0 {
		SendResponse(ctx, nil, fmt.Errorf("stderr: %v", execErr.String()))
		return
	}

	SendResponse(ctx, nil, execOut.String())
	return

}
