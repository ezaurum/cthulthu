package admin

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/test"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
)

func getRegisterFormPostData() url.Values {
	form := make(url.Values)
	form.Set("accountName", "test")
	form.Set("accountPassword", "test")
	return form
}

func initializeTest(manager *database.Manager, expires int) *gin.Engine {
	return initialize(&Config{
		DBManager:               manager,
		TemplateDir:             testTemplateDir,
		StaticDir:               testStaticDir,
		NodeNumber:              0,
		SessionExpiresInSeconds: expires,
		AutoMigrates:            []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}},
		AuthorizerConfig:        testAuthorizerConfig,
	})
}

func TestRoleAccess(t *testing.T) {

	testDB := database.Test()
	db := testDB.Connect()
	defer db.Close()
	r := initializeTest(testDB, session.DefaultSessionExpires)

	client := test.HttpClient{}
	w0 := client.GetRequest(r, "/")
	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/login", test.GetRedirect(w0))
}
