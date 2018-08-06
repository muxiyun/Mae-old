package handler


import (
	"fmt"
	"errors"

	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/casbin"

)


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

	if user.Role=="admin"{
		fmt.Println("username",user.UserName)
		casbin.AttachToAdmin(user.UserName)
		casbin.AttachToUser(user.UserName)
		casbin.AttachToAnonymous(user.UserName)
	}else if user.Role=="user" {
		casbin.AttachToUser(user.UserName)
		casbin.AttachToAnonymous(user.UserName)
	}

	SendResponse(ctx, nil, iris.Map{"id": user.ID, "username": user.UserName})
}


func DeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetInt64("id")
	if err := model.DeleteUser(uint(id)); err != nil {
		SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
		return;
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
	fmt.Println("newuser",newuser)

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
	limit := ctx.URLParamIntDefault("limit", 25)    //how many if limit=0,default=20
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
