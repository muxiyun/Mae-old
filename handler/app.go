package handler

import (
	"fmt"
	"time"
	"errors"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/mail"
)

//create a new app
func CreateApp(ctx iris.Context) {
	var app model.App
	ctx.ReadJSON(&app)
	if app.AppName == "" {
		SendResponse(ctx, errno.New(errno.ErrCreateApp, errors.New("app_name can't be empty")), nil)
		return
	}
	if err := app.Create(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrCreateApp, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"message": app.AppName + " created"})
}

// get app info
func GetApp(ctx iris.Context) {
	app_name := ctx.Params().Get("appname")
	app, err := model.GetAppByName(app_name)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrGetApp, err), nil)
		return
	}
	SendResponse(ctx, nil, app)
}


//update the info of a app
func UpdateApp(ctx iris.Context) {
	var newapp model.App
	ctx.ReadJSON(&newapp)

	id, _ := ctx.Params().GetInt64("id")
	app, err := model.GetAppByID(id)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	//update app name
	if newapp.AppName != "" {
		app.AppName = newapp.AppName
	}

	//update app desc
	if newapp.AppDesc != "" {
		app.AppDesc = newapp.AppDesc
	}

	if err = app.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	SendResponse(ctx, nil, iris.Map{"message": "update ok"})
}


//get app list
func GetAppList(ctx iris.Context) {
	limit := ctx.URLParamIntDefault("limit", 20)    //how many if limit=0,default=20
	offsize := ctx.URLParamIntDefault("offsize", 0) // from where

	apps, count, err := model.ListApp(offsize, limit)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"count": count, "apps": apps})
}

//check whether a app name exist in db
func AppNameDuplicateChecker(ctx iris.Context) {
	appname := ctx.URLParamDefault("appname", "")

	app, _ := model.GetAppByName(appname)
	if app.AppName == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(fmt.Sprintf("app %s not exist", appname))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.WriteString(fmt.Sprintf("app %s already exist", appname))
	return

}



// delete an app, dangerous action. it will delete all the resources which belongs to this app.
// such as services,versions and the deployment,service in the cluster
func DeleteApp(ctx iris.Context) {
	app_id, _ := ctx.Params().GetInt64("id")//get the app id

	app,_:=model.GetAppByID(app_id)
	//get services which belongs to the app
	var services []model.Service
	d := model.DB.RWdb.Where("app_id = ?", app_id).Find(&services)
	if d.Error!=nil{
		SendResponse(ctx,errno.New(errno.ErrDatabase,d.Error),nil)
		return
	}

	for _,service:=range services{

		// get current active version of the service and delete the deployments and services
		// in the cluster of the currently active version

		// current_service have active version
		if service.CurrentVersion!=""{
			version:=&model.Version{}
			d := model.DB.RWdb.Where("version_name = ?", service.CurrentVersion).First(&version)
			if d.Error!=nil{
				SendResponse(ctx,errno.New(errno.ErrDatabase,d.Error),nil)
				return
			}

			//unmarshal the config
			var version_config model.VersionConfig
			json.Unmarshal([]byte(version.VersionConfig), &version_config)

			if err:=DeleteDeploymentAndServiceInCluster(version_config);err!=nil{
				SendResponse(ctx,errno.New(errno.ErrDeleteResourceInCluster,err),nil)
				return
			}
		}
		// delete versions record belongs to the service
		d=model.DB.RWdb.Unscoped().Delete(model.Version{}, "svc_id = ?", service.ID)
		if d.Error!=nil{
			SendResponse(ctx,errno.New(errno.ErrDatabase,d.Error),nil)
			return
		}
	}

	//delete service record belongs to the app
	d=model.DB.RWdb.Unscoped().Delete(model.Service{}, "app_id = ?", app_id)
	if d.Error!=nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, d.Error), nil)
		return
	}

	// finally,delete the app record from the database
	if err := model.DeleteApp(uint(app_id)); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	notification:=mail.NotificationEvent{
		Level:"Warning",
		UserName:"Admin user",
		Who:ctx.Values().GetString("current_user_name"),
		Action:" delete ",
		What:" app ["+app.AppName+"]",
		When:time.Now().String(),
	}

	receptions:=[]string{}
	var adminUsers []model.User
	d = model.DB.RWdb.Where("role = ?", "admin").Find(&adminUsers)
	if d.Error!=nil{
		SendResponse(ctx,errno.New(errno.ErrDatabase,d.Error),nil)
		return
	}
	for _,admin:=range adminUsers{
		receptions=append(receptions,admin.Email)
	}

	mail.SendNotificationEmail(notification,receptions)

	SendResponse(ctx, nil, iris.Map{"id": app_id})
}

