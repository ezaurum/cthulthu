package admin

import (
	"github.com/ezaurum/cthulthu/database"
	"time"
	"github.com/ezaurum/cthulthu/authenticator"
)

type CookieIDToken struct {
	database.Model
	IdentityID	int64
	expires         time.Time
	Token           string
}

func (l CookieIDToken) TokenString() string {
	return l.Token
}

//TODO 이건 필요없나?
func (l CookieIDToken) IsPersisted() bool {
	return true
}

func (l CookieIDToken) IsExpired() bool {
	return time.Now().After(l.expires)
}

func (l CookieIDToken) IdentityKey() int64 {
	return l.IdentityID
}

var _ authenticator.IDToken = CookieIDToken{}
