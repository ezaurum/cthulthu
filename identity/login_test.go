package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/helper"
	itest "github.com/ezaurum/cthulthu/identity/test"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/test"

	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/ezaurum/cthulthu/generators/snowflake"
)

func TestLogin(t *testing.T) {
	targets := []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}}
	gens := snowflake.GetGenerators(0, targets...)

	testDB := itest.DB(gens)
	defer testDB.Close()
	r, conf := initializeTest(gens, testDB, authenticator.DefaultSessionExpires)

	route.InitRoute(r, conf.Routes...)

	form := getRegisterFormPostData()

	client := test.HttpClient{}
	w0 := client.PostFormRequest(r, "/login", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))
}

func TestCookiePersistLogin(t *testing.T) {

	targets := []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}}
	gens := snowflake.GetGenerators(0, targets...)

	testDB := itest.DB(gens)
	defer testDB.Close()
	r, conf := initializeTest(gens, testDB, 1)
	route.InitRoute(r, conf.Routes...)

	form := getRegisterFormPostData()

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)

	form.Set("rememberLogin", "checked")
	client.PostFormRequest(r, "/login", form)

	assert.True(t, helper.Contains(client.Cookies, authenticator.PersistedIDTokenCookieName))

	time.Sleep(time.Second)

	w := client.GetRequest(r, "/")

	assert.True(t, !test.IsRedirect(w))
}
