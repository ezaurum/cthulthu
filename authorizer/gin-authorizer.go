package authorizer

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	IdentityContextKey = "identity context key"
)

type AuthorizeMiddleware struct {
	authorizer *Authorizer
}

func InitWithAuthenticator(r *gin.Engine) {
	authenticator.Init(r)
	Init(r)
}

func Init(r *gin.Engine) {


	auth := AuthorizeMiddleware{
		authorizer: Default(),
	}

	r.Use(auth.Handler())

}

func (a *AuthorizeMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.CheckPermission(c) {
			return
		}
		writeRequirePermission(c)
	}
}

//TODO 변경이 가능해야지...

func redirectToLogin(c *gin.Context) {
	s := c.Request.URL.Path
	redirect := "/login"
	if s != "/" {
		redirect = redirect + "?redirect=" + s
	}
	c.Redirect(http.StatusFound, redirect)
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

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *AuthorizeMiddleware) CheckPermission(c *gin.Context) bool {

	path := c.Request.URL.Path
	method := c.Request.Method
	userRole := "*"

	identity, exists := c.Get(IdentityContextKey)
	if exists {
		userRole = identity.(authenticator.Identity).Role()
	}

	return a.authorizer.CheckPermission(path, method, userRole)
}
