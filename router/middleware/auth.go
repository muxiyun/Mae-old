package middleware

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"

	"github.com/muxiyun/Mae/handler"
	"github.com/muxiyun/Mae/model"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/token"
)

func TokenChecker(ctx iris.Context) {
	auth_info := ctx.GetHeader("Authorization")
	if auth_info == "" {
		ctx.Values().Set("current_user_name", "")
		ctx.Values().Set("token_used", "0")
		ctx.Next()
		return
	} else {
		result, err := base64.StdEncoding.DecodeString(auth_info[6:])
		if err != nil {
			handler.SendResponse(ctx, errno.New(errno.ErrDecodeToken, err), nil)
			return
		}
		auth_strs := strings.Split(string(result), ":")
		token_or_username := auth_strs[0]
		password := auth_strs[1]
		//fmt.Println("token_or_username:",token_or_username,"passwd",password)
		if token_or_username == "" {
			ctx.Values().Set("current_user_name", "")
			ctx.Values().Set("token_used", "0")
			ctx.Next()
			return
		}

		if password == "" {
			//token
			tk := token.NewJWToken("")
			tkinfo, err := tk.ParseJWToken(token_or_username)
			if err != nil {
				handler.SendResponse(ctx, errno.New(errno.ErrDecodeToken, nil), nil)
				return
			}

			username := tkinfo["username"].(string)
			signTime := tkinfo["signTime"].(float64)
			validdeltatime := tkinfo["validdeltatime"].(float64)

			if time.Now().Unix() > int64(signTime+validdeltatime) {
				handler.SendResponse(ctx, errno.New(errno.ErrTokenExpired, errors.New("expired")), nil)
				return
			}

			user, err := model.GetUserByName(username)
			if err != nil {
				handler.SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
				return
			}

			ctx.Values().Set("current_user_role", user.Role)
			ctx.Values().Set("current_user_id", string(user.ID))
			ctx.Values().Set("current_user_name", user.UserName)
			ctx.Values().Set("token_used", "1")
			ctx.Next()
			return
		}

		//username + password
		user, err := model.GetUserByName(token_or_username)
		if err != nil {
			handler.SendResponse(ctx, errno.New(errno.ErrDatabase, err), nil)
			return
		}

		if err := user.Compare(password); err != nil {
			handler.SendResponse(ctx, errno.New(errno.ErrPasswordIncorrect, err), nil)
			return
		}
		ctx.Values().Set("current_user_name", user.UserName)
		ctx.Values().Set("current_user_role", user.Role)
		ctx.Values().Set("current_user_id", string(user.ID))
		ctx.Values().Set("token_used", "0")
		ctx.Next()
		return
	}
}

//
//func UsernamePasswordRequired(ctx iris.Context){
//	if ctx.Values().Get("current_user_id")!="" && ctx.Values().Get("token_used")=="0"{
//		ctx.Next()
//	}else{
//		handler.SendResponse(ctx,errno.New(errno.ErrUsernamePasswordRequired,errors.New("need username,password,but you give token")),nil)
//		return
//	}
//}
//
//
//func TokenRequired(ctx iris.Context){
//	if ctx.Values().Get("token_used")=="1"{
//		ctx.Next()
//	}else{
//		handler.SendResponse(ctx,errno.New(errno.ErrTokenRequired,errors.New("need token,but you give username and password")),nil)
//		return
//	}
//}
//
//
//
//func AdminRequired(ctx iris.Context){
//
//}
//
//func PermissionRequired(ctx iris.Context){
//
//}
