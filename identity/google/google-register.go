package google

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/identity"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"github.com/jinzhu/gorm"
)

const (
	ProviderName = "Google"
)

func Register() route.Routes {
	return make(route.Routes).
		POST("/google/register", route.GetProcess("/", CreateIdentity))
}

func CreateIdentity(c *gin.Context, s session.Session, m *gorm.DB) (int, interface{}) {

	var form identity.OAuthIDToken
	err := c.Bind(&form)
	if nil != err {
		panic(err)
	}

	form.Provider = ProviderName

	if len(form.Token) < 1 {
		panic(form)
	}

	if len(form.TokenID) < 1 {
		panic(form)
	}

	//TODO token validation

	tk, b := identity.FindOAuthToken(form, m)
	if b {
		// 토큰은 이미 valid 하니까 덮어씌우자
		//TODo expires
		expires := time.Now().Add(time.Hour * 24 * 365)
		identity.UpdateOAuthToken(tk, form.TokenString(), expires, m)
	} else {
		tk = identity.CreateIdentityByOAuth(form, m, nil)
	}

	//TODO

	ac := authenticator.GetAuthenticator(c)
	ac.Authenticate(c, s, tk)
	return http.StatusFound, "/"

}
