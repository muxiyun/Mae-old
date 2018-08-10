package handler

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
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
	ctx.WriteString("ok")
}


// the ApplyVersion handler to apply a version config that created before.
// this handler will replace a version with the specified one,
// that is to remove the resources of the previous version config and then
// create the resources about the specified version config
func ApplyVersion(ctx iris.Context){

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
