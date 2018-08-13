//service 服务增删改查测试文件

package main

import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
)

func TestServiceCRUD(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")

	// create two users and get their token
	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	CreateAdminForTest(e, "andrewadmin", "andrewadmin123", "3480437308@qq.com")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)
	andrewadmin_token := GetTokenForTest(e, "andrewadmin", "andrewadmin123", 60*60)

	//a normal user(andrew) to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "学而",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//an admin(andrewadmin) to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "华师匣子",
		"app_desc": "华师校园助手",
	}).WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	// anonymous to create a service
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1, // Be careful, it's type is int,but not string
		"svc_name": "xueer_be",
		"svc_desc": "the backend part of xueer",
	}).Expect().Status(httptest.StatusForbidden)

	// a normal user to create a service
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "xueer_be",
		"svc_desc": "the backend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// an admin to create a service
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "xueer_fe",
		"svc_desc": "frontend part of xueer",
	}).WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	// anonymous to get a service
	e.GET("/api/v1.0/service/{svc_name}").WithPath("svc_name", "xueer_be").
		Expect().Status(httptest.StatusForbidden)

	// a normal user to get a service
	e.GET("/api/v1.0/service/{svc_name}").WithPath("svc_name", "xueer_be").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// a admin user to get a service
	e.GET("/api/v1.0/service/{svc_name}").WithPath("svc_name", "xueer_be").
		WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	// anonymous to update a service
	e.PUT("/api/v1.0/service/{id}").WithPath("id", 1).WithJSON(map[string]interface{}{
		"app_id":   2,
		"svc_name": "XUEER_BE",
		"svc_desc": "xueer backend",
	}).Expect().Status(httptest.StatusForbidden)

	// a normal user to update a service
	e.PUT("/api/v1.0/service/{id}").WithPath("id", 1).WithJSON(map[string]interface{}{
		"app_id":   2,
		"svc_name": "XUEER_BE",
		"svc_desc": "xueer backend",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// a admin user to update a service
	e.PUT("/api/v1.0/service/{id}").WithPath("id", 1).WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "Xueer_Be",
	}).WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	//anonymous to list services
	e.GET("/api/v1.0/service").Expect().Status(httptest.StatusForbidden) // list all
	e.GET("/api/v1.0/service").WithQuery("app_id", 1).
		Expect().Status(httptest.StatusForbidden) // list service belongs to an app

	//a normal user to list service
	e.GET("/api/v1.0/service").WithBasicAuth(andrew_token, "").Expect().Status(httptest.StatusForbidden) // list all
	e.GET("/api/v1.0/service").WithQuery("app_id", 1).WithBasicAuth(andrew_token, "").
		Expect().Body().Contains("OK") // list service belongs to an app

	// admin user to list service
	e.GET("/api/v1.0/service").WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK") // list all
	e.GET("/api/v1.0/service").WithQuery("app_id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK") // list service belongs to an app

	// anonymous to delete a service
	e.DELETE("/api/v1.0/service/{id}").WithPath("id", 1).Expect().Status(httptest.StatusForbidden)

	// a normal user to delete a service
	e.DELETE("/api/v1.0/service/{id}").WithPath("id", 1).WithBasicAuth(andrew_token, "").
		Expect().Status(httptest.StatusForbidden)

	e.DELETE("/api/v1.0/service/{id}").WithPath("id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	// delete a app who has service
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	//这里有一个问题：当app下面还存在service时，不影响app的删除，该app删除之后，对y的service依然存在于数据库中
	//这不是我们想要的，这里应该是app删除之后对应的service也要删除。在service和version的情况中也应该是这样的，
	//我们可以通过设置级联删除或者是在应用中来实现级联删除来达到上面的要求。
}
