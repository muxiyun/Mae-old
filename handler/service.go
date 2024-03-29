package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/mail"
	"time"
)

// create a service
func CreateService(ctx iris.Context) {
	var svc model.Service
	ctx.ReadJSON(&svc)

	if svc.AppID == 0 || svc.SvcName == "" {
		SendResponse(ctx, errno.New(errno.ServiceNameEmptyorAppIDTypeError, errors.New("")), nil)
		return
	}

	if err := svc.Create(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrCreateService, err), nil)
	}

	SendResponse(ctx, nil, iris.Map{"id": svc.ID})
}

// get a service info by svc_name
func GetService(ctx iris.Context) {
	svc_name := ctx.Params().Get("svc_name")
	fmt.Println(svc_name)
	svc, err := model.GetServiceByName(svc_name)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrGetService, err), nil)
		return
	}
	SendResponse(ctx, nil, svc)
}

// update app_id or/and svc_name or/and svc_desc
func UpdateService(ctx iris.Context) {
	var newsvc model.Service
	ctx.ReadJSON(&newsvc)

	id, _ := ctx.Params().GetInt64("id")
	svc, err := model.GetServiceByID(id)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	//update the app_id of a service(move a service to another app)
	if newsvc.AppID != 0 {
		svc.AppID = newsvc.AppID
	}

	//update service name
	if newsvc.SvcName != "" {
		svc.SvcName = newsvc.SvcName
	}

	//update service desc
	if newsvc.SvcDesc != "" {
		svc.SvcDesc = newsvc.SvcDesc
	}

	if err = svc.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	SendResponse(ctx, nil, iris.Map{"message": "update ok"})
}

// delete a service by id, dangerous action. If the service has active version,
// then we will delete the active version's deployment and service in the cluster,
// and all the version records of the service including the service record itself.
// If the service not have a active version,that is to say there is no deployment
// and service of the service that asked to delete in the cluster, so we will just
// to delete all the version records of the service including the service record
// itself.
func DeleteService(ctx iris.Context) {
	service_id, _ := ctx.Params().GetInt64("id")

	// get the current service object
	service, err := model.GetServiceByID(service_id)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	// current service have active version
	if service.CurrentVersion != "" {
		version := &model.Version{}
		d := model.DB.RWdb.Where("version_name = ?", service.CurrentVersion).Find(&version)
		if d.Error != nil {
			SendResponse(ctx, errno.New(errno.ErrDatabase, d.Error), nil)
			return
		}

		//unmarshal the config
		var version_config model.VersionConfig
		json.Unmarshal([]byte(version.VersionConfig), &version_config)

		if err := DeleteDeploymentAndServiceInCluster(version_config); err != nil {
			SendResponse(ctx, errno.New(errno.ErrDeleteResourceInCluster, err), nil)
			return
		}

	}

	//delete versions which belongs to current service
	d := model.DB.RWdb.Unscoped().Delete(model.Version{}, "svc_id = ?", service_id)
	if d.Error != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, d.Error), nil)
		return
	}

	//delete the service record itself
	if err := model.DeleteService(uint(service_id)); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	notification := mail.NotificationEvent{
		Level:    "Warning",
		UserName: "Admin user",
		Who:      ctx.Values().GetString("current_user_name"),
		Action:   " delete ",
		What:     " service [" + service.SvcName + "]",
		When:     time.Now().String(),
	}

	receptions := []string{}
	var adminUsers []model.User
	d = model.DB.RWdb.Where("role = ?", "admin").Find(&adminUsers)
	if d.Error != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, d.Error), nil)
		return
	}
	for _, admin := range adminUsers {
		receptions = append(receptions, admin.Email)
	}

	mail.SendNotificationEmail(notification, receptions)

	SendResponse(ctx, nil, iris.Map{"id": service_id})
}

//get all services or services that belongs to an app
func GetServiceList(ctx iris.Context) {
	limit := ctx.URLParamIntDefault("limit", 20)    //how many if limit=0,default=20
	offsize := ctx.URLParamIntDefault("offsize", 0) // from where
	app_id := ctx.URLParamIntDefault("app_id", 0)

	var (
		svcs  []*model.Service
		count uint64
		err   error
	)

	if app_id == 0 { //list all, admin only
		if ctx.Values().GetString("current_user_role") == "admin" {
			svcs, count, err = model.ListService(offsize, limit)
		} else {
			ctx.StatusCode(iris.StatusForbidden)
			return
		}
	} else { // list service belongs to an app,login only
		svcs, count, err = model.ListServiceByAppID(offsize, limit, uint(app_id))
	}

	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"count": count, "svcs": svcs})
}

//check whether a service name exist in db
func ServiceNameDuplicateChecker(ctx iris.Context) {
	svcName := ctx.URLParamDefault("svcname", "")

	svc, _ := model.GetServiceByName(svcName)
	if svc.SvcName == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(fmt.Sprintf("service %s not exist", svcName))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.WriteString(fmt.Sprintf("service %s already exist", svc.SvcName))
	return

}