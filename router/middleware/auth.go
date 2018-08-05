package middleware

import (
	"github.com/kataras/iris"
	"github.com/Andrewpqc/Mae/handler"
	"fmt"
	"encoding/base64"
	"strings"
	"github.com/Andrewpqc/Mae/pkg/errno"
	"github.com/Andrewpqc/Mae/pkg/token"
	"github.com/Andrewpqc/Mae/model"
	"time"
	"github.com/kataras/iris/core/errors"
)

func TokenChecker(ctx iris.Context){
	tokenStr:=ctx.GetHeader("Authorization")
	if tokenStr==""{
		ctx.Values().Set("current_user_id","")
		ctx.Values().Set("token_used","0")
		ctx.Next()
	}else {
		result, err := base64.StdEncoding.DecodeString(tokenStr[6:])
		if err != nil {
			handler.SendResponse(ctx,errno.New(errno.ErrDecodeToken,err),nil)
			return
		}
		auth_strs := strings.Split(string(result), ":")
		token_or_username := auth_strs[0]
		password := auth_strs[1]
		fmt.Println("token_or_username:",token_or_username,"passwd",password)
		if token_or_username==""{
			ctx.Values().Set("current_user_id","")
			ctx.Values().Set("token_used","0")
			ctx.Next()
		}

		if password==""{
			//token
			tk:=token.NewJWToken("")
			tkinfo,err:=tk.ParseJWToken(token_or_username)
			if err!=nil{
				handler.SendResponse(ctx,errno.New(errno.ErrDecodeToken,nil),nil)
				return
			}

			id:=tkinfo["id"].(uint64)
			signTime:=tkinfo["signTime"].(int64)
			validdeltatime:=tkinfo["validdeltatime"].(int64)

			if time.Now().Unix()>signTime+validdeltatime{
				handler.SendResponse(ctx,errno.New(errno.ErrTokenExpired,errors.New("expired")),nil)
				return
			}

			ctx.Values().Set("current_user_id",id)
			ctx.Values().Set("token_used","1")
			ctx.Next()
		}

		//username + password
		user,err:=model.GetUserByName(token_or_username)
		if err!=nil{
			handler.SendResponse(ctx,errno.New(errno.ErrDatabase,err),nil)
			return
		}

		if err:=user.Compare(password);err!=nil{
			handler.SendResponse(ctx,errno.New(errno.ErrPasswordIncorrect,err),nil)
			return
		}
		ctx.Values().Set("current_user_id",user.ID)
		ctx.Values().Set("token_used","0")
		ctx.Next()
	}
}




func UsernamePasswordRequired(ctx iris.Context){
	if ctx.Values().Get("current_user_id")!="" && ctx.Values().Get("token_used")=="0"{
		ctx.Next()
	}else{
		handler.SendResponse(ctx,errno.New(errno.ErrUsernamePasswordRequired,errors.New("need username,password,but you give token")),nil)
		return
	}
}


func TokenRequired(ctx iris.Context){
	if ctx.Values().Get("token_used")=="1"{
		ctx.Next()
	}else{
		handler.SendResponse(ctx,errno.New(errno.ErrTokenRequired,errors.New("need token,but you give username and password")),nil)
		return
	}
}



func AdminRequired(ctx iris.Context){

}

func PermissionRequired(ctx iris.Context){

}