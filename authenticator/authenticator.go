package authenticator

import (
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
)

const (
	IDTokenSessionKey = "ID token session key tekelli-li"
	IdentitySessionKey= "ID session key tekelli-li"
)

type IDTokenLoader func(string) (IDToken, bool)
type IDLoader func(IDToken) (Identity, bool)

type Identity interface {
	Role() string
}

type IDToken interface {
	TokenString() string
	IsPersisted() bool
	IdentityKey() int64
}

func SetIDToken(session session.Session, token IDToken) {
	session.Set(IDTokenSessionKey, token)

	// TODO identifier 가 멀쩡한지 검증 안 해도 되나?

}

func GetIDToken(session session.Session) IDToken {
	a, b := session.Get(IDTokenSessionKey)
	if b {
		return a.(IDToken)
	}
	return nil
}

func HasIDToken(session session.Session) bool {
	return nil != GetIDToken(session)
}

func IsAuthenticated(session session.Session) bool {
	_, b := session.Get(IdentitySessionKey)
	return b
}

func GetIdentity(session session.Session) Identity {
	i, _ :=  session.Get(IdentitySessionKey)
	return i.(Identity)
}

func FindIdentity(session session.Session) (Identity, bool) {
	id, e :=  session.Get(IdentitySessionKey)
	if e {
		return id.(Identity), true
	}
	return nil, false
}

func SetIdentity(session session.Session, identity Identity) {
	session.Set(IdentitySessionKey, identity)
}

func GetAuthenticator(c *gin.Context) Authenticator {
	return c.MustGet(ContextKey).(Authenticator)
}

