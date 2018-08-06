package main


import (
	"testing"
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
)

func TestGetToken(t *testing.T){
	e:=httptest.New(t,newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")

	//create a user
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew",
		"password":"123456",
		"email":"3480437308@qq.com",
		"role":"user",//optional, default is 'user'
	}).Expect().Body().Contains("OK")

	//create admin
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":"andrew2",
		"password":"ppppsssswwwwdddd",
		"email":"andrewpqc@mails.ccnu.edu.cn",
		"role":"admin",
	}).Expect().Body().Contains("OK")

	//Anonymous to get token
	e.GET("/api/v1.0/token").Expect().Status(httptest.StatusForbidden)
	e.GET("/api/v1.0/token").WithBasicAuth("andrew","123456").
		WithQuery("ex",12*60*60).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/token").WithBasicAuth("andrew2","ppppsssswwwwdddd").
		WithQuery("ex",5).Expect().Body().Contains("OK")
}