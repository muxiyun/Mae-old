//app ,应用增删改查测试文件
package main

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
)

func TestAppCRUD(t *testing.T){
	e := httptest.New(t, newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")

	CreateUserForTest(e,"andrew","andrew123","andrewpqc@mails.ccnu.edu.cn")
	CreateAdminForTest(e,"andrewadmin","andrewadmin123","3480437308@qq.com")
	andrew_token:=GetTokenForTest(e,"andrew","andrew123",60*60)
	andrewadmin_token:=GetTokenForTest(e,"andrewadmin","andrewadmin123",60*60)

	// Anonymous to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name":"学而",
		"app_desc":"华师课程挖掘机",
	}).Expect().Status(httptest.StatusForbidden)

	//a normal user to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name":"学而",
		"app_desc":"华师课程挖掘机",
	}).WithBasicAuth(andrew_token,"").Expect().Body().Contains("OK")

	//an admin to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name":"华师匣子",
		"app_desc":"华师校园助手",
	}).WithBasicAuth(andrewadmin_token,"").Expect().Body().Contains("OK")


	//Anonymous to get an app
	e.GET("/api/v1.0/app/{appname}").WithPath("appname","学而").Expect().Status(httptest.StatusForbidden)

	//a normal user to get an app
	e.GET("/api/v1.0/app/{appname}").WithPath("appname","学而").WithBasicAuth(andrew_token,"").
		Expect().Body().Contains("OK")

	// an admin user to get an app
	e.GET("/api/v1.0/app/{appname}").WithPath("appname","学而").WithBasicAuth(andrewadmin_token,"").
		Expect().Body().Contains("OK")


	// Anonymous to update a app
	e.PUT("/api/v1.0/app/{id}").WithPath("id",1).WithJSON(map[string]interface{}{
		"app_name":"xueer",
	}).Expect().Status(httptest.StatusForbidden)

	//a normal user to update an app
	e.PUT("/api/v1.0/app/{id}").WithPath("id",1).WithJSON(map[string]interface{}{
		"app_name":"Xueer",
		"app_desc":"华师课程挖掘机鸡鸡鸡鸡",
	}).WithBasicAuth(andrew_token,"").Expect().Body().Contains("OK")

	//an admin user to update an app
	e.PUT("/api/v1.0/app/{id}").WithPath("id",1).WithJSON(map[string]interface{}{
		"app_name":"xueer",
		"app_desc":"山东蓝想挖掘机学学校",
	}).WithBasicAuth(andrewadmin_token,"").Expect().Body().Contains("OK")


	//anonymous to list apps
	e.GET("/api/v1.0/app").Expect().Status(httptest.StatusForbidden)

	//a normal user to list apps
	e.GET("/api/v1.0/app").WithBasicAuth(andrew_token,"").Expect().Body().Contains("OK")

	// an admin user to list apps
	e.GET("/api/v1.0/app").WithBasicAuth(andrewadmin_token,"").Expect().Body().Contains("OK")


	// anonymous to delete an app
	e.DELETE("/api/v1.0/app/{id}").WithPath("id",1).Expect().Status(httptest.StatusForbidden)

	//a normal user to delete an app
	e.DELETE("/api/v1.0/app/{id}").WithPath("id",1).WithBasicAuth(andrew_token,"").
		Expect().Status(httptest.StatusForbidden)

	//an admin user to delete an app
	e.DELETE("/api/v1.0/app/{id}").WithPath("id",1).WithBasicAuth(andrewadmin_token,"").
		Expect().Body().Contains("OK")

	// anonymous test app_name duplicate checker
	e.GET("/api/v1.0/app/duplicate").WithQuery("appname","华师匣子").
		Expect().Status(httptest.StatusForbidden)

	// a normal user to test app_name duplicate checker
	e.GET("/api/v1.0/app/duplicate").WithQuery("appname","华师匣子").
		WithBasicAuth(andrew_token,"").Expect().Body().NotContains("record not found")

	// an admin user to test app_name duplicate checker
	e.GET("/api/v1.0/app/duplicate").WithQuery("appname","木小犀机器人").
		WithBasicAuth(andrewadmin_token,"").Expect().Body().Contains("record not found")

}

