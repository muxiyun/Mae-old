package middleware

import (
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	cm "github.com/muxiyun/Mae/pkg/casbin"
)

var (
	Enforcer         *casbin.Enforcer
	CasbinMiddleware *cm.Casbin
)

func init() {
	a := gormadapter.NewAdapter("mysql", "root:pqc19960320@tcp(127.0.0.1:3306)/")
	Enforcer = casbin.NewEnforcer("conf/casbinmodel.conf", a)
	CasbinMiddleware = cm.New(Enforcer)

	// load policy to memory
	Enforcer.LoadPolicy()
}


//目前计划４个角色: admin,user,anonymous

//admin角色拥有全部权限
func CasbinSetAdminPolicy() {
	Enforcer.AddPolicy("admin", "domain_all","*", "*")
}


//Enforcer.AddPolicy("admin", "app", "/app/1", "GET")
//Enforcer.AddGroupingPolicy("alice", "admin", "app")



//高级用户所独有的权限
func CasbinSetUserPolicy() {
	//更新用户信息
	Enforcer.AddPolicy("user","domain_user", "/api/v1.0/user/*","PUT")
	//获取特定用户信息
	Enforcer.AddPolicy("user","domain_user","/api/v1.0/user/*","GET")


	//获取token
	Enforcer.AddPolicy("user","domain_token","/api/v1.0/token","GET")
}



//匿名用户独有的权限
func CasbinSetAnonymousPolicy() {
	//创建用户的权限，即注册
	Enforcer.AddPolicy("anonymous", "user","/api/v1.0/user", "POST")

	//查看邮箱或用户名是否已经存在
	Enforcer.AddPolicy("anonymous", "user","/api/v1.0/user/duplicate", "GET")
}


//默认用户注册可以获得Anonymous+User两个部分的权限