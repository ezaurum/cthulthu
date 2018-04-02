package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r, _ := initializeTest(testDB, session.DefaultSessionExpires)

	route.InitRoute(r)

	form := getRegisterFormPostData()

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)
	w0 := client.PostFormRequest(r, "/login", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))
}

func TestCookiePersistLogin(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r, conf := initializeTest(testDB, 1)
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
