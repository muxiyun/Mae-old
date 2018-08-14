package main

import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
	"time"
)

func TestCreateAndDeleteNamespace(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")

	//create a user and an admin
	CreateUserForTest(e, "andrew", "123456", "3480437308@qq.com")
	CreateAdminForTest(e, "andrewadmin", "123456", "admin@qq.com")

	//test an anonymous to create test-ns-1
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "test-ns-1").
		Expect().Status(httptest.StatusForbidden)

	//test a user to create test-ns-2
	andrew_token := GetTokenForTest(e, "andrew", "123456", 60*60)
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "test-ns-2").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//test a admin to create test-ns-3
	andrewadmin_token := GetTokenForTest(e, "andrewadmin", "123456", 60*60)
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "test-ns-3").
		WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	//test an anonymous to delete  test-ns-2
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "test-ns-2").Expect().
		Status(httptest.StatusForbidden)

	time.Sleep(3*time.Second)

	//test a normal user to delete test-ns-2
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "test-ns-2").WithBasicAuth(andrew_token, "").
		Expect().Status(httptest.StatusForbidden)

	// test an admin user to delete test-ns-2 and test-ns-3
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "test-ns-2").WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "test-ns-3").WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

}

func TestListNamespace(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")

	//create a user and an admin
	CreateUserForTest(e, "andrew", "123456", "3480437308@qq.com")
	CreateAdminForTest(e, "andrewadmin", "123456", "admin@qq.com")

	//test an anonymous to list namespaces
	e.GET("/api/v1.0/ns").Expect().Status(httptest.StatusForbidden)

	//test a user to list namespaces
	andrew_token := GetTokenForTest(e, "andrew", "123456", 60*60)
	e.GET("/api/v1.0/ns").WithBasicAuth(andrew_token, "").Expect().
		Body().Contains("OK").NotContains("kube-system")

	//test an admin to list namespaces
	andrewadmin_token := GetTokenForTest(e, "andrewadmin", "123456", 60*60)
	e.GET("/api/v1.0/ns").WithBasicAuth(andrewadmin_token, "").Expect().
		Body().Contains("OK").Contains("kube-system")
}
