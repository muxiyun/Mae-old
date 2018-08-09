package handler

import (
	//"fmt"
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/kataras/iris/core/errors"
	"fmt"
)

// create a service
func CreateService(ctx iris.Context) {
	var svc model.Service
	ctx.ReadJSON(&svc)

	if svc.AppID==0 || svc.SvcName==""{
		SendResponse(ctx,errno.New(errno.ServiceNameEmptyorAppIDTypeError,errors.New("")),nil)
		return
	}

	if err:=svc.Create();err!=nil{
		SendResponse(ctx,errno.New(errno.ErrCreateService,err),nil)
	}

	SendResponse(ctx,nil,iris.Map{"id":svc.ID})
}

// get a service info by svc_name
func GetService(ctx iris.Context) {
	svc_name:=ctx.Params().Get("svc_name")
	fmt.Println(svc_name)
	svc,err:=model.GetServiceByName(svc_name)
	if err!=nil{
		SendResponse(ctx, errno.New(errno.ErrGetService, err), nil)
		return
	}
	SendResponse(ctx,nil,svc)
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
	if newsvc.AppID!=0{
		svc.AppID=newsvc.AppID
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

// delete a service by id
func DeleteService(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	if err := model.DeleteService(uint(id)); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"id": id})
}

//get all services or services that belongs to an app
func GetServiceList(ctx iris.Context) {
	limit := ctx.URLParamIntDefault("limit", 20)    //how many if limit=0,default=20
	offsize := ctx.URLParamIntDefault("offsize", 0) // from where
	app_id :=ctx.URLParamIntDefault ("app_id",0)

	var(
		svcs []*model.Service
	    count uint64
	    err error
	)

	if app_id==0{//list all, admin only
		if ctx.Values().GetString("current_user_role")=="admin" {
			svcs, count, err = model.ListService(offsize, limit)
		}else{
			ctx.StatusCode(iris.StatusForbidden)
			return
		}
	}else{// list service belongs to an app,login only
		svcs, count, err = model.ListServiceByAppID(offsize, limit,uint(app_id))
	}

	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"count": count, "svcs": svcs})
}

