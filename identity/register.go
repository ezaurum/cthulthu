package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

func Register() route.Routes {
	rt := make(route.Routes)
	rt.AddPage("/register", "common/register").
		POST("/register", route.GetProcess("/", CreateIdentity))
	return rt
}

func CreateFormIdentityWithRole(m *database.Manager,
	account string, password string, role string) Identity {

	ft := FormIDToken{}
	if m.IsExist(&ft, FormIDToken{
		AccountName:account,
	}) {
		var i Identity
		m.Find(&i, ft.IdentityID)
		return i
	}

	identity := GetNewIdentity(m)
	identity.IdentityRole = role

	form := FormIDToken{
		AccountName:     account,
		AccountPassword: password,
		Model:           identity.Model,
		IdentityID:      identity.ID,
		//TODO 토큰 암호화
		Token: strconv.FormatInt(identity.ID, 10),
		//TODO 설정 필요
		expires: time.Now().Add(time.Hour * 24 * 365),
	}

	m.CreateAll(&identity, &form)

	return identity
}


func CreateIdentity(c *gin.Context, s session.Session, m *database.Manager) (int, interface{}) {

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
		tk := CreateIdentityByForm(loginForm, m)
		ac := authenticator.GetAuthenticator(c)
		ac.Authenticate(c, s, tk)
		return http.StatusFound, "/"
		break
	case nil:
		// return
		return http.StatusFound, "/register?err=duplicate"
	default:
		panic(findErr)
	}

	return http.StatusFound, "/"
}

func CreateIdentityByForm(registerForm FormIDToken, m *database.Manager) FormIDToken {

	id := GetNewIdentity(m)

	form := FormIDToken{
		AccountName:     registerForm.AccountName,
		AccountPassword: registerForm.AccountPassword,
		Model:           id.Model,
		IdentityID:      id.ID,
		RememberLogin:   registerForm.RememberLogin,
		Token:           strconv.FormatInt(id.ID, 10),
		expires:         time.Now().Add(time.Hour * 24 * 365),
	}

	m.CreateAll(&id, &form)
	return form
}

func CreateIdentityByOAuth(form OAuthIDToken, m *database.Manager) OAuthIDToken {

	id := GetNewIdentity(m)

	f := OAuthIDToken{
		IdentityID: id.ID,
		Provider:   form.Provider,
		Token:      form.Token,
		TokenID:    form.TokenID,
		expires:    time.Now().Add(time.Hour * 24 * 365),
	}

	m.CreateAll(&id, &f)
	return f
}

func GetNewIdentity(m *database.Manager) Identity {
	i := database.Model{
		ID: m.GenerateByType(&Identity{}),
	}
	id := Identity{
		Model: i,
		IdentityRole: "User",
	}
	return id
}

func CreateUsersIfNotExist(manager *database.Manager, defaultUsers []FormIDToken) {
	for _, userForm := range defaultUsers {
		CreateUserIfNotExist(manager, userForm)
	}
}

func CreateUserIfNotExist(manager *database.Manager, token FormIDToken) bool {
	ft := FormIDToken{}
	if manager.IsExist(&ft, token) {
		return false
	}

	CreateIdentityByForm(token, manager)
	return true
}
