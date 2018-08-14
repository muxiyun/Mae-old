//version 版本配置增删改查
package main

import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
	//"fmt"
	"time"
)

func TestCreateApplyUnapplyVersion(t *testing.T) {
	time.Sleep(15*time.Second)
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")

	// create a user and get his token
	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)

	// create an admin user and get his token
	CreateAdminForTest(e, "andrewadmin", "admin123", "andrewadmin@gamil.com")
	admin_token := GetTokenForTest(e, "andrewadmin", "admin123", 60*60)

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

	// Anonymous to create a version
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
	}).Expect().Status(httptest.StatusForbidden)

	// normal user  to create a version which belongs to xueer_be
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
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//Anonymous to apply xueer-be-v1
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		Expect().Status(httptest.StatusForbidden)

	// a normal user to apply xueer-be-v1
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// an admin user create another version which belongs to xueer_be
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
	}).WithBasicAuth(admin_token, "").Expect().Body().Contains("OK")

	// an admin user to apply xueer-be-v2
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v2").
		WithBasicAuth(admin_token, "").Expect().Body().Contains("OK")

	// Anonymous to  unapply xueer-be-v2
	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-be-v2").
		Expect().Status(httptest.StatusForbidden)

	// a normal user to unapply xueer-be-v2
	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-be-v2").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// a admin user to delete mae-test namespace to clear the test context.
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")
}

func TestDeleteVersion(t *testing.T) {
	time.Sleep(15*time.Second)
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")

	// create a normal user and get his token
	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)

	// create an admin user and delete namespace xueer
	CreateAdminForTest(e, "andrewadmin", "admin123", "andrewadmin@gamil.com")
	admin_token := GetTokenForTest(e, "andrewadmin", "admin123", 60*60)

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
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// apply the version xueer-be-v1
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create another version xueer-be-v2 which belongs to service xueer-be
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
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// apply xueer-be-v2
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v2").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// anonymous to delete an undeployed version(just to delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 1).Expect().Status(httptest.StatusForbidden)

	// a normal user to delete an undeployed version(just to delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 1).WithBasicAuth(andrew_token, "").
		Expect().Status(httptest.StatusForbidden)

	// an admin user to delete an undeployed version(just to delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 1).WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")

	// anonymous to delete a deployed version(delete the deployment,service in cluster and delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 1).Expect().Status(httptest.StatusForbidden)

	// a normal user to delete a deployed version(delete the deployment,service in cluster and delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 2).WithBasicAuth(andrew_token, "").
		Expect().Status(httptest.StatusForbidden)

	// an admin user to delete a deployed version(delete the deployment,service in cluster and delete the database record)
	e.DELETE("/api/v1.0/version/{id}").WithPath("id", 2).WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")

	// delete namespace mae-test to clear the test context
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")
}

func TestGetVersionAndGetVersionList(t *testing.T) {
	time.Sleep(15*time.Second)
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")

	// create a normal user and get his token
	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)

	// create a admin user and get his token
	CreateAdminForTest(e, "andrewadmin", "admin123", "andrewadmin@gamil.com")
	admin_token := GetTokenForTest(e, "andrewadmin", "admin123", 60*60)

	//create an app
	e.POST("/api/v1.0/app").WithJSON(map[string]interface{}{
		"app_name": "xueer",
		"app_desc": "华师课程挖掘机",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// create a service 'xueer_be' which belongs to　华师匣子　app
	e.POST("/api/v1.0/service").WithJSON(map[string]interface{}{
		"app_id":   1,
		"svc_name": "xueer_be",
		"svc_desc": "the backend part of xueer",
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//create a service 'xueer_fe' which belongs to 华师匣子　app
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
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//apply version "xueer-be-v1"
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-be-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

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
	}).WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// apply version 'xueer-fe-v1'
	e.GET("/api/v1.0/version/apply").WithQuery("version_name", "xueer-fe-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// anonymous get a single version's information
	e.GET("/api/v1.0/version/{version_name}").WithPath("version_name", "xueer-be-v1").
		Expect().Status(httptest.StatusForbidden)

	// a normal user to get a single version's information
	e.GET("/api/v1.0/version/{version_name}").WithBasicAuth(andrew_token, "").
		WithPath("version_name", "xueer-be-v1").Expect().Body().Contains("OK")

	// an admin user to get  a single version's information
	e.GET("/api/v1.0/version/{version_name}").WithBasicAuth(admin_token, "").
		WithPath("version_name", "xueer-be-v1").Expect().Body().Contains("OK")

	//anonymous get all the versions which are belongs to service xueer_be
	e.GET("/api/v1.0/version").WithQuery("service_id", 1).
		Expect().Status(httptest.StatusForbidden)

	// a normal user to get all the versions which are belongs to service xueer_be
	e.GET("/api/v1.0/version").WithQuery("service_id", 1).WithBasicAuth(andrew_token, "").
		Expect().Body().Contains("OK")

	// an admin user to get all the versions which are belongs to service xueer_be
	e.GET("/api/v1.0/version").WithQuery("service_id", 1).WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")

	// anonymous get all the versions in database
	e.GET("/api/v1.0/version").Expect().Status(httptest.StatusForbidden)

	// a normal user to get all the versions in the database
	e.GET("/api/v1.0/version").WithBasicAuth(andrew_token, "").Expect().Status(httptest.StatusForbidden)

	// an admin user to get all the versions in the database
	e.GET("/api/v1.0/version").WithBasicAuth(admin_token, "").Expect().Body().Contains("OK")

	// to unapply xueer-be-v1 (that is to delete the deploy and svc of xueer-be-v1 in the cluster)
	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-be-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// to unapply xueer-fe-v1 (that is to delete the deploy and svc of xueer-fe-v1 in the cluster)
	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-fe-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// delete namespace mae-test to clear the test context
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")
}
