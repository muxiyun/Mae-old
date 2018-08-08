// some common util function for test

package main

import (
	"encoding/json"
	"github.com/iris-contrib/httpexpect"
	"fmt"
)

type token struct {
	Token string `json:"token"`
}

type tokenResponse struct{
	Code uint `json:"code"`
	Data token `json:"data"`
	Msg string `json:"msg"`
}

//get token by username and password for test
func GetTokenForTest(e *httpexpect.Expect,username,password string)string{
	body:=e.GET("/api/v1.0/token").WithBasicAuth(username,password).Expect().Body().Raw()
	var mytokenResponse tokenResponse
	json.Unmarshal([]byte(body), &mytokenResponse)
	fmt.Println("tttttttttkkkkkkkk",mytokenResponse.Data.Token)
	return mytokenResponse.Data.Token
}


//create a normal user for test
func CreateUserForTest(e *httpexpect.Expect,username,password,email string){
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":username,
		"password":password,
		"email":email,
		"role":"user", //optional, default is 'user'
	}).Expect().Body().Contains("OK")
}


//create a admin user for test
func CreateAdminForTest(e *httpexpect.Expect,username,password,email string){
	e.POST("/api/v1.0/user").WithJSON(map[string]interface{}{
		"username":username,
		"password":password,
		"email":email,
		"role":"admin",
	}).Expect().Body().Contains("OK")
}