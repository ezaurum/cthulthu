package cthulthu

import (
	ezs "github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
)

const (
	DefaultSessionContextKey = "session context key tekeli-li tekeli-li"
	DefaultSessionExpires    = 60 * 15
)

func GetSession(c *gin.Context) ezs.Session {
	return c.MustGet(DefaultSessionContextKey).(ezs.Session)
}

func SetSession(c *gin.Context, session ezs.Session) {
	c.Set(DefaultSessionContextKey, session)
}

