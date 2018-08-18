// some common util function for test

package main

import (
	"encoding/json"
	"github.com/iris-contrib/httpexpect"
	"github.com/muxiyun/Mae/handler"

	"strings"
	// "fmt"
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




type LinkData struct{
	ID int `json:"id"`
	Username string `json:"username"`
	Link string `json:"link"`
}

type createUserReturnData struct{
	Code int `json:"code"`
	Msg int `json:"msg"`
	Data LinkData `json:"data"`
}

//create a normal user and confirm the email for test
func CreateUserForTest(e *httpexpect.Expect, username, password, email string) {
	returnData:=e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username": username,
		"password": password,
		"email":    email,
		"role":     "user", //optional, default is 'user'
	}).Expect().Body().Raw()

	var rd createUserReturnData
	json.Unmarshal([]byte(returnData),&rd)

	req:=strings.Split(rd.Data.Link[21:],"?")
	e.GET(req[0]).WithQuery("tk",req[1][3:]).Expect().Body().Contains("验证成功")

}

//create a admin user　and confirm the email for test
func CreateAdminForTest(e *httpexpect.Expect, username, password, email string) {
	returnData:=e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username": username,
		"password": password,
		"email":    email,
		"role":     "admin",
	}).Expect().Body().Raw()

	var rd createUserReturnData
	json.Unmarshal([]byte(returnData),&rd)

	req:=strings.Split(rd.Data.Link[21:],"?")
	e.GET(req[0]).WithQuery("tk",req[1][3:]).Expect().Body().Contains("验证成功")

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
	return res.Data[0].PodName, res.Data[0].Containers[0]
}
