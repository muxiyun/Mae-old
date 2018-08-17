package handler

import (
	"errors"
	"time"

	"github.com/kataras/iris"
	"github.com/muxiyun/Mae/pkg/errno"
	"github.com/muxiyun/Mae/pkg/token"
)

func SignToken(ctx iris.Context) {
	validdeltatime := ctx.URLParamInt64Default("ex", 60*60) //validity period,default a hour
	current_user_name := ctx.Values().GetString("current_user_name")

	if current_user_name == "" {
		SendResponse(ctx, errno.New(errno.ErrUnauth, errors.New("need username and password to access")), nil)
		return
	}

	if ctx.Values().Get("token_used") == "0" { //只能使用用户名密码来获取token
		tk := token.NewJWToken("")
		tokenString, err := tk.GenJWToken(map[string]interface{}{
			"username":       current_user_name,
			"signTime":       time.Now().Unix(),
			"validdeltatime": validdeltatime,
		})

		if err != nil {
			SendResponse(ctx, errno.New(errno.ErrToken, err), nil)
			return
		}
		SendResponse(ctx, nil, iris.Map{"token": tokenString})
	} else {
		SendResponse(ctx, errno.New(errno.ErrUsernamePasswordRequired,
			errors.New("need username,password,but you give token")), nil)
	}

}
