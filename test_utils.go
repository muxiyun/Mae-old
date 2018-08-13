// some common util function for test

package main

import (
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/httpexpect"
	"github.com/muxiyun/Mae/handler"
)

type token struct {
	Token string `json:"token"`
}

type tokenResponse struct {
	Code uint   `json:"code"`
	Data token  `json:"data"`
	Msg  string `json:"msg"`
}

//get token by username and password for test
func GetTokenForTest(e *httpexpect.Expect, username, password string, ex int) string {
	body := e.GET("/api/v1.0/token").WithQuery("ex", ex).
		WithBasicAuth(username, password).Expect().Body().Raw()

	var mytokenResponse tokenResponse
	json.Unmarshal([]byte(body), &mytokenResponse)

	return mytokenResponse.Data.Token
}

//create a normal user for test
func CreateUserForTest(e *httpexpect.Expect, username, password, email string) {
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username": username,
		"password": password,
		"email":    email,
		"role":     "user", //optional, default is 'user'
	}).Expect().Body().Contains("OK")
}

//create a admin user for test
func CreateAdminForTest(e *httpexpect.Expect, username, password, email string) {
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username": username,
		"password": password,
		"email":    email,
		"role":     "admin",
	}).Expect().Body().Contains("OK")
}

type NsResopnse struct {
	Code uint                 `json:"code"`
	Data []handler.PodMessage `json:"data"`
	Msg  interface{}          `json:"msg"`
}

//这里获取指定ns下的一个pod的名称以及该pod下的一个container名称，用于测试
func GetPodAndContainerNameForTest(e *httpexpect.Expect, ns, token string) (string, string) {
	body := e.GET("/api/v1.0/pod/{ns}").WithPath("ns", ns).WithBasicAuth(token, "").
		Expect().Body().Raw()

	var res NsResopnse
	json.Unmarshal([]byte(body), &res)
	fmt.Println(">>>>res:", res)
	fmt.Println("------------->", res.Data[0].PodName, res.Data[0].Containers[0])
	return res.Data[0].PodName, res.Data[0].Containers[0]

}
