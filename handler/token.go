package handler

import (
	"time"
	"errors"

	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/token"
	"github.com/muxiyun/Mae/pkg/errno"
)


func SignToken(ctx iris.Context){
	validdeltatime:=ctx.URLParamInt64Default("ex",60*60)//validity period,default a hour
	current_user_id:=ctx.Values().Get("current_user_id")

	if ctx.Values().Get("token_used")=="0"{//只能使用用户名密码来获取token
		tk:=token.NewJWToken("")
		tokenString,err:=tk.GenJWToken(map[string]interface{}{
			"id":current_user_id,
			"signTime":time.Now().Unix(),
			"validdeltatime":validdeltatime,
		})

		if err!=nil{
			SendResponse(ctx,errno.New(errno.ErrToken,err),nil)
			return
		}

		SendResponse(ctx,nil,iris.Map{"token":tokenString})
	} else {
		SendResponse(ctx,errno.New(errno.ErrUsernamePasswordRequired,
			errors.New("need username,password,but you give token")),nil)
	}

}


