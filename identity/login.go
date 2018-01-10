package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Login() route.Routes {
	rt := make(route.Routes)
	rt.AddPageWith("/login", "common/login", gin.H{"GoogleClientID":"629871792762-uvt14107uj1shd35lq9i0sgodp20vd77.apps.googleusercontent.com"}).
		POST("/login", route.GetProcess("/",
			func(c *gin.Context, s session.Session, m *database.Manager) (int, interface{}) {

				//TODO 흠? 에러가 나면 걍 무시를 때려야 하나?
				var loginForm FormIDToken
				err := c.Bind(&loginForm)
				if nil != err {
					panic(err)
				}

				var token FormIDToken
				findErr := m.Find(&token, &FormIDToken{AccountName: loginForm.AccountName})
				//TODO 에러를 감싸든가...
				switch findErr {
				case gorm.ErrRecordNotFound:
					return http.StatusFound, "/login?err=not"
					break
				case nil:
					ac := authenticator.GetAuthenticator(c)
					token.RememberLogin = loginForm.RememberLogin
					ac.Authenticate(c, s, &token)
					return http.StatusFound, "/"
				default:
					panic(findErr)
				}

				return http.StatusFound, "/"
			}))
	return rt
}
