package route

import (
	auth "github.com/ezaurum/cthulthu"
	"github.com/ezaurum/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SessionHandlerFunc func(session session.Session) (int, interface{})
type SessionContextHandlerFunc func(c *gin.Context, session session.Session) (int, interface{})

func MakeHTML(gameHandler SessionHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := auth.GetSession(c)
		code, result := gameHandler(session)
		//TODO 페이지 뽑아오는 로직이 필요함
		page := "index"
		c.HTML(code, page, result)
	}
}

func MakeRedirect(sessionHandler SessionContextHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := auth.GetSession(c)
		code, page := sessionHandler(c, session)
		c.Redirect(code, page.(string))
	}
}

func RunRedirect(c *gin.Context, sessionHandler SessionContextHandlerFunc) {
	session := auth.GetSession(c)
	code, page := sessionHandler(c, session)
	c.Redirect(code, page.(string))
}

func RunJson(c *gin.Context, gameHandler SessionContextHandlerFunc) {
	session := auth.GetSession(c)
	code, result := gameHandler(c, session)
	c.JSON(code, result)
}

func MakeJustHTML(page string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, page, nil)
	}
}

func MakeJSON(gameHandler SessionHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := auth.GetSession(c)
		code, result := gameHandler(session)
		c.JSON(code, result)
	}
}

func HTMLOnlyAuth(sessionHandler SessionHandlerFunc, page string, redirect string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := auth.GetSession(c)
		if auth.IsAuthenticated(session) {
			code, result := sessionHandler(session)
			c.HTML(code, page, result)
		} else {
			c.Redirect(http.StatusFound, redirect)
		}
	}
}
