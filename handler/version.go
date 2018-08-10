package handler

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/kataras/iris/core/errors"
	"github.com/muxiyun/Mae/pkg/k8sclient"

	apiv1 "k8s.io/api/core/v1"
	appsv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)


// the CreateVersion handler just to create a version config and store it to database
// and it doesn't create any resources in the cluster, you can create resources by applying
// one version config that created before. That is to send a request to the ApplyVersion handler
func CreateVersion(ctx iris.Context) {
	var (
		rv model.ReqVersion
		v model.Version
	)

	ctx.ReadJSON(&rv)
	v.ServiceID=rv.ServiceID
	v.VersionName=rv.VersionName
	v.VersionDesc=rv.VersionDesc

	r,_:=json.Marshal(rv)
	fmt.Println(string(r))

	//将配置部分序列化，然后存到VersionConfig字段中
	configbytearray,err:=json.Marshal(rv.VersionConf)
	if err!=nil{
		SendResponse(ctx,errno.New(errno.ErrVersionConfigMarshal,err),nil)
		return
	}
	v.VersionConfig=string(configbytearray)


	if err:=v.Create();err!=nil{
		SendResponse(ctx,errno.New(errno.ErrCreateVersion,err),nil)
	}
	SendResponse(ctx,nil,nil)
}


// the ApplyVersion handler to apply a version config that created before.
// this handler will replace a version with the specified one,
// that is to remove the resources of the previous version config and then
// create the resources about the specified version config
func ApplyVersion(ctx iris.Context){
	version_name:=ctx.URLParam("version_name")
	if version_name==""{
		SendResponse(ctx,errno.New(errno.ErrVersionNameEmpty,errors.New("")),nil)
		return
	}

	v,err:=model.GetVersionByName(version_name)
	if err!=nil{
		SendResponse(ctx,errno.New(errno.ErrDatabase,err),nil)
		return
	}

	var version_config model.VersionConfig
	json.Unmarshal([]byte(v.VersionConfig),&version_config)

	// get the deployment client for the specified namespace
	deploymentClient:=k8sclient.ClientSet.ExtensionsV1beta1().
		Deployments(version_config.Deployment.NameSapce)

	// config a deployment
	deployment:=&appsv1b1.Deployment{
		ObjectMeta:metav1.ObjectMeta{
			Name:version_config.Deployment.DeployName,
		},
		Spec:appsv1b1.DeploymentSpec{
			Replicas:int32Ptr(int32(version_config.Deployment.Replicas)),
			Template:apiv1.PodTemplateSpec{
				ObjectMeta:metav1.ObjectMeta{
					Labels:version_config.Deployment.Labels,
				},
				Spec:apiv1.PodSpec{
					Volumes:bindVolumeSource(version_config),
					Containers:bindContainers(version_config),
						},
				},
			},
		}

	// create the deployment
	result,err:=deploymentClient.Create(deployment)
	if err != nil {
		SendResponse(ctx,errno.New(errno.ErrCreateDeployment,err),nil)
		return
	}

	// update the service's current_version field
	service,err:=model.GetServiceByID(int64(v.ServiceID))
	if err!=nil{
		SendResponse(ctx,errno.New(errno.ErrDatabase,err),nil)
		return
	}
	service.CurrentVersion=v.VersionName

	if err=service.Update();err!=nil{
		SendResponse(ctx,errno.New(errno.ErrDatabase,err),nil)
		return
	}

	// change the version's active field to true
	v.Active=true
	if err=v.Update();err!=nil{
		SendResponse(ctx,errno.New(errno.ErrDatabase,err),nil)
		return
	}

	SendResponse(ctx,nil,"Created deployment"+result.GetObjectMeta().GetName())
}

// the UnapplyVersion handler just to STOP a being used Version config,
// that is to delete the resources　of the version config in the cluster,
// but the version config information is still exist in the database, so
// you can apply it to create the resources again.
func UnapplyVersion(ctx iris.Context){

}

// the GetVersion handler is to get version config information
func GetVersion(ctx iris.Context) {

}

// the UpdateVersion handler
func UpdateVersion(ctx iris.Context) {

}

// the DeleteVersion handler will delete the version config information first,
// if the version config that requested to be deleted is being used at that time,
// then the resources about the version config will be removed in the cluster too.
func DeleteVersion(ctx iris.Context) {

}

//the GetVersionList handler can get all the Version config information in the database
// or just to get the version config that belongs to one service.
func GetVersionList(ctx iris.Context) {

}



func int32Ptr(i int32) *int32 { return &i }