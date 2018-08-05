package middleware


import (
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql"
	cm "github.com/muxiyun/Mae/pkg/casbin"
)

var (
	Enforcer *casbin.Enforcer
	CasbinMiddleware *cm.Casbin
)

func init() {
	a := gormadapter.NewAdapter("mysql", "root:pqc19960320@tcp(127.0.0.1:3306)/") // Your driver and data source.
	Enforcer = casbin.NewEnforcer("../../conf/casbinmodel.conf", a)
	CasbinMiddleware = cm.New(Enforcer)
}