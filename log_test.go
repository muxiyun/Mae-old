package main

import (
	"github.com/kataras/iris/httptest"
	"github.com/muxiyun/Mae/model"
	"testing"
	"time"
)



func TestGetLog(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	defer model.DB.RWdb.DropTableIfExists("users")
	defer model.DB.RWdb.DropTableIfExists("casbin_rule")
	defer model.DB.RWdb.DropTableIfExists("apps")
	defer model.DB.RWdb.DropTableIfExists("services")
	defer model.DB.RWdb.DropTableIfExists("versions")

	CreateUserForTest(e, "andrew", "andrew123", "andrewpqc@mails.ccnu.edu.cn")
	andrew_token := GetTokenForTest(e, "andrew", "andrew123", 60*60)

	CreateAdminForTest(e, "andrew_admin", "andrewadmin123", "3480437308@qq.com")
	admin_token := GetTokenForTest(e, "andrew_admin", "andrewadmin123", 60*60)

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

	time.Sleep(5*time.Second)
	//get a pod's name and a container's name in mae-test namespace
	mae_test_pod_name, mae_test_container_name := GetPodAndContainerNameForTest(e, "mae-test", andrew_token)

	//anonymous to query log in mae-test namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns", "mae-test").
		WithPath("pod_name", mae_test_pod_name).WithPath("container_name", mae_test_container_name).
		Expect().Status(httptest.StatusForbidden)

	//a normal user to query log in mae-test namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns", "mae-test").
		WithPath("pod_name", mae_test_pod_name).WithPath("container_name", mae_test_container_name).
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	//get a pod's name and a container's name in kube-system namespace
	kube_system_pod_name, kube_system_container_name := GetPodAndContainerNameForTest(e, "kube-system", admin_token)

	// a normal user to query log in kube-system
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns", "kube-system").
		WithPath("pod_name", kube_system_pod_name).WithPath("container_name", kube_system_container_name).
		WithBasicAuth(andrew_token, "").Expect().Status(httptest.StatusForbidden)

	// an admin user to query log in mae-test namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns", "mae-test").
		WithPath("pod_name", mae_test_pod_name).WithPath("container_name", mae_test_container_name).
		WithBasicAuth(admin_token, "").Expect().Body().Contains("OK")

	// an admin user to query log in kube-system namespace
	e.GET("/api/v1.0/log/{ns}/{pod_name}/{container_name}").WithPath("ns", "kube-system").
		WithPath("pod_name", kube_system_pod_name).WithPath("container_name", kube_system_container_name).
		WithBasicAuth(admin_token, "").Expect().Body().Contains("OK")

	// to unapply xueer-be-v1 (that is to delete the deploy and svc of xueer-be-v1 in the cluster),for clear test context
	e.GET("/api/v1.0/version/unapply").WithQuery("version_name", "xueer-be-v1").
		WithBasicAuth(andrew_token, "").Expect().Body().Contains("OK")

	// delete namespace mae-test to clear test context
	e.DELETE("/api/v1.0/ns/{ns}").WithPath("ns", "mae-test").WithBasicAuth(admin_token, "").
		Expect().Body().Contains("OK")
}
