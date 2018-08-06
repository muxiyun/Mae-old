package casbin

import (
	"fmt"
)

//3个角色/分组(group): roleAdmin,roleUser,roleAnonymous

func AdminRolePolicy() {

	//删除特定用户
	CasbinMiddleware.enforcer.AddPolicy("roleAdmin", "dom_user","/api/v1.0/user/*", "DELETE")
	//获取用户列表
	CasbinMiddleware.enforcer.AddPolicy("roleAdmin","dom_user","/api/v1.0/user","GET")


	//查看cpu状态
	CasbinMiddleware.enforcer.AddPolicy("roleAdmin", "dom_sd","/api/v1.0/sd/cpu","GET")
	//查看磁盘状态
	CasbinMiddleware.enforcer.AddPolicy("roleAdmin", "dom_sd","/api/v1.0/sd/disk","GET")
	//查看内存状态
	CasbinMiddleware.enforcer.AddPolicy("roleAdmin", "dom_sd","/api/v1.0/sd/mem","GET")

}



func UserRolePolicy() {
	//更新用户信息
	CasbinMiddleware.enforcer.AddPolicy("roleUser","dom_user", "/api/v1.0/user/*","PUT")
	//获取特定用户信息
	CasbinMiddleware.enforcer.AddPolicy("roleUser","dom_user","/api/v1.0/user/*","GET")


	//获取token
	CasbinMiddleware.enforcer.AddPolicy("roleUser","dom_token","/api/v1.0/token","GET")
}



//匿名用户可以进行的操作
func AnonymousRolePolicy() {
	//创建用户的权限，即注册
	CasbinMiddleware.enforcer.AddPolicy("roleAnonymous", "dom_user","/api/v1.0/user", "POST")
	//查看邮箱或用户名是否已经存在
	CasbinMiddleware.enforcer.AddPolicy("roleAnonymous", "dom_user","/api/v1.0/user/duplicate", "GET")
	//健康检查，用于pingServer
	CasbinMiddleware.enforcer.AddPolicy("roleAnonymous","dom_sd","/api/v1.0/sd/health","GET")
}




func AttachRoleDomain2User(username,rolename,domain string){
	if ok:=CasbinMiddleware.enforcer.AddGroupingPolicy(username,rolename,domain);!ok{
		fmt.Println("the policy already exist")
	}
}

//将RoleAdmin分组的所有权限授予给username
func AttachToAdmin(username string){
	AttachRoleDomain2User(username,"roleAdmin","dom_user")
	AttachRoleDomain2User(username,"roleAdmin","dom_sd")
}

//将RoleUser分组的所有权限授予给username
func AttachToUser(username string){
	AttachRoleDomain2User(username,"roleUser","dom_user")
	AttachRoleDomain2User(username,"roleUser","dom_token")
}

//将RoleAnonymous分组的所有权限授予给username
func AttachToAnonymous(username string){
	AttachRoleDomain2User(username,"roleAnonymous","dom_user")
	AttachRoleDomain2User(username,"roleAnonymous","dom_sd")
}





