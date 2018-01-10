package google

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"github.com/ezaurum/cthulthu/identity"
	"time"
)
const (
	ProviderName = "Google"
)


func Register() route.Routes {
	rt := make(route.Routes)
	rt.POST("/register/google", route.GetProcess("/",
			func(c *gin.Context, s session.Session, m *database.Manager) (int, interface{}) {

				//TODO 흠? 에러가 나면 걍 무시를 때려야 하나?
				var form identity.OAuthIDToken
				err := c.Bind(&form)
				if nil != err {
					panic(err)
				}

				//TODO token validation

				tk, b := FindToken(form, m)
				if b {
					// 토큰은 이미 valid 하니까 덮어씌우자
					identity.UpdateOAuthToken(tk)
					ac := authenticator.GetAuthenticator(c)
					ac.Authenticate(c, s, form)
					return http.StatusFound, "/"
				}

				tk = identity.CreateUserByOAuth(form, m)
				ac := authenticator.GetAuthenticator(c)
				ac.Authenticate(c, s, tk)

				return http.StatusFound, "/"
			}))
	return rt
}
/*
func CreateUserByForm(registerForm FormIDToken, m *database.Manager) authenticator.IDToken {
	i := database.Model{
		ID: m.GenerateByType(&Identity{}),
	}
	id := Identity{
		Model:        i,
		IdentityRole: "User",
	}

	form := FormIDToken{
		AccountName:     registerForm.AccountName,
		AccountPassword: registerForm.AccountPassword,
		Model:           i,
		IdentityID:      id.ID,
		RememberLogin:   registerForm.RememberLogin,
		Token:           strconv.FormatInt(i.ID, 10),
		expires:         time.Now().Add(time.Hour * 24 * 365),
	}

	m.CreateAll(&id, &form)
	return form
}*/
