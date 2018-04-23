package identity

import (
	"fmt"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/test"
	itest "github.com/ezaurum/cthulthu/identity/test"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
	"github.com/ezaurum/cthulthu/database"
)

func TestRegister(t *testing.T) {

	targets := []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}}
	gens := getGenerators(targets...)
	testDB := itest.DB(gens)
	defer testDB.Close()
	r, conf := initializeTest(gens, testDB, authenticator.DefaultSessionExpires)
	route.InitRoute(r, conf.Routes...)

	//loginForm := webtest.GetStatusOKDoc(r, redirectLocation, t)
	//assert.Equal(t, 1, loginForm.Find("form").Length())
	//assert.Equal(t, 1, loginForm.Find("form").Find("input[name='userID']").Length())

	form := make(url.Values)
	form.Set("accountName", "test"+fmt.Sprintf("%v", time.Now()))
	form.Set("accountPassword", "test")

	client := test.HttpClient{}
	w0 := client.PostFormRequest(r, "/register", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))

	var result FormIDToken
	b := database.IsExist(testDB, &result, &FormIDToken{AccountPassword: "test", AccountName: "test"})
	var identity Identity
	b0 := database.IsExist(testDB, &identity, result.IdentityKey())

	assert.True(t, b)
	assert.True(t, b0)
}

func TestAfterRegisterAuthenticated(t *testing.T) {

	targets := []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}}
	gens := getGenerators(targets...)
	testDB := itest.DB(gens)
	defer testDB.Close()
	r, conf := initializeTest(gens, testDB, authenticator.DefaultSessionExpires)
	route.InitRoute(r, conf.Routes...)

	form := getRegisterFormPostData()
	form.Set("accountName", "test"+fmt.Sprintf("%v", time.Now()))

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)

	w := client.GetRequest(r, "/")
	assert.True(t, !test.IsRedirect(w))
}
