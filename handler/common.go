package handler

import (
	"net/http"

	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/k8sclient"
	"errors"
)

func SendResponse(c iris.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.StatusCode(http.StatusOK)
	c.JSON(iris.Map{
		"code": code,
		"msg":  message,
		"data": data,
	})
}



func DeleteDeploymentAndServiceInCluster(version_config model.VersionConfig)error{
	//delete the deployment
	deploymentClient := k8sclient.ClientSet.ExtensionsV1beta1().
		Deployments(version_config.Deployment.NameSapce)
	if err := deploymentClient.Delete(version_config.Deployment.DeployName, nil);err != nil {
		return errors.New("error delete deployment, "+err.Error())
	}
	//delete the service
	ServiceClient := k8sclient.ClientSet.CoreV1().Services(version_config.Deployment.NameSapce)
	if err:=ServiceClient.Delete(version_config.Svc.SvcName, nil);err != nil {
		return errors.New("error delete service, "+err.Error())
	}

	return nil
}