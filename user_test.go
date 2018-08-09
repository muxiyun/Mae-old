// user 增删改查测试文件


package  main

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	//"fmt"
	//"fmt"
)

func Test404(t *testing.T) {
	e:=httptest.New(t,newApp(),)
	e.GET("/a/unexist/url").Expect().Status(httptest.StatusNotFound)
}



func TestCreateUser(t *testing.T) {
	e := httptest.New(t, newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")


	//test bad request
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"password":"123456",
		"email": "3480437308@qq.com",
	}).Expect().Body().Contains("Bad request")

	//test ok
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew",
		"password":"123456",
		"email":"3480437308@qq.com",
	}).Expect().Body().Contains("OK")

	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew",
		"password":"123456789",
		"email":"123456789@qq.com",
	}).Expect().Body().Contains("Duplicate")

}


func TestDeleteUser(t *testing.T) {
	e := httptest.New(t, newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")

	//test ok
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew",
		"password":"123456",
		"email":"3480437308@qq.com",
	}).Expect().Body().Contains("OK")

	e.DELETE("/api/v1.0/user/1").Expect().Status(httptest.StatusForbidden)

	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrewpqc",
		"password":"123456",
		"email":"andrewpqc@gmail.com",
		"role":"admin",
	}).Expect().Body().Contains("OK")

	e.DELETE("/api/v1.0/user/1").WithBasicAuth("andrewpqc","123456").Expect().Body().Contains("OK")
}

func TestUpdateUser(t *testing.T) {
	e := httptest.New(t, newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")

	//test ok
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew",
		"password":"123456",
		"email":"3480437308@qq.com",
	}).Expect().Body().Contains("OK")

	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"jim",
		"password":"jimpassword",
		"email":"jim@gmail.com",
	}).Expect().Body().Contains("OK")

	e.PUT("/api/v1.0/user/1").WithJSON(map[string]interface{}{
		"username":"andrew2",
		"password":"ppppsssswwwwdddd",
		"email":"andrewpqc@mails.ccnu.edu.cn",
	}).Expect().Status(httptest.StatusForbidden)

	e.PUT("/api/v1.0/user/1").WithBasicAuth("andrew","123456").
		WithJSON(map[string]interface{}{
		"username":"jim",
		"email":"jim@qq.com",
	}).Expect().Body().Contains("Duplicate")

	e.PUT("/api/v1.0/user/1000").WithBasicAuth("andrew","123456").
		WithJSON(map[string]interface{}{
		"username":"hhh",
	}).Expect().Body().Contains("not found")

	e.PUT("/api/v1.0/user/1").WithBasicAuth("andrew","123456").
		WithJSON(map[string]interface{}{
			"username":"andrewpqc",
	}).Expect().Body().Contains("OK")
}


func TestGetUser(t *testing.T) {
	e := httptest.New(t, newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")

	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew",
		"password":"123456",
		"email":"3480437308@qq.com",
	}).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/user/andrewpqc").Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/user/andrew").WithBasicAuth("andrew","123456").
		Expect().Body().Contains("OK")
}


func TestGetUserList(t *testing.T) {
	e := httptest.New(t, newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")

	//add tom
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"tom",
		"password":"tompassword",
		"email":"tom@qq.com",
	}).Expect().Body().Contains("OK")

	//add jim
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"jim",
		"password":"123456jim",
		"email":"jim@qq.com",
	}).Expect().Body().Contains("OK")

	//add bob
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"bob",
		"password":"123456bobpass",
		"email":"bob@qq.com",
		"role":"admin",
	}).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/user").WithQuery("limit",2).WithQuery("offsize",1).
				Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/user").WithQuery("limit",2).Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/user").WithQuery("offsize",1).Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/user").Expect().Status(httptest.StatusForbidden)

	e.GET("/api/v1.0/user").WithQuery("limit",2).WithQuery("offsize",1).
		WithBasicAuth("bob","123456bobpass").Expect().Body().Contains("OK")
	e.GET("/api/v1.0/user").WithQuery("limit",2).WithBasicAuth("bob","123456bobpass").
		Expect().Body().Contains("OK")
	e.GET("/api/v1.0/user").WithQuery("offsize",1).WithBasicAuth("bob","123456bobpass").
		Expect().Body().Contains("OK")
	e.GET("/api/v1.0/user").WithBasicAuth("bob","123456bobpass").
		Expect().Body().Contains("OK")

}

func TestUserInfoDuplicateCheck(t *testing.T) {
	e := httptest.New(t, newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")

	//add tom
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"tom",
		"password":"tompassword",
		"email":"tom@qq.com",
	}).Expect().Body().Contains("OK")

	//add jim
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"jim",
		"password":"123456jim",
		"email":"jim@qq.com",
	}).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/user/duplicate").WithQuery("username","jim").Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/user/duplicate").WithQuery("username","andrew").Expect().Status(httptest.StatusNotFound)
	e.GET("/api/v1.0/user/duplicate").WithQuery("email","jim@qq.com").Expect().Status(httptest.StatusOK)
	e.GET("/api/v1.0/user/duplicate").WithQuery("email","andrewpqc@qq.com").Expect().Status(httptest.StatusNotFound)
}

