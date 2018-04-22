package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/authorizer"
	"github.com/ezaurum/cthulthu/database"
	"github.com/gin-gonic/gin"
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
	sessionExpiresInSeconds int, authorizerConfig ...interface{}) {
	// authenticator 를 초기화한다
	ca := authenticator.NewMem(nodeNumber, sessionExpiresInSeconds)
	ca.SetActions(GetLoadCookieIDToken(manager.DB()),
		GetLoadIdentityByCookie(manager.DB()),
		GetPersistToken(manager.DB()))
	if len(authorizerConfig) > 0 {
		au := authorizer.Init(authorizerConfig...)
		r.Use(database.Handler(manager.DB()), ca.Handler(), au.Handler())
	} else {
		r.Use(manager.Handler(), ca.Handler())
	}
}
