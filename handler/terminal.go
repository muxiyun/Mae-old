// xterm.js & ws

package handler

import (
	"bytes"
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/k8sclient"
	"log"

	"github.com/gorilla/websocket"
	"github.com/kataras/iris/core/errors"
	"github.com/muxiyun/Mae/pkg/errno"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetCommandOutput(ns, pod_name, container_name string, cmd []string) (string, error) {
	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	req := k8sclient.Restclient.Post().
		Resource("pods").
		Name(pod_name).
		Namespace(ns).
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Container: container_name,
		Command:   cmd,
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(k8sclient.Config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    false,
	})

	if err != nil {
		return "", err
	}

	if execErr.Len() > 0 {
		return "", errors.New(execErr.String())
	}

	return execOut.String(), nil
}

func Terminal(ctx iris.Context) {
	ns := ctx.Params().Get("ns")
	pod_name := ctx.Params().Get("pod_name")
	container_name := ctx.Params().Get("container_name")

	current_user_role := ctx.Values().Get("current_user_role")
	if current_user_role == "user" && (ns == "default" || ns == "kube-public" || ns == "kube-system") {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.WriteString("Forbidden")
		return
	}

	//get the websocket conn
	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)

	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrorGetWSConn, err), nil)
		return
	}

	// Interaction with client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// binary message will be ignore
		if messageType == websocket.TextMessage {
			var goodCmd []string
			for _, cmd := range strings.Split(string(p), " ") {
				c := strings.Trim(cmd, " ")
				goodCmd = append(goodCmd, c)
			}

			//exec the command and get the command output
			output, err := GetCommandOutput(ns, pod_name, container_name, goodCmd)
			if err != nil {
				SendResponse(ctx, errno.New(errno.ErrCannotExec, err), nil)
				return
			}
			//push the output to the client
			if err := conn.WriteMessage(messageType, []byte(output)); err != nil {
				SendResponse(ctx, errno.New(errno.ErrPush, err), nil)
				return
			}
		}
	}
	SendResponse(ctx, nil, "session over")
}
