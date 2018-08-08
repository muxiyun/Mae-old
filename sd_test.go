// 健康检查,cpu,memory,disk状态获取api测试文件

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
	CreateUserForTest(e,"andrew","123456","3480437308@qq.com")
	andrew_token:=GetTokenForTest(e,"andrew","123456",60*60)

	e.GET("/api/v1.0/sd/health").WithBasicAuth(andrew_token,"").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/cpu").WithBasicAuth(andrew_token,"").
		Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/sd/mem").WithBasicAuth(andrew_token,"").
		Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/sd/disk").WithBasicAuth(andrew_token,"").
		Expect().Status(httptest.StatusForbidden)

	CreateAdminForTest(e,"andrewadmin","123456","admin@qq.com")
	andrewadmin_token:=GetTokenForTest(e,"andrewadmin","123456",60*60)
	e.GET("/api/v1.0/sd/health").WithBasicAuth(andrewadmin_token,"").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/cpu").WithBasicAuth(andrewadmin_token,"").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/mem").WithBasicAuth(andrewadmin_token,"").
		Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/sd/disk").WithBasicAuth(andrewadmin_token,"").
		Expect().Status(httptest.StatusOK)
}
