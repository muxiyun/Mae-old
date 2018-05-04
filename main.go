package main

import (
	"github.com/Andrewpqc/MAE/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"


)

func newApp() *iris.Application{
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/hi", func(ctx iris.Context) {
		ctx.JSON(models.Logininfo{Username: "小明", Password: "hhhhh2333"})
	})
	return app
}
func main() {
	app:=newApp()
	app.Configure(iris.WithConfiguration(iris.YAML("./config.yml")))
	app.Run(iris.Addr("127.0.0.1:8000"),iris.WithoutVersionChecker,iris.WithoutServerError(iris.ErrServerClosed))
}