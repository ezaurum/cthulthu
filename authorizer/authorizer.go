package authorizer

import (
	"github.com/casbin/casbin"
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

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *Authorizer) CheckPermission(path string, method string, userRole string) bool {

	if len(userRole) < 1 {
		userRole = "*"
	}

	return a.enforcer.Enforce(userRole, path, method)
}
