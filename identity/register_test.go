package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"net/url"
	"time"
	"fmt"
)

func TestRegister(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r, conf := initializeTest(testDB, session.DefaultSessionExpires)
	route.InitRoute(r, conf.Routes...)

	//loginForm := webtest.GetStatusOKDoc(r, redirectLocation, t)
	//assert.Equal(t, 1, loginForm.Find("form").Length())
	//assert.Equal(t, 1, loginForm.Find("form").Find("input[name='userID']").Length())

	form := make(url.Values)
	form.Set("accountName", "test"+fmt.Sprintf("%v",time.Now()))
	form.Set("accountPassword", "test")

	client := test.HttpClient{}
	w0 := client.PostFormRequest(r, "/register", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))

	var result FormIDToken
	b := testDB.IsExist(&result, &FormIDToken{AccountPassword: "test", AccountName: "test"})
	var identity Identity
	b0 := testDB.IsExist(&identity, result.IdentityKey())

	assert.True(t, b)
	assert.True(t, b0)
}

func TestAfterRegisterAuthenticated(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r, conf := initializeTest(testDB, session.DefaultSessionExpires)
	route.InitRoute(r, conf.Routes...)

	form := getRegisterFormPostData()
	form.Set("accountName", "test"+fmt.Sprintf("%v",time.Now()))

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)

	w := client.GetRequest(r, "/")
	assert.True(t, !test.IsRedirect(w))
}
