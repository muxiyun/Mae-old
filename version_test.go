//version 版本配置增删改查
package main


import (
	"github.com/kataras/iris/httptest"
	//"github.com/muxiyun/Mae/model"
	"testing"
	//"fmt"
)

func TestCreateVersion(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))
	//defer model.DB.RWdb.DropTableIfExists("users")
	//defer model.DB.RWdb.DropTableIfExists("casbin_rule")

	e.POST("/api/v1.0/version").WithJSON(map[string]interface{}{
		"svc_id":1,
		"version_name":"v1",
		"version_desc":"xueer be version 1",
		"version_conf":map[string]interface{}{
			"deployment":map[string]interface{}{
				"deploy_name":"xueer-be-v1-deployment",
				"name_space":"test",
				"replicas":1,
				"labels":[]string{"run=xueer-be","env=test"},
				"pod_labels":[]string{"run=xueer-be","env=test"},
				"containers":[](map[string]interface{}){
					map[string]interface{}{
						"ctr_name":"xueer_be_v1_ct",
						"image_url":"pqcsdockerhub/kube-test",
						"start_cmd":"python manage.py runserver",
						"envs":[](map[string]interface{}){
							map[string]interface{}{
								"env_key":"MYSQL_ORM",
								"env_val":"sb:xxx@x.x.x.x:3306/db",
							},
							map[string]interface{}{
								"env_key":"CONFIG_PATH",
								"env_val":"/path/to/config/file",
							},
						},
						"volumes":[](map[string]interface{}){
							map[string]interface{}{
								"host_path":"/path/in/host",
								"target_path":"/path/in/container",
							},
							map[string]interface{}{
								"host_path":"/path/in/host2",
								"target_path":"/path/in/container2",
							},
						},
						"ports":[](map[string]interface{}){
							map[string]interface{}{
								"image_port":80,
								"target_part":80,
							},
							map[string]interface{}{
								"image_port":443,
								"target_part":443,
							},
						},
					},
				},

			},
			"svc":map[string]interface{}{
				"svc_name":"xueer-be-v1-service",
				"selector":[]string{"run=xueer-be"},
				"svc_type":"clusterip",
			},
		},

	}).Expect().Body().Contains("ok")
}