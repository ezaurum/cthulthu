package authorizer

import (
	"github.com/casbin/casbin"
	"github.com/gin-contrib/sessions"
)

const (
	SessionUserRoleKey = "template/UserRole"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func New(params ...interface{}) *Authorizer {
	e := casbin.NewEnforcer(params...)
	return &Authorizer{
		enforcer: e,
	}
}

func Default() *Authorizer {
	return New("model.conf", "policy.csv")
}

// Authorizer stores the casbin handler
type Authorizer struct {
	enforcer *casbin.Enforcer
}

func SetUserRole(session sessions.Session, userRole string) {
	session.Set(SessionUserRoleKey, userRole)
	session.Save()
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *Authorizer) CheckPermission(path string, method string, userRole string) bool {

	if len(userRole) < 1 {
		userRole = "*"
	}

	return a.enforcer.Enforce(userRole, path, method)
}
