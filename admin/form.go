package admin

import (
	"github.com/ezaurum/cthulthu/database"
	"time"
	"github.com/ezaurum/cthulthu/authenticator"
)

type FormIDToken struct {
	database.Model
	AccountName     string `form:"accountName" binding:"required"`
	AccountPassword string `form:"accountPassword" binding:"required"`
	RememberLogin   bool   `form:"rememberLogin" gorm:"-"`
	IdentityID		int64
	expires         time.Time
	Token           string
}

func (l FormIDToken) TokenString() string {
	return l.Token
}

func (l FormIDToken) IsPersisted() bool {
	return l.RememberLogin
}

func (l FormIDToken) IsExpired() bool {
	return time.Now().After(l.expires)
}

func (l FormIDToken) IdentityKey() int64 {
	return l.IdentityID
}

var _ authenticator.IDToken = FormIDToken{}
