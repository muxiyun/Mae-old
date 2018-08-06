//sd 健康检查,cpu,memory,disk状态获取api测试文件

package main

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	//"fmt"
)

func TestSystemCheck(t *testing.T) {
	e:=httptest.New(t,newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")

	//anonymous
	e.GET("/api/v1.0/sd/health").Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/cpu").Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/sd/mem").Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/sd/disk").Expect().Status(httptest.StatusForbidden)

	//user, first to register a user
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew",
		"password":"123456",
		"email":"3480437308@qq.com",
		"role":"user", //optional, the default role is user
	}).Expect().Body().Contains("OK")
	e.GET("/api/v1.0/sd/health").WithBasicAuth("andrew","123456").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/cpu").WithBasicAuth("andrew","123456").
		Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/sd/mem").WithBasicAuth("andrew","123456").
		Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/sd/disk").WithBasicAuth("andrew","123456").
		Expect().Status(httptest.StatusForbidden)

	//admin
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrewadmin",
		"password":"123456",
		"email":"admin@qq.com",
		"role":"admin", //optional, the default role is user
	}).Expect().Body().Contains("OK")
	e.GET("/api/v1.0/sd/health").WithBasicAuth("andrewadmin","123456").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/cpu").WithBasicAuth("andrewadmin","123456").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/mem").WithBasicAuth("andrewadmin","123456").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/disk").WithBasicAuth("andrewadmin","123456").
		Expect().Status(httptest.StatusOK)
}
