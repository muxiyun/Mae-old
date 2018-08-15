//app ,应用增删改查测试文件
package main

import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
	"time"
)

func TestAppCRUD(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("versions")
	defer model.DB.RWdb.DropTableIfExists("services")

	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	CreateAdminForTest(e, "andrewadmin", "andrewadmin123", "3480437308@qq.com")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)
	andrewadmin_token := GetTokenForTest(e, "andrewadmin", "andrewadmin123", 60*60)

	// Anonymous to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "学而",
		"app_desc": "华师课程挖掘机",
	}).Expect().Status(httptest.StatusForbidden)

	//a normal user to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "学而",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//an admin to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "华师匣子",
		"app_desc": "华师校园助手",
	}).WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	//Anonymous to get an app
	e.GET("/api/v1.0/app/{appname}").WithPath("appname", "学而").Expect().Status(httptest.StatusForbidden)

	//a normal user to get an app
	e.GET("/api/v1.0/app/{appname}").WithPath("appname", "学而").WithBasicAuth(andrew_token, "").
		Expect().Body().Contains("OK")

	// an admin user to get an app
	e.GET("/api/v1.0/app/{appname}").WithPath("appname", "学而").WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	// Anonymous to update a app
	e.PUT("/api/v1.0/app/{id}").WithPath("id", 1).WithJSON(map[string]interface{}{
		"app_name": "xueer",
	}).Expect().Status(httptest.StatusForbidden)

	//a normal user to update an app
	e.PUT("/api/v1.0/app/{id}").WithPath("id", 1).WithJSON(map[string]interface{}{
		"app_name": "Xueer",
		"app_desc": "华师课程挖掘机鸡鸡鸡鸡",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//an admin user to update an app
	e.PUT("/api/v1.0/app/{id}").WithPath("id", 1).WithJSON(map[string]interface{}{
		"app_name": "xueer",
		"app_desc": "山东蓝想挖掘机学学校",
	}).WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	//anonymous to list apps
	e.GET("/api/v1.0/app").Expect().Status(httptest.StatusForbidden)

	//a normal user to list apps
	e.GET("/api/v1.0/app").WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// an admin user to list apps
	e.GET("/api/v1.0/app").WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	// anonymous to delete an app
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 1).Expect().Status(httptest.StatusForbidden)

	//a normal user to delete an app
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 1).WithBasicAuth(andrew_token, "").
		Expect().Status(httptest.StatusForbidden)

	//an admin user to delete an app
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	// anonymous test app_name duplicate checker
	e.GET("/api/v1.0/app/duplicate").WithQuery("appname", "华师匣子").
		Expect().Status(httptest.StatusForbidden)

	// a normal user to test app_name duplicate checker
	e.GET("/api/v1.0/app/duplicate").WithQuery("appname", "华师匣子").
		WithBasicAuth(andrew_token, "").Expect().Body().NotContains("record not found")

	// an admin user to test app_name duplicate checker
	e.GET("/api/v1.0/app/duplicate").WithQuery("appname", "木小犀机器人").
		WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("record not found")
}



func TestRecursiveDeleteApp(t *testing.T){
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("versions")
	defer model.DB.RWdb.DropTableIfExists("services")

	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	CreateAdminForTest(e, "andrewadmin", "andrewadmin123", "3480437308@qq.com")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)
	andrewadmin_token := GetTokenForTest(e, "andrewadmin", "andrewadmin123", 60*60)

	//a normal user to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "学而1",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// delete an app which has no service
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	time.Sleep(3*time.Second)

	//a normal user to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "学而2",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// a normal user to create a service
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   2,
		"svc_name": "xueer_be2",
		"svc_desc": "the backend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// an admin to create a service
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   2,
		"svc_name": "xueer_fe2",
		"svc_desc": "frontend part of xueer",
	}).WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	// delete an app which has two services, but there are no versions of each service
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 2).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	time.Sleep(3*time.Second)

	//a normal user to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "学而3",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// a normal user to create a service
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   3,
		"svc_name": "xueer_be3",
		"svc_desc": "the backend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// an admin to create a service
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   3,
		"svc_name": "xueer_fe3",
		"svc_desc": "frontend part of xueer",
	}).WithBasicAuth(andrewadmin_token, "").Expect().Body().Contains("OK")

	// create a namespace mae-test-g
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test-g").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//create a version which belongs to service xueer_be
	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       3,
		"version_name": "xueer-be-v1",
		"version_desc": "xueer be version 1",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-be-v1-deployment",
				"name_space":  "mae-test-g",
				"replicas":    1,
				"labels":      map[string]string{"run": "xueer-be"},
				"containers": [](map[string]interface{}){
					map[string]interface{}{
						"ctr_name":  "xueer-be-v1-ct",
						"image_url": "pqcsdockerhub/kube-test",
						"start_cmd": []string{"gunicorn", "app:app", "-b", "0.0.0.0:8080", "--log-level", "DEBUG"},
						"ports": [](map[string]interface{}){
							map[string]interface{}{
								"image_port":  8080,
								"target_port": 8090,
								"protocol":    "TCP",
							},
						},
					},
				},
			},
			"svc": map[string]interface{}{
				"svc_name": "xueer-be-v1-service",
				"selector": map[string]string{"run": "xueer-be"},
				"labels":   map[string]string{"run": "xueer-be"},
			},
		},
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//apply version "xueer-be-v1"
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	time.Sleep(3*time.Second)

	// to recursive delete the app and the service of the app and the versions of the service.
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 3).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test-g").WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")
}