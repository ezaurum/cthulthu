package admin

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/test"
	"net/url"
)

const (
	testTemplateDir = "templates"
	testStaticDir = "templates"
)

var (
testConfig = []interface{}{"model.conf", "policy.csv"}
)

func TestRoleAccess(t *testing.T) {

	testDB := database.Test()
	db := testDB.Connect()
	defer db.Close()
	r := initialize(testDB, testTemplateDir, testStaticDir,testConfig...)

	client := test.HttpClient{}
	w0 := client.GetRequest(r, "/")
	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/login", test.GetRedirect(w0))
}

func TestLogin(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r := initialize(testDB, testTemplateDir, testStaticDir,testConfig...)

	form := getRegisterFormPostData()

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)
	w0 := client.PostFormRequest(r, "/login", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))
}

func TestRegister(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r := initialize(testDB, testTemplateDir, testStaticDir,testConfig...)

	//loginForm := webtest.GetStatusOKDoc(r, redirectLocation, t)
	//assert.Equal(t, 1, loginForm.Find("form").Length())
	//assert.Equal(t, 1, loginForm.Find("form").Find("input[name='userID']").Length())

	form := getRegisterFormPostData()

	client := test.HttpClient{}
	w0 := client.PostFormRequest(r, "/register", form)

	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/", test.GetRedirect(w0))

	var result FormIDToken
	b := testDB.IsExist(&result, &FormIDToken{AccountPassword:"test", AccountName:"test"})
	var identity Identity
	b0 := testDB.IsExist(&identity, result.IdentityKey())

	assert.True(t, b)
	assert.True(t, b0)
}

func TestAfterRegisterAuthenticated(t *testing.T) {

	testDB := database.TestNew()
	db := testDB.Connect()
	defer db.Close()
	r := initialize(testDB, testTemplateDir, testStaticDir,testConfig...)

	form := getRegisterFormPostData()

	client := test.HttpClient{}
	client.PostFormRequest(r, "/register", form)

	w := client.GetRequest(r, "/")
	assert.True(t, !test.IsRedirect(w))
}
func getRegisterFormPostData() url.Values {
	form := make(url.Values)
	form.Set("accountName", "test")
	form.Set("accountPassword", "test")
	return form
}
