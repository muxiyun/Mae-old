package main


import (
	"testing"
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"time"
)


func TestGetToken(t *testing.T){
	e:=httptest.New(t,newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")

	//create a user and an admin
	CreateUserForTest(e,"andrew","123456","3480437308@qq.com")
	CreateAdminForTest(e,"andrew2","ppppsssswwwwdddd","andrewpqc@mails.ccnu.edu.cn")

	//Anonymous to get token
	e.GET("/api/v1.0/token").Expect().Status(httptest.StatusForbidden)
	//normal user to get token
	e.GET("/api/v1.0/token").WithBasicAuth("andrew","123456").
		WithQuery("ex",12*60*60).Expect().Body().Contains("OK")
	//admin user to get token
	e.GET("/api/v1.0/token").WithBasicAuth("andrew2","ppppsssswwwwdddd").
		Expect().Body().Contains("OK")
}


func TestTokenExpire(t *testing.T){
	e:=httptest.New(t,newApp(),httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")

	CreateAdminForTest(e,"andrew","123456","3480437308@qq.com")
	token:=GetTokenForTest(e,"andrew","123456",3)

	time.Sleep(4*time.Second) //sleep for 4 seconds,wait for the token expires

	//admin user can get the ns list,but the token expired,so it can not get the
	//ns list now
	e.GET("/api/v1.0/ns").WithBasicAuth(token,"").Expect().
		Body().Contains("Token expired")
}
