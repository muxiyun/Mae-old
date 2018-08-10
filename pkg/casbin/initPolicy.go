package casbin

import (
	"fmt"
)

//3个角色/分组(group): roleAdmin,roleUser,roleAnonymous

func myAddPolicy(group, domain, url, method string) {
	if ok := CasbinMiddleware.enforcer.AddPolicy(group, domain, url, method); !ok {
		fmt.Println("[", group, domain, url, method, "]", "already exist")
	}
}

func attachRoleDomain2User(username, rolename, domain string) {
	if ok := CasbinMiddleware.enforcer.AddGroupingPolicy(username, rolename, domain); !ok {
		fmt.Println("[", username, rolename, domain, "]", "already exist")
	}
}

func InitPolicy() {

	//init roleAdmin policy
	for _, policy := range [][]string{
		[]string{"dom_user", "/api/v1.0/user/*", "DELETE"},
		[]string{"dom_user", "/api/v1.0/user", "GET"},

		[]string{"dom_sd", "/api/v1.0/sd/cpu", "GET"},
		[]string{"dom_sd", "/api/v1.0/sd/disk", "GET"},
		[]string{"dom_sd", "/api/v1.0/sd/mem", "GET"},

		[]string{"dom_ns", "/api/v1.0/ns/*", "DELETE"},

		[]string{"dom_app", "/api/v1.0/app/*", "DELETE"},

		[]string{"dom_service","/api/v1.0/service/*","DELETE"},

	} {
		myAddPolicy("roleAdmin", policy[0], policy[1], policy[2])
	}
	fmt.Println("init roleAdmin policy----ok")

	//init roleUser policy
	for _, policy := range [][]string{
		[]string{"dom_user", "/api/v1.0/user/*", "PUT"},
		[]string{"dom_user", "/api/v1.0/user/*", "GET"},

		[]string{"dom_token", "/api/v1.0/token", "GET"},

		[]string{"dom_ns", "/api/v1.0/ns", "GET"},
		[]string{"dom_ns", "/api/v1.0/ns/*", "POST"},

		[]string{"dom_app", "/api/v1.0/app", "POST"},
		[]string{"dom_app", "/api/v1.0/app/*", "GET"},
		[]string{"dom_app", "/api/v1.0/app", "GET"},
		[]string{"dom_app", "/api/v1.0/app/*", "PUT"},
		[]string{"dom_app", "/api/v1.0/app/duplicate", "GET"},

		[]string{"dom_service","/api/v1.0/service","POST"},
		[]string{"dom_service","/api/v1.0/service/*","GET"},
		[]string{"dom_service","/api/v1.0/service/*","PUT"},
		[]string{"dom_service","/api/v1.0/service","GET"},



	} {
		myAddPolicy("roleUser", policy[0], policy[1], policy[2])
	}
	fmt.Println("init roleUser policy----ok")

	//init roleAnonymous policy
	for _, policy := range [][]string{
		[]string{"dom_user", "/api/v1.0/user", "POST"},
		[]string{"dom_user", "/api/v1.0/user/duplicate", "GET"},

		[]string{"dom_sd", "/api/v1.0/sd/health", "GET"},

		//temp rule placed here for test
		[]string{"dom_version","/api/v1.0/version","POST"},
		[]string{"dom_version","/api/v1.0/version/apply","GET"},
	} {
		myAddPolicy("roleAnonymous", policy[0], policy[1], policy[2])
	}

	fmt.Println("init roleAnonymous policy----ok")
}

//将RoleAdmin分组的所有权限授予给username
func AttachToAdmin(username string) {
	for _, dom := range []string{
		"dom_user",
		"dom_sd",
		"dom_ns",
		"dom_app",
		"dom_service",
	} {
		attachRoleDomain2User(username, "roleAdmin", dom)
	}
}

//将RoleUser分组的所有权限授予给username
func AttachToUser(username string) {
	for _, dom := range []string{
		"dom_user",
		"dom_token",
		"dom_ns",
		"dom_app",
		"dom_service",
	} {
		attachRoleDomain2User(username, "roleUser", dom)
	}
}

//将RoleAnonymous分组的所有权限授予给username
func AttachToAnonymous(username string) {
	for _, dom := range []string{
		"dom_user",
		"dom_sd",
	} {
		attachRoleDomain2User(username, "roleAnonymous", dom)
	}
}
