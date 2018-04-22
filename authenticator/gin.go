package authenticator

import (
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
)

const (
	ContextKey = "Authenticator-context-key-CTHULTHU"
)

func SetAuthenticator(c *gin.Context, ca Authenticator) {
	c.Set(ContextKey, ca)
}

type Authenticator interface {
	Authenticate(c *gin.Context, session session.Session, idToken IDToken)
	PersistIDToken(c *gin.Context, session session.Session, idToken IDToken)
	SetActions(loadIDToken IDTokenLoader, loadIdentity IDLoader, persistIDToken TokenSaver)
}

func Init(r *gin.Engine) Authenticator {
	ca := Default()
	r.Use(ca.Handler())
	return ca
}
