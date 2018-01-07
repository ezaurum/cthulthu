package admin

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/test"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

const (
	testTemplateDir = "templates"
	testStaticDir   = "templates"
)

var (
	testConfig = []interface{}{"model.conf", "policy.csv"}
)

func getRegisterFormPostData() url.Values {
	form := make(url.Values)
	form.Set("accountName", "test")
	form.Set("accountPassword", "test")
	return form
}

func TestRoleAccess(t *testing.T) {

	testDB := database.Test()
	db := testDB.Connect()
	defer db.Close()
	r := initialize(testDB, testTemplateDir, testStaticDir, testConfig...)

	client := test.HttpClient{}
	w0 := client.GetRequest(r, "/")
	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/login", test.GetRedirect(w0))
}
