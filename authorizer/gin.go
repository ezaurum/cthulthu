package authorizer

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthorizeMiddleware struct {
	authorizer *Authorizer
}

func InitWithAuthenticator(r *gin.Engine, config ...interface{}) (authenticator.Authenticator, *AuthorizeMiddleware) {
	ac := authenticator.Init(r)
	auth := Init(config...)
	return ac, auth
}

func Init(config ...interface{}) *AuthorizeMiddleware {
	auth := GetAuthorizer(config...)
	return &auth
}

func GetAuthorizer(config ...interface{}) AuthorizeMiddleware {
	var authorizer *Authorizer
	if len(config) < 1 {
		authorizer = Default()
	} else {
		authorizer = New(config...)
	}
	auth := AuthorizeMiddleware{
		authorizer: authorizer,
	}
	return auth
}

func (a *AuthorizeMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.CheckPermission(c, authenticator.GetSession(c)) {
			return
		}
		writeRequirePermission(c)
	}
}

//TODO 변경이 가능해야지...

func redirectToLogin(c *gin.Context) {
	s := c.Request.URL.Path
	redirect := "/login"
	switch s {
	case "/":
	case "/login":
	case "/register":
		break
	default:
		redirect = redirect + "?redirect=" + s
		break
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
func (a *AuthorizeMiddleware) CheckPermission(c *gin.Context, session session.Session) bool {

	path := c.Request.URL.Path
	method := c.Request.Method
	userRole := "*"

	if authenticator.IsAuthenticated(session) {
		userRole = authenticator.GetIdentity(session).Role()
	}

	return a.authorizer.CheckPermission(path, method, userRole)
}
