package session

// gin 과 관련해서 세션 관리하는 정보 모음
//TODO 나중에 분리할 때는 세션 라이브러리와 이건 따로 분리해야 함.

import (
	"github.com/gin-gonic/gin"
)

const (
	DefaultSessionContextKey = "session context key tekeli-li tekeli-li"
	DefaultSessionExpires    = 60 * 15
)

func GetSession(c *gin.Context) Session {
	return c.MustGet(DefaultSessionContextKey).(Session)
}

func SetSession(c *gin.Context, session Session) {
	c.Set(DefaultSessionContextKey, session)
}
