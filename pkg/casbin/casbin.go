//改编自iris-casbin

package casbin

import (
	"net/http"

	"github.com/kataras/iris"
	"github.com/casbin/casbin"
	"github.com/muxiyun/Mae/model"
)

func New(e *casbin.Enforcer) *Casbin {
	return &Casbin{enforcer: e}
}


func (c *Casbin) ServeHTTP(ctx iris.Context) {
	if !c.Check(ctx) {
		ctx.StatusCode(http.StatusForbidden) // Status Forbiden
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

type Casbin struct {
	enforcer *casbin.Enforcer
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
func (c *Casbin) Check(ctx iris.Context) bool {
	username := Username(ctx)
	method := ctx.Method()
	path := ctx.Path()
	if username==""{
		return true
	}
	return c.enforcer.Enforce(username, path, method)//授权通过则返回true
}

// Username gets the username from db according current_user_id
func Username(ctx iris.Context) string {
	current_user_id,_:=ctx.Values().GetInt("current_user_id")
	current_user,_:=model.GetUserByID(uint(current_user_id))
	return current_user.UserName
}
