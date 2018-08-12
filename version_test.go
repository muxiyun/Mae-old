//version 版本配置增删改查
package main

import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
	//"fmt"
)

func TestCreateApplyUnapplyVersion(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")

	// create a user and get his token
	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)

	//create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "xueer",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a service which belongs to　华师匣子　app
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "xueer_be",
		"svc_desc": "the backend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a namespace mae-test
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a version which belongs to xueer_be
	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       1,
		"version_name": "xueer-be-v1",
		"version_desc": "xueer be version 1",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-be-v1-deployment",
				"name_space":  "mae-test",
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
	}).Expect().Body().Contains("OK")

	// apply xueer-be-v1
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		Expect().Body().Contains("OK")

	// create another version which belongs to xueer_be
	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       1,
		"version_name": "xueer-be-v2",
		"version_desc": "xueer be version 2",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-be-v2-deployment",
				"name_space":  "mae-test",
				"replicas":    2,
				"labels":      map[string]string{"run": "xueer-be"},
				"containers": [](map[string]interface{}){
					map[string]interface{}{
						"ctr_name":  "xueer-be-v2-ct",
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
				"svc_name": "xueer-be-v2-service",
				"selector": map[string]string{"run": "xueer-be"},
				"labels":   map[string]string{"run": "xueer-be"},
			},
		},
	}).Expect().Body().Contains("OK")

	// apply xueer-be-v2
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v2").
		Expect().Body().Contains("OK")

	// unapply xueer-be-v2
	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-be-v2").
		Expect().Body().Contains("OK")

	// create an admin user and delete namespace mae-test
	CreateAdminForTest(e, "andrewadmin", "admin123", "andrewadmin@gamil.com")
	admin_token := GetTokenForTest(e, "andrewadmin", "admin123", 60*60)
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")
}

func TestDeleteVersion(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")

	// create a user and get his token
	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)

	//create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "xueer",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a service which belongs to　华师匣子　app
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "xueer_be",
		"svc_desc": "the backend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a namespace mae-test
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a version which belongs to xueer_be
	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       1,
		"version_name": "xueer-be-v1",
		"version_desc": "xueer be version 1",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-be-v1-deployment",
				"name_space":  "mae-test",
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
	}).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		Expect().Body().Contains("OK")

	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       1,
		"version_name": "xueer-be-v2",
		"version_desc": "xueer be version 1",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-be-v2-deployment",
				"name_space":  "mae-test",
				"replicas":    1,
				"labels":      map[string]string{"run": "xueer-be"},
				"containers": [](map[string]interface{}){
					map[string]interface{}{
						"ctr_name":  "xueer-be-v2-ct",
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
				"svc_name": "xueer-be-v2-service",
				"selector": map[string]string{"run": "xueer-be"},
				"labels":   map[string]string{"run": "xueer-be"},
			},
		},
	}).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v2").
		Expect().Body().Contains("OK")

	// delete an undeployed version(just to delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 1).Expect().Body().Contains("OK")

	//delete a deployed version(delete the deployment,service in cluster and delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 2).Expect().Body().Contains("OK")

	// create an admin user and delete namespace xueer
	CreateAdminForTest(e, "andrewadmin", "admin123", "andrewadmin@gamil.com")
	admin_token := GetTokenForTest(e, "andrewadmin", "admin123", 60*60)
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")
}

func TestGetVersionAndGetVersionList(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")

	// create a user and get his token
	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)

	//create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "xueer",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a service which belongs to　华师匣子　app
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "xueer_be",
		"svc_desc": "the backend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "xueer_fe",
		"svc_desc": "the frontend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a namespace mae-test
	e.POST("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//create a version which belongs to service xueer_be
	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       1,
		"version_name": "xueer-be-v1",
		"version_desc": "xueer be version 1",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-be-v1-deployment",
				"name_space":  "mae-test",
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
	}).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		Expect().Body().Contains("OK")

	//create a version which belongs to service xueer_fe
	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":       2,
		"version_name": "xueer-fe-v1",
		"version_desc": "xueer fe version 1",
		"version_conf": map[string]interface{}{
			"deployment": map[string]interface{}{
				"deploy_name": "xueer-fe-v1-deployment",
				"name_space":  "mae-test",
				"replicas":    1,
				"labels":      map[string]string{"run": "xueer-fe"},
				"containers": [](map[string]interface{}){
					map[string]interface{}{
						"ctr_name":  "xueer-fe-v1-ct",
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
				"svc_name": "xueer-fe-v1-service",
				"selector": map[string]string{"run": "xueer-fe"},
				"labels":   map[string]string{"run": "xueer-fe"},
			},
		},
	}).Expect().Body().Contains("OK")

	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-fe-v1").
		Expect().Body().Contains("OK")

	// get a single version's information
	e.GET("/api/v1.0/version/{version_name}").WithPath("version_name", "xueer-be-v1").
		Expect().Body().Contains("OK")

	//get all the versions which are belongs to service xueer_be
	e.GET("/api/v1.0/version").WithQuery("service_id", 1).
		Expect().Body().Contains("OK")

	// get all the versions in database(admin only)
	e.GET("/api/v1.0/version").Expect().Status(httptest.StatusForbidden)

	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-be-v1").
		Expect().Body().Contains("OK")

	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-fe-v1").
		Expect().Body().Contains("OK")

	CreateAdminForTest(e, "andrewadmin", "admin123", "andrewadmin@gamil.com")
	admin_token := GetTokenForTest(e, "andrewadmin", "admin123", 60*60)
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")
}
