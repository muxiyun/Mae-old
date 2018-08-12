package handler

import (
	"fmt"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/k8sclient"
)

// the CreateVersion handler just to create a version config and store it to database
// and it doesn't create any resources in the cluster, you can create resources by applying
// one version config that created before. That is to send a request to the ApplyVersion handler
func CreateVersion(ctx iris.Context) {
	var (
		rv model.ReqVersion
		v  model.Version
	)

	ctx.ReadJSON(&rv)
	v.ServiceID = rv.ServiceID
	v.VersionName = rv.VersionName
	v.VersionDesc = rv.VersionDesc

	r, _ := json.Marshal(rv)
	fmt.Println(string(r))

	//将配置部分序列化，然后存到VersionConfig字段中
	configbytearray, err := json.Marshal(rv.VersionConf)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrVersionConfigMarshal, err), nil)
		return
	}
	v.VersionConfig = string(configbytearray)

	if err := v.Create(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrCreateVersion, err), nil)
	}
	SendResponse(ctx, nil, nil)
}

// the ApplyVersion handler to apply a version config that created before.
// this handler will replace a version with the specified one,
// that is to remove the resources of the previous version config and then
// create the resources about the specified version config
func ApplyVersion(ctx iris.Context) {
	version_name := ctx.URLParam("version_name")
	if version_name == "" {
		SendResponse(ctx, errno.New(errno.ErrVersionNameEmpty, errors.New("")), nil)
		return
	}

	v, err := model.GetVersionByName(version_name)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	// get service that the version belongs to
	service, err := model.GetServiceByID(int64(v.ServiceID))
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	//service.CurrentVersion=="" 发行某一service的第一版
	// service.CurrentVersion==v.VersionName 重复应用
	//  (service.CurrentVersion!="") &&(service.CurrentVersion!=v.VersionName) 切换版本

	//重复apply,什么都不做，直接返回
	if service.CurrentVersion==v.VersionName {
		SendResponse(ctx,nil,nil)
		return
	}

	//切换版本，需要首先删除旧版本的deployment,service
	if (service.CurrentVersion!="") &&(service.CurrentVersion!=v.VersionName){
		old_version, err := model.GetVersionByName(service.CurrentVersion)
		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
			return
		}

		//unmarshal the old version config
		var old_version_config model.VersionConfig
		json.Unmarshal([]byte(old_version.VersionConfig), &old_version_config)

		//get the deploymentClient and delete the old deployment
		deploymentClient := k8sclient.ClientSet.ExtensionsV1beta1().
			Deployments(old_version_config.Deployment.NameSapce)
		err=deploymentClient.Delete(old_version_config.Deployment.DeployName,nil)
		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrDeleteDeployment, err), nil)
			return
		}

		// get the serviceClient and delete the old service
		ServiceClient := k8sclient.ClientSet.CoreV1().Services(old_version_config.Deployment.NameSapce)
		err=ServiceClient.Delete(old_version_config.Svc.SvcName, nil)
		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrDeleteService, err), nil)
			return
		}

		// turn old_version's active field to false
		old_version.Active=false
		if err = old_version.Update(); err != nil {
			SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
			return
		}
	}


	//部署新版本
	var version_config model.VersionConfig
	json.Unmarshal([]byte(v.VersionConfig), &version_config)

	// get the deployment client for the specified namespace
	deploymentClient := k8sclient.ClientSet.ExtensionsV1beta1().
		Deployments(version_config.Deployment.NameSapce)

	// config the deployment with the version_config struct
	deployment := ConfigDeployment(version_config)

	// create the deployment
	deploy_result, err := deploymentClient.Create(deployment)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrCreateDeployment, err), nil)
		return
	}

	// get the k8s service client of a specified namespace
	ServiceClient := k8sclient.ClientSet.CoreV1().Services(version_config.Deployment.NameSapce)

	// config the ks8 service with the version_config struct
	svc := ConfigService(version_config)

	// create the k8s service
	svc_result, err := ServiceClient.Create(svc)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrCreateService, err), nil)
		return
	}

	// update the mae service's current_version field

	service.CurrentVersion = v.VersionName

	if err = service.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	// change the version's active field to true
	v.Active = true
	if err = v.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	SendResponse(ctx, nil, "Created"+deploy_result.GetObjectMeta().GetName()+
		" and "+svc_result.GetObjectMeta().GetName())
}

// the UnapplyVersion handler just to STOP a being used Version config,
// that is to delete the resources　of the version config in the cluster,
// but the version config information is still exist in the database, so
// you can apply it to create the resources again.
func UnapplyVersion(ctx iris.Context) {
	//get the version_name url param
	version_name := ctx.URLParam("version_name")
	if version_name == "" {
		SendResponse(ctx, errno.New(errno.ErrVersionNameEmpty, errors.New("")), nil)
		return
	}

	//get version information from database by version_name field
	v, err := model.GetVersionByName(version_name)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	// unmarshal the VersionConfig field to a version_config struct
	var version_config model.VersionConfig
	json.Unmarshal([]byte(v.VersionConfig), &version_config)

	// get the deployment client for the specified namespace
	deploymentClient := k8sclient.ClientSet.ExtensionsV1beta1().
		Deployments(version_config.Deployment.NameSapce)
	err = deploymentClient.Delete(version_config.Deployment.DeployName, nil)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDeleteDeployment, err), nil)
		return
	}

	ServiceClient := k8sclient.ClientSet.CoreV1().Services(version_config.Deployment.NameSapce)
	ServiceClient.Delete(version_config.Svc.SvcName, nil)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDeleteService, err), nil)
		return
	}

	//get the service of mae
	service, err := model.GetServiceByID(int64(v.ServiceID))
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	// now there is no current version
	service.CurrentVersion = ""

	if err = service.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	// change the version's active field to false
	v.Active = false
	if err = v.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	SendResponse(ctx, nil, "remove the deployment and service")
}

// the GetVersion handler is to get detail information of a version
func GetVersion(ctx iris.Context) {
	version_name := ctx.Params().Get("version_name")

	v, err := model.GetVersionByName(version_name)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	SendResponse(ctx, nil, v)
}

func Scale(ctx iris.Context) {

}

// the UpdateVersion handler
func UpdateVersion(ctx iris.Context) {

}

// the DeleteVersion handler will delete the version config information first,
// if the version config that requested to be deleted is being used at that time,
// then the resources about the version config will be removed in the cluster too.
func DeleteVersion(ctx iris.Context) {
	version_id, _ := ctx.Params().GetInt64("id")

	version, err := model.GetVersionByID(version_id)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	var version_config model.VersionConfig
	json.Unmarshal([]byte(version.VersionConfig), &version_config)

	//if the version is active ,we need to delete
	// the service and the database record together,
	// else just to delete the database record.
	if version.Active == true {
		deploymentClient := k8sclient.ClientSet.ExtensionsV1beta1().
			Deployments(version_config.Deployment.NameSapce)
		err = deploymentClient.Delete(version_config.Deployment.DeployName, nil)
		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrDeleteDeployment, err), nil)
			return
		}

		ServiceClient := k8sclient.ClientSet.CoreV1().Services(version_config.Deployment.NameSapce)
		ServiceClient.Delete(version_config.Svc.SvcName, nil)
		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrDeleteService, err), nil)
			return
		}

		//get the service of mae
		service, err := model.GetServiceByID(int64(version.ServiceID))
		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
			return
		}
		// now there is no current version
		service.CurrentVersion = ""

		if err = service.Update(); err != nil {
			SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
			return
		}

	}

	//delete database record
	err = model.DeleteVersion(uint(version_id))
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	SendResponse(ctx, nil, nil)
}

//the GetVersionList handler can get all the Version config information in the database
// or just to get the version config that belongs to one service.
func GetVersionList(ctx iris.Context) {
	limit := ctx.URLParamIntDefault("limit", 20)    //how many if limit=0,default=20
	offsize := ctx.URLParamIntDefault("offsize", 0) // from where
	service_id := ctx.URLParamIntDefault("service_id", 0)

	var (
		versions []*model.Version
		count    uint64
		err      error
	)

	if service_id == 0 { //list all, admin only
		if ctx.Values().GetString("current_user_role") == "admin" {
			versions, count, err = model.ListVersion(offsize, limit)
		} else {
			ctx.StatusCode(iris.StatusForbidden)
			return
		}
	} else { // list versions belongs to an service,login
		versions, count, err = model.ListVersionByServiceID(offsize, limit, uint(service_id))
	}

	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"count": count, "versions": versions})
}


func int32Ptr(i int32) *int32 { return &i }
