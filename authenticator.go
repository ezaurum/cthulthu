package cthulthu

import (
	ezs "github.com/ezaurum/session"
	"time"
)

const (
	IDTokenSessionKey     = "ID token session key tekelli-li"
)

type IDToken interface {
	TokenString() string
	IsPersisted() bool
}

func SetIDToken(session ezs.Session, token IDToken) {
	session.Set(IDTokenSessionKey, token)

	// TODO identifier 가 멀쩡한지 검증 안 해도 되나?

}

func GetIDToken(session ezs.Session) IDToken {
	a :=  session.Get(IDTokenSessionKey)
	if nil != a {
		return a.(IDToken)
	}
	return nil
}

func HasIDToken(session ezs.Session) bool {
	return nil != GetIDToken(session)
}

type LoginIdentity struct {
	id int64
	UserID       string
	UserPassword string
	isPersisted bool
	expires		time.Time
	Token string
}

func (l LoginIdentity) TokenString() string {
	return l.Token
}

func (l LoginIdentity) IsPersisted() bool {
	return l.isPersisted
}

func (l LoginIdentity) IsExpired() bool {
	return time.Now().After(l.expires)
}
