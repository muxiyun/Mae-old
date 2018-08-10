package router

import (
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/handler"
	"github.com/muxiyun/Mae/pkg/casbin"
	"github.com/muxiyun/Mae/router/middleware"
)

func Load(app *iris.Application) *iris.Application {

	app.UseGlobal(middleware.TokenChecker)
	app.Use(middleware.NoCache)
	app.Use(middleware.Options)
	app.Use(middleware.Secure)
	app.Use(casbin.CasbinMiddleware.ServeHTTP)

	//routers setup here
	app.OnErrorCode(iris.StatusNotFound, handler.Handle404)
	app.Get("/api/v1.0/token", handler.SignToken)

	user_app := app.Party("/api/v1.0/user")
	{
		user_app.Post("", handler.CreateUser)
		user_app.Delete("/{id:long}", handler.DeleteUser)
		user_app.Put("/{id:long}", handler.UpdateUser)
		user_app.Get("/{username:string}", handler.GetUser)
		user_app.Get("", handler.GetUserList)
		user_app.Get("/duplicate", handler.UserInfoDuplicateChecker)
	}

	sd_app := app.Party("/api/v1.0/sd")
	{
		sd_app.Get("/health", handler.HealthCheck)
		sd_app.Get("/cpu", handler.CPUCheck)
		sd_app.Get("/disk", handler.DiskCheck)
		sd_app.Get("/mem", handler.RAMCheck)
	}

	ns_app := app.Party("/api/v1.0/ns")
	{
		ns_app.Get("", handler.ListNS)
		ns_app.Post("/{ns:string}", handler.CreateNS)
		ns_app.Delete("/{ns:string}", handler.DeleteNS)
	}

	app_app := app.Party("/api/v1.0/app")
	{
		app_app.Post("", handler.CreateApp)
		app_app.Put("/{id:long}", handler.UpdateApp)
		app_app.Delete("/{id:long}", handler.DeleteApp)
		app_app.Get("/{appname:string}", handler.GetApp)
		app_app.Get("", handler.GetAppList)
		app_app.Get("/duplicate", handler.AppNameDuplicateChecker)

	}

	service_app := app.Party("/api/v1.0/service")
	{
		service_app.Post("", handler.CreateService)
		service_app.Put("/{id:long}", handler.UpdateService)
		service_app.Delete("/{id:long}", handler.DeleteService)
		service_app.Get("/{svc_name:string}", handler.GetService)
		service_app.Get("", handler.GetServiceList)
	}

	version_app := app.Party("/api/v1.0/version")
	{
		version_app.Post("", handler.CreateVersion)
		version_app.Put("/{id:long}", handler.UpdateVersion)
		version_app.Delete("/{id:long}", handler.DeleteVersion)
		version_app.Get("/{version_name:string}", handler.GetVersion)
		version_app.Get("", handler.GetVersionList)
		version_app.Get("/apply",handler.ApplyVersion)
		version_app.Get("/unapply",handler.UnapplyVersion)
	}

	return app
}
