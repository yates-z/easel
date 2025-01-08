package test

import (
	"fmt"
	"testing"

	"github.com/yates-z/easel/auth/authorization/rbac"
)

func Test_RBAC(t *testing.T) {
	adapter := rbac.NewCSVAdapter("policy.csv")
	enforcer, err := rbac.NewEnforcer(adapter)
	if err != nil {
		fmt.Println("Failed to create enforcer:", err)
		return
	}
	enforcer.AddPolicy("admin", "/data", "read")
	enforcer.AddPolicy("admin", "/data", "write")
	enforcer.AddPolicy("user", "/data", "read")
	err = enforcer.Save()
	if err != nil {
		fmt.Println("Failed to add policy:", err)
		return
	}

	fmt.Println(enforcer.Enforce("admin", "/data", "write")) // true
	fmt.Println(enforcer.Enforce("user", "/data", "write"))  // false
}
