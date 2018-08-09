package casbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
)

func New(e *casbin.Enforcer) *Casbin {
	return &Casbin{enforcer: e}
}

func (c *Casbin) ServeHTTP(ctx iris.Context) {
	if !c.check(ctx) {
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
func (c *Casbin) check(ctx iris.Context) bool {
	username := username(ctx)
	method := ctx.Method()
	path := ctx.Path()
	domain := "dom_" + strings.Split(path, "/")[3] //dom_*
	fmt.Println(username, domain, path, method)
	return c.enforcer.Enforce(username, domain, path, method)
}

// Username gets the username from db according current_user_id
func username(ctx iris.Context) string {
	current_user_name := ctx.Values().GetString("current_user_name")
	if current_user_name == "" {
		return "roleAnonymous"
	}
	return current_user_name
}

var (
	CasbinMiddleware *Casbin
)

func Init() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))

	a := gormadapter.NewAdapter("mysql", dataSource, true)

	Enforcer := casbin.NewEnforcer(viper.GetString("casbinmodel"), a)

	CasbinMiddleware = New(Enforcer)

	// load policy to memory
	CasbinMiddleware.enforcer.LoadPolicy()
}
