package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"time"
)

type OAuthIDToken struct {
	database.Model
	IdentityID int64
	expires    time.Time
	TokenID  string `form:"tokenID" binding:"required"`
	Token string `form:"token" binding:"required"`
	Provider string
}

func (l OAuthIDToken) TokenString() string {
	return l.Token
}

func (l OAuthIDToken) IsPersisted() bool {
	return true
}

func (l OAuthIDToken) IsExpired() bool {
	return time.Now().After(l.expires)
}

func (l OAuthIDToken) IdentityKey() int64 {
	return l.IdentityID
}

var _ authenticator.IDToken = OAuthIDToken{}

type CookieIDToken struct {
	database.Model
	IdentityID int64
	Expires    time.Time
	Token      string
}

func (l CookieIDToken) TokenString() string {
	return l.Token
}

//TODO 이건 필요없나?
func (l CookieIDToken) IsPersisted() bool {
	return true
}

func (l CookieIDToken) IsExpired() bool {
	return time.Now().After(l.Expires)
}

func (l CookieIDToken) IdentityKey() int64 {
	return l.IdentityID
}

var _ authenticator.IDToken = CookieIDToken{}

type FormIDToken struct {
	database.Model
	AccountName     string `form:"accountName" binding:"required"`
	AccountPassword string `form:"accountPassword" binding:"required"`
	RememberLogin   string   `form:"rememberLogin" gorm:"-"`
	IdentityID      int64
	expires         time.Time
	Token           string
}

func (l FormIDToken) TokenString() string {
	return l.Token
}

func (l FormIDToken) IsPersisted() bool {
	return "checked" == l.RememberLogin
}

func (l FormIDToken) IsExpired() bool {
	return time.Now().After(l.expires)
}

func (l FormIDToken) IdentityKey() int64 {
	return l.IdentityID
}

var _ authenticator.IDToken = FormIDToken{}


