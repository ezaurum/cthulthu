package cthulthu

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
)

func Init(r *gin.Engine) {
	handlerFunc, _ := NewAuthorizer("model.conf", "policy.csv")
	r.Use(handlerFunc)
}

const (
	SessionUserRoleKey = "template/UserRole"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func NewAuthorizer(params ... interface{}) ( gin.HandlerFunc, *casbin.Enforcer) {

	e := casbin.NewEnforcer(params...)

	return func(c *gin.Context) {
		a := &BasicAuthorizer{enforcer: e}

		if a.CheckPermission(c) { return }

		writeRequirePermission(c)

	}, e
}
func writeRequirePermission(c *gin.Context) {

	if c.Request.Method != http.MethodGet {
		c.Status(http.StatusForbidden)
	} else {
		redirectToLogin(c)
	}

	println("ABORT")
	c.Abort()
}

func redirectToLogin(c *gin.Context) {
	s := c.Request.URL.Path
	redirect := "/login"
	if s != "/" {
		redirect = redirect + "?redirect=" + s
	}
	c.Redirect(http.StatusFound, redirect)
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}


func SetUserRole(session sessions.Session, userRole string) {
	session.Set(SessionUserRoleKey, userRole)
	session.Save()
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(c *gin.Context) bool {

	session := GetSession(c)
	userRole := session.Get(SessionUserRoleKey)

	if userRole == "" || userRole == nil {
		userRole = "*"
	}

	r := c.Request
	method := r.Method
	path := r.URL.Path

	return a.enforcer.Enforce(userRole, path, method)
}
