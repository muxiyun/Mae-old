
package casbin

import (
	"net/http"

	"github.com/kataras/iris"
	"github.com/casbin/casbin"
	"github.com/muxiyun/Mae/model"
	"github.com/casbin/gorm-adapter"

	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strings"
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
	fmt.Println(username)
	if username==""{
		username="anonymous"
	}
	domain:="dom_"+strings.Split(path,"/")[3] //dom_*
	return c.enforcer.Enforce(username,domain,path,method)
}

// Username gets the username from db according current_user_id
func Username(ctx iris.Context) string {
	current_user_id,_:=ctx.Values().GetInt("current_user_id")
	current_user,_:=model.GetUserByID(uint(current_user_id))
	return current_user.UserName
}



var (
	CasbinMiddleware *Casbin
)

func init() {
	a := gormadapter.NewAdapter("mysql", "root:pqc19960320@tcp(127.0.0.1:3306)/")
	Enforcer:= casbin.NewEnforcer("conf/casbinmodel.conf", a)
	CasbinMiddleware = New(Enforcer)

	// load policy to memory
	CasbinMiddleware.enforcer.LoadPolicy()
}






