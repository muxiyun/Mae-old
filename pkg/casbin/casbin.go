
package casbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kataras/iris"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
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
	domain:="dom_"+strings.Split(path,"/")[3] //dom_*
	fmt.Println(username,domain,path,method)
	return c.enforcer.Enforce(username,domain,path,method)
}

// Username gets the username from db according current_user_id
func Username(ctx iris.Context) string {
	current_user_name:=ctx.Values().GetString("current_user_name")
	if current_user_name==""{
		return "roleAnonymous"
	}
	return current_user_name
}

var (
	CasbinMiddleware *Casbin
)

func init() {
	//a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/mae", true)
	a := gormadapter.NewAdapter("mysql", "root:pqc19960320@tcp(127.0.0.1:3306)/mae",true)
	Enforcer:= casbin.NewEnforcer("conf/casbinmodel.conf", a)
	CasbinMiddleware = New(Enforcer)

	// load policy to memory
	CasbinMiddleware.enforcer.LoadPolicy()
}






