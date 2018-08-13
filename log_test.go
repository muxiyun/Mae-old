package main


import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
)

//注意目前的测试依赖于集群的状态，即集群中手动部署了kube-test命名空间以及deploy,service

func TestGetLog(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")

	CreateUserForTest(e,"andrew","andrew123","andrewpqc@mails.ccnu.edu.cn")
	andrew_token:=GetTokenForTest(e,"andrew","andrew123",60*60)

	CreateAdminForTest(e,"andrew_admin","andrewadmin123","3480437308@qq.com")
	admin_token:=GetTokenForTest(e,"andrew_admin","andrewadmin123",60*60)

	//anonymous to query log in kube-test namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns","kube-test").
		WithPath("pod_name","kube-test-deploy-3598112474-175fp").WithPath("container_name","kube-test-ct").
		Expect().Status(httptest.StatusForbidden)

	//a normal user to query log in kube-test namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns","kube-test").
		WithPath("pod_name","kube-test-deploy-3598112474-175fp").WithPath("container_name","kube-test-ct").
		WithBasicAuth(andrew_token,"").Expect().Body().Contains("OK")

	// a normal user to query log in default
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns","default").
		WithPath("pod_name","my-nginx-3297553814-xp5rf").WithPath("container_name","my-nginx").
		WithBasicAuth(andrew_token,"").Expect().Status(httptest.StatusForbidden)

	// an admin user to query log in kube-test namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns","kube-test").
		WithPath("pod_name","kube-test-deploy-3598112474-175fp").WithPath("container_name","kube-test-ct").
		WithBasicAuth(admin_token,"").Expect().Body().Contains("OK")

	// an admin user to query log in default namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns","default").
		WithPath("pod_name","my-nginx-3297553814-xp5rf").WithPath("container_name","my-nginx").
		WithBasicAuth(admin_token,"").Expect().Body().Contains("OK")

}