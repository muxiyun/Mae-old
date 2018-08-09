package handler

import (
	"errors"

	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
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

// delete a app,dangerous
func DeleteApp(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	if err := model.DeleteApp(uint(id)); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return;
	}
	SendResponse(ctx, nil, iris.Map{"id": id})
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

	if appname != "" {
		app, err := model.GetAppByName(appname)
		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrDatabase, err),
				iris.Map{"message": app.AppName + " not exists"})
			return
		}
		SendResponse(ctx, nil, iris.Map{"message": appname + " exists"})
		return
	}

	SendResponse(ctx, errno.New(errno.ErrAppNameNotProvide, errors.New("")), nil)
}
