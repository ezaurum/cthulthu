package route

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SessionHandlerFunc func(session session.Session) (int, interface{})
type SessionContextHandlerFunc func(c *gin.Context, session session.Session) (int, interface{})
type FullContextHandlerFunc func(c *gin.Context, session session.Session, manager *database.Manager) (int, interface{})

func GetProcess(page string, f FullContextHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbm := database.GetDatatbase(c)
		s := session.GetSession(c)
		code, result := f(c, s, dbm)

		//TODO
		switch code {
		case http.StatusFound:
			c.Redirect(code, result.(string))
			break
		default:
			c.HTML(code, page, result)
			break
		}
	}
}

func MakeHTML(gameHandler SessionHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := session.GetSession(c)
		code, result := gameHandler(session)
		//TODO 페이지 뽑아오는 로직이 필요함
		page := "index"
		c.HTML(code, page, result)
	}
}

func MakeRedirect(sessionHandler SessionContextHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := session.GetSession(c)
		code, page := sessionHandler(c, session)
		c.Redirect(code, page.(string))
	}
}

func RunRedirect(c *gin.Context, sessionHandler SessionContextHandlerFunc) {
	session := session.GetSession(c)
	code, page := sessionHandler(c, session)
	c.Redirect(code, page.(string))
}

func RunJson(c *gin.Context, gameHandler SessionContextHandlerFunc) {
	session := session.GetSession(c)
	code, result := gameHandler(c, session)
	c.JSON(code, result)
}

func MakeJustHTML(page string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, page, nil)
	}
}

func MakeHTMLWith(page string, obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, page, obj)
	}
}



func MakeJSON(gameHandler SessionHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := session.GetSession(c)
		code, result := gameHandler(session)
		c.JSON(code, result)
	}
}

func HTMLOnlyAuth(sessionHandler SessionHandlerFunc, page string, redirect string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := session.GetSession(c)
		//TODO 이거 자체가 잘못 되어 있다. authorizer는 다른 곳에.
		if authenticator.HasIDToken(session) {
			code, result := sessionHandler(session)
			c.HTML(code, page, result)
		} else {
			c.Redirect(http.StatusFound, redirect)
		}
	}
}
