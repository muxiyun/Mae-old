package handler


import (
	"net/http"
	"github.com/kataras/iris"
)


func Handle404(ctx iris.Context){
	ctx.StatusCode(http.StatusNotFound)
	ctx.WriteString("Not Found")
}
