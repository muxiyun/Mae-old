package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/casbin"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/mail"
	"github.com/muxiyun/Mae/pkg/token"
)

//获取验证链接
func getConfirmLink(ctx iris.Context, user model.User) (string, error) {
	confirmLink := ""
	if ctx.Request().TLS != nil {
		confirmLink += "https://"
	} else {
		confirmLink += "http://"
	}

	confirmLink += ctx.Host()

	confirmLink += "/api/v1.0/user/confirm?tk="

	// generate token
	tk := token.NewJWToken("")
	tokenString, err := tk.GenJWToken(map[string]interface{}{
		"username":       user.UserName,
		"signTime":       time.Now().Unix(),
		"validdeltatime": 30, // 30 minutes
	})
	if err != nil {
		return "", err
	}

	confirmLink += tokenString

	return confirmLink, nil
}


//获取重发请求链接
func getResendLink(ctx iris.Context, username string) string {
	reSendRequestLink := ""
	if ctx.Request().TLS != nil {
		reSendRequestLink += "https://"
	} else {
		reSendRequestLink += "http://"
	}

	reSendRequestLink += ctx.Host()

	reSendRequestLink += "/api/v1.0/user/resend?u="

	//用base64将用户名编码，更安全
	encodeUsername := base64.StdEncoding.EncodeToString([]byte(username))
	reSendRequestLink += encodeUsername

	return reSendRequestLink
}


//邮箱验证邮件过期，重发
func ResendMail(ctx iris.Context) {
	encodeUserName:=ctx.URLParam("u")
	byteUserName,err:=base64.StdEncoding.DecodeString(encodeUserName)
	if err!=nil{
		SendResponse(ctx,errno.New(errno.ErrDecodeToken,err),nil)
		return
	}

	user,err:=model.GetUserByName(string(byteUserName))
	if err!=nil{
		SendResponse(ctx,errno.New(errno.ErrDatabase,err),nil)
		return
	}

	if user.Confirm==true{
		ctx.HTML("<center><h2>邮箱已验证成功，无需再次发送邮件</h2></center>")
		return
	}

	link,err:=getConfirmLink(ctx,*user)
	if err!=nil{
		SendResponse(ctx,errno.New(errno.ErrgenernateConfirmLink,err),nil)
		return
	}

	mail.SendConfirmEmail(mail.ConfirmEvent{UserName:user.UserName,ConfirmLink:link},user.Email)

	SendResponse(ctx,nil,nil)

}


func CreateUser(ctx iris.Context) {
	var user model.User
	ctx.ReadJSON(&user)
	if user.UserName == "" || user.Email == "" || user.PasswordHash == "" {
		SendResponse(ctx, errno.New(errno.BadRequest,
			errors.New("username,email,password can't be null")), nil)
		return
	}

	if err := user.Validate(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrValidation, err), nil)
		return
	}

	if err := user.Encrypt(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrEncrypt, err), nil)
		return
	}

	if err := user.Create(); err != nil {
		fmt.Println(err)
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	link, err := getConfirmLink(ctx, user)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrgenernateConfirmLink, err), nil)
		return
	}

	ce := mail.ConfirmEvent{UserName: user.UserName, ConfirmLink: link}

	mail.SendConfirmEmail(ce, user.Email)

	// Be careful,confirm link can't be return to client in production environment
	// here we return it just for test
	SendResponse(ctx, nil, iris.Map{"id": user.ID, "username": user.UserName,"link":link})
}


func DeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	if err := model.DeleteUser(uint(id)); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"id": id})
}

func GetUser(ctx iris.Context) {
	username := ctx.Params().Get("username")
	user, err := model.GetUserByName(username)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, user)
}


func UpdateUser(ctx iris.Context) {
	var newuser model.User
	ctx.ReadJSON(&newuser)

	id, _ := ctx.Params().GetInt64("id")
	user, err := model.GetUserByID(uint(id))
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	//update emial
	if newuser.Email != "" {
		user.Email = newuser.Email
	}
	//update username
	if newuser.UserName != "" {
		user.UserName = newuser.UserName
	}
	//update password
	if newuser.PasswordHash != "" {
		user.PasswordHash = newuser.PasswordHash
		if err := user.Encrypt(); err != nil {
			SendResponse(ctx, errno.New(errno.ErrEncrypt, err), nil)
			return
		}
	}

	if err = user.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	SendResponse(ctx, nil, iris.Map{"id": user.ID})
}


func GetUserList(ctx iris.Context) {
	limit := ctx.URLParamIntDefault("limit", 20)    //how many if limit=0,default=20
	offsize := ctx.URLParamIntDefault("offsize", 0) // from where

	users, count, err := model.ListUser(offsize, limit)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}
	SendResponse(ctx, nil, iris.Map{"count": count, "users": users})
}


func UserInfoDuplicateChecker(ctx iris.Context) {
	username := ctx.URLParamDefault("username", "")
	email := ctx.URLParamDefault("email", "")
	//检查用户名是否占用
	if username != "" {
		user, _ := model.GetUserByName(username)
		if user.UserName != "" {
			ctx.StatusCode(iris.StatusOK)
			ctx.WriteString(fmt.Sprintf("%s already exist", username))
			return
		}
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(fmt.Sprintf("%s not exist", username))
		return
	}

	//检查邮箱是否占用
	if email != "" {
		user, _ := model.GetUserByEmail(email)
		if user.UserName != "" {
			ctx.StatusCode(iris.StatusOK)
			ctx.WriteString(fmt.Sprintf("%s already exist", email))
			return
		}
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(fmt.Sprintf("%s not exist", email))
		return
	}
}


func ConfirmUser(ctx iris.Context) {
	tokenString := ctx.URLParam("tk")
	tk := token.NewJWToken("")
	tkinfo, err := tk.ParseJWToken(tokenString)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDecodeToken, nil), nil)
		return
	}
	username := tkinfo["username"].(string)
	signTime := tkinfo["signTime"].(float64)
	validdeltatime := tkinfo["validdeltatime"].(float64)

	if time.Now().Unix() > int64(signTime+validdeltatime) {
		link := getResendLink(ctx, username)
		ctx.StatusCode(iris.StatusOK)
		ctx.HTML(fmt.Sprintf("<center><h2>链接已过期</h2><br><a href='%s'>点击重新获取邮件</a></center>", link))
		return
	}

	user, err := model.GetUserByName(username)
	if err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	// confirm user
	user.Confirm = true
	if err = user.Update(); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return
	}

	// to give rights to the user
	if user.Role == "admin" {
		casbin.AttachToAdmin(user.UserName)
		casbin.AttachToUser(user.UserName)
		casbin.AttachToAnonymous(user.UserName)
	} else if user.Role == "user" {
		casbin.AttachToUser(user.UserName)
		casbin.AttachToAnonymous(user.UserName)
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.HTML("<center><h2>恭喜您，邮箱验证成功</h2></center>")

}
