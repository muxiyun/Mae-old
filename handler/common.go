package handler

import (
	"net/http"

	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/errno"
)

func SendResponse(c iris.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.StatusCode(http.StatusOK)
	c.JSON(iris.Map{
		"code": code,
		"msg":  message,
		"data": data,
	})
}
