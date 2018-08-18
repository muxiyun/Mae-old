package handler

import (
	"github.com/kataras/iris"

)

func Handle404(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNotFound)
	ctx.WriteString("Not Found")
}
