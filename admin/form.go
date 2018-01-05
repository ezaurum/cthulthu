package admin

import (
	"github.com/ezaurum/cthulthu/database"
	"time"
)

type LoginForm struct {
}

type FormIDToken struct {
	database.Model
	AccountName     string `form:"accountName" binding:"required"`
	AccountPassword string `form:"accountPassword" binding:"required"`
	RememberLogin   bool   `form:"rememberLogin" gorm:"-"`
	Identity        Identity
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
