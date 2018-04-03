package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/authorizer"
)

type Identity struct {
	database.Model
	IdentityRole string
}

func (i Identity) Role() string {
	return i.IdentityRole
}

func DefaultMiddleware(
	nodeNumber int64, manager *database.Manager, r *gin.Engine,
	sessionExpiresInSeconds int, authorizerConfig...interface{}) {
	// authenticator 를 초기화한다
	ca := authenticator.NewMem(nodeNumber, sessionExpiresInSeconds)
	ca.SetActions(GetLoadCookieIDToken(manager),
		GetLoadIdentityByCookie(manager),
		GetPersistToken(manager))
	if len(authorizerConfig) > 0 {
		au := authorizer.Init(authorizerConfig...)
		r.Use(manager.Handler(), ca.Handler(), au.Handler())
	} else {
		r.Use(manager.Handler(), ca.Handler())
	}
}
