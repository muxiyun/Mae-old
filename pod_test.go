package main


import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
)

func TestListPod(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")

	CreateUserForTest(e,"andrew","andrew123","andrewpqc@mails.ccnu.edu.cn")
	andrew_token:=GetTokenForTest(e,"andrew","andrew123",60*60)

	CreateAdminForTest(e,"andrew_admin","andrewadmin123","3480437308@qq.com")
	admin_token:=GetTokenForTest(e,"andrew_admin","andrewadmin123",60*60)

	// anonymous to get pods in kube-test namespace
	e.GET("/api/v1.0/pod/{ns}").WithPath("ns","kube-test").Expect().Status(httptest.StatusForbidden)

	// a normal user to get pods in kube-test namespace
	e.GET("/api/v1.0/pod/{ns}").WithPath("ns","kube-test").WithBasicAuth(andrew_token,"").
		Expect().Body().Contains("OK")

	// a normal user to get pods in kube-public namespace
	e.GET("/api/v1.0/pod/{ns}").WithPath("ns","kube-public").WithBasicAuth(andrew_token,"").
		Expect().Status(httptest.StatusForbidden)

	// an admin user to get pods in kube-test namespace
	e.GET("/api/v1.0/pod/{ns}").WithPath("ns","kube-test").WithBasicAuth(admin_token,"").
		Expect().Body().Contains("OK")

	// an admin user to get pods in kube-public namespace
	e.GET("/api/v1.0/pod/{ns}").WithPath("ns","kube-public").WithBasicAuth(admin_token,"").
		Expect().Body().Contains("OK")
}