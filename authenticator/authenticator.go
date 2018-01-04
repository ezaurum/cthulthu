package authenticator

import (
	"github.com/ezaurum/cthulthu/session"
)

const (
	IDTokenSessionKey = "ID token session key tekelli-li"
)

type Identity interface {
	Role() string
}

type IDToken interface {
	TokenString() string
	IsPersisted() bool
}

func SetIDToken(session session.Session, token IDToken) {
	session.Set(IDTokenSessionKey, token)

	// TODO identifier 가 멀쩡한지 검증 안 해도 되나?

}

func GetIDToken(session session.Session) IDToken {
	a := session.Get(IDTokenSessionKey)
	if nil != a {
		return a.(IDToken)
	}
	return nil
}

func HasIDToken(session session.Session) bool {
	return nil != GetIDToken(session)
}
