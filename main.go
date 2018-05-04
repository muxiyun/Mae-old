package main

import (
	"MAE/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/hi", func(ctx iris.Context) {
		ctx.JSON(models.Logininfo{Username: "小明", Password: "hhhhh2333"})
	})

	app.Run(iris.Addr("127.0.0.1:8000"),iris.WithoutVersionChecker,iris.WithoutServerError(iris.ErrServerClosed))

}