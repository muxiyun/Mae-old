//service 服务增删改查测试文件

package main

import (
	"testing"
	"github.com/muxiyun/Mae/model"
	"github.com/kataras/iris/httptest"
	"time"
)

func TestServiceCRUD(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")


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

	// an admin user to delete a service
	e.DELETE("/api/v1.0/service/{id}").WithPath("id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	// delete a app who has service
	e.DELETE("/api/v1.0/app/{id}").WithPath("id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

}


func TestRecursiveDeleteService(t *testing.T){
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	//defer model.DB.RWdb.DropTableIfExists("users")
	//defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	//defer model.DB.RWdb.DropTableIfExists("apps")
	//defer model.DB.RWdb.DropTableIfExists("versions")
	//defer model.DB.RWdb.DropTableIfExists("services")

	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	CreateAdminForTest(e, "andrewadmin", "andrewadmin123", "3480437308@qq.com")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)
	andrewadmin_token := GetTokenForTest(e, "andrewadmin", "andrewadmin123", 60*60)

	//a normal user to create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "学而1",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")


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


	// create a namespace mae-test-g
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test-h").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	time.Sleep(3*time.Second)

	//create a version which belongs to service xueer_be
	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       1,
		"version_name": "xueer-be-v1",
		"version_desc": "xueer be version 1",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-be-v1-deployment",
				"name_space":  "mae-test-h",
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

	// an admin user to delete a service
	e.DELETE("/api/v1.0/service/{id}").WithPath("id", 2).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	// an admin user to delete a service
	e.DELETE("/api/v1.0/service/{id}").WithPath("id", 1).WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")

	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test-h").WithBasicAuth(andrewadmin_token, "").
		Expect().Body().Contains("OK")
}


