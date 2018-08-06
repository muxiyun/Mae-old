package casbin

//目前计划3个角色: admin,user,anonymous

//admin角色拥有全部权限
func AdminRolePolicy() {
	CasbinMiddleware.enforcer.AddPolicy("admin", "*","*", "*")
}


//Enforcer.AddPolicy("admin", "app", "/app/1", "GET")




//高级用户所独有的权限
func UserRolePolicy() {
	//更新用户信息
	CasbinMiddleware.enforcer.AddPolicy("user","dom_user", "/api/v1.0/user/*","PUT")
	//获取特定用户信息
	CasbinMiddleware.enforcer.AddPolicy("user","dom_user","/api/v1.0/user/*","GET")


	//获取token
	CasbinMiddleware.enforcer.AddPolicy("user","dom_token","/api/v1.0/token","GET")
}



//匿名用户独有的权限
func AnonymousRolePolicy() {
	//创建用户的权限，即注册
	CasbinMiddleware.enforcer.AddPolicy("anonymous", "dom_user","/api/v1.0/user", "POST")

	//查看邮箱或用户名是否已经存在
	CasbinMiddleware.enforcer.AddPolicy("anonymous", "dom_user","/api/v1.0/user/duplicate", "GET")

	//健康检查，用于pingServer
	CasbinMiddleware.enforcer.AddPolicy("anonymous","dom_sd","/api/v1.0/sd/health","GET")
}


func AttachAdminToAdminRole(username string){
	CasbinMiddleware.enforcer.AddGroupingPolicy(username,"admin","*")
}


func AttachUserToUserRole(username string){
	CasbinMiddleware.enforcer.AddGroupingPolicy(username,"user","dom_user")
	CasbinMiddleware.enforcer.AddGroupingPolicy(username,"user","dom_token")
	CasbinMiddleware.enforcer.AddGroupingPolicy(username,"anonymous","dom_user")
	CasbinMiddleware.enforcer.AddGroupingPolicy(username,"anonymous","dom_sd")

}

