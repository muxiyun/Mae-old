package handler

import (
	"github.com/kataras/iris"
	"net/http"
)

func Handle404(ctx iris.Context) {
	ctx.StatusCode(http.StatusNotFound)
	ctx.WriteString("Not Found")
}
