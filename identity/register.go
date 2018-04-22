package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

func Register(generator generators.IDGenerator) route.Routes {
	rt := make(route.Routes)
	rt.AddPage("/register", "common/register").
		POST("/register", route.GetProcess("/", CreateIdentity(generator)))
	return rt
}

func CreateFormIdentityWithRole(m *gorm.DB,
	generator generators.IDGenerator,
	account string, password string, role string) Identity {

	ft := FormIDToken{}
	if database.IsExist(m, &ft, FormIDToken{
		AccountName: account,
	}) {
		var i Identity
		m.Find(&i, ft.IdentityID)
		return i
	}

	identity := GetNewIdentity(generator)
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

	database.CreateAll(m, &identity, &form)

	return identity
}

func CreateIdentity(generator generators.IDGenerator) func(c *gin.Context, s session.Session, m *gorm.DB) (int, interface{}) {
	return func(c *gin.Context, s session.Session, m *gorm.DB) (int, interface{}) {
		//TODO 흠? 에러가 나면 걍 무시를 때려야 하나?
		var loginForm FormIDToken
		err := c.Bind(&loginForm)
		if nil != err {
			panic(err)
		}

		var token FormIDToken
		findErr := m.Find(&token, &FormIDToken{AccountName: loginForm.AccountName})
		//TODO 에러를 감싸든가...
		switch findErr.Error {
		case gorm.ErrRecordNotFound:
			tk := CreateIdentityByForm(loginForm, m, generator)
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
}

func CreateIdentityByForm(registerForm FormIDToken, m *gorm.DB, generator generators.IDGenerator) FormIDToken {

	id := GetNewIdentity(generator)

	form := FormIDToken{
		AccountName:     registerForm.AccountName,
		AccountPassword: registerForm.AccountPassword,
		Model:           id.Model,
		IdentityID:      id.ID,
		RememberLogin:   registerForm.RememberLogin,
		Token:           strconv.FormatInt(id.ID, 10),
		expires:         time.Now().Add(time.Hour * 24 * 365),
	}

	database.CreateAll(m, &id, &form)
	return form
}

func CreateIdentityByOAuth(form OAuthIDToken, m *gorm.DB, generator generators.IDGenerator) OAuthIDToken {

	id := GetNewIdentity(generator)

	f := OAuthIDToken{
		IdentityID: id.ID,
		Provider:   form.Provider,
		Token:      form.Token,
		TokenID:    form.TokenID,
		expires:    time.Now().Add(time.Hour * 24 * 365),
	}

	database.CreateAll(m, &id, &f)
	return f
}

func GetNewIdentity(generator generators.IDGenerator) Identity {
	i := database.Model{
		ID: generator.GenerateInt64(),
	}
	id := Identity{
		Model:        i,
		IdentityRole: "User",
	}
	return id
}

func CreateUsersIfNotExist(db *gorm.DB, generator generators.IDGenerator, defaultUsers []FormIDToken) {
	for _, userForm := range defaultUsers {
		CreateUserIfNotExist(db, userForm, generator)
	}
}

func CreateUserIfNotExist(db *gorm.DB, token FormIDToken, generator generators.IDGenerator) bool {
	ft := FormIDToken{}
	if database.IsExist(db, &ft, token) {
		return false
	}

	CreateIdentityByForm(token, db, generator)
	return true
}
