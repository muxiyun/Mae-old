package router

import (
	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/handler"
	"github.com/muxiyun/Mae/router/middleware"

	"net/http"
)




func Load(app *iris.Application) *iris.Application {

	app.UseGlobal(middleware.TokenChecker)
	app.Use(middleware.NoCache)
	app.Use(middleware.Options)
	app.Use(middleware.Secure)
	app.Use(middleware.CasbinMiddleware.ServeHTTP)

	//routers setup here

	app.OnErrorCode(http.StatusNotFound, handler.Handle404)

	user_app:=app.Party("/api/v1.0/user")
	{
		user_app.Post("",handler.CreateUser)
		user_app.Delete("/{id:long}",handler.DeleteUser)
		user_app.Put("/{id:long}",handler.UpdateUser)
		user_app.Get("/{username:string}",handler.GetUser)
		user_app.Get("",handler.GetUserList)
		user_app.Get("/duplicate",handler.UserInfoDuplicateChecker)
	}

	sd_app:=app.Party("/api/v1.0/sd/")
	{
		sd_app.Get("/health",handler.HealthCheck)
		sd_app.Get("/cpu",handler.CPUCheck)
		sd_app.Get("/disk",handler.DiskCheck)
		sd_app.Get("/mem",handler.RAMCheck)
	}

	app.Get("/api/v1.0/token",handler.SignToken)


	return app
}

