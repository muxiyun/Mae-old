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
	return c.enforcer.Enforce(username, path, method)
}

// Username gets the username from db according current_user_id
func Username(ctx iris.Context) string {
	current_user_id:=ctx.Values().Get("current_user_id")

	current_user,_:=model.GetUserByID(current_user_id.(uint))

	return current_user.UserName
}
