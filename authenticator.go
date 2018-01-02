package cthulthu

import (
	ezs "github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
)

const (
	DefaultSessionContextKey = "default session context key"
	IdentifierSessionKey     = "Identifier session key"
	DefaultSessionExpires = 60*15
)

type Identifier interface {
	ID() int64
	IsPersisted() bool
}

func GetIdentifier(session ezs.Session) Identifier {
	a := session.Get(IdentifierSessionKey)
	return a.(Identifier)
}

func SetIdentifier(session ezs.Session, identifier Identifier) {
	session.Set(IdentifierSessionKey, identifier)
}

func IsAuthenticated(session ezs.Session) bool {
	a := session.Get(IdentifierSessionKey)

	//TODO identifier 가 멀쩡한지 검증 안 해도 되나?

	return a != nil
}

func Authenticate(session ezs.Session, identifier Identifier) {
	SetIdentifier(session, identifier)

	//TODO 로그나 뭐 그런 거?
}

func GetSession(c *gin.Context) ezs.Session {
	return c.MustGet(DefaultSessionContextKey).(ezs.Session)
}
