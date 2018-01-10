package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func getRegisterFormPostData() url.Values {
	form := make(url.Values)
	form.Set("accountName", "test")
	form.Set("accountPassword", "test")
	return form
}

var testConfig = &Config{
	//DBManager:               manager,
	TemplateDir: testTemplateDir,
	StaticDir:   testStaticDir,
	NodeNumber:  0,
	//SessionExpiresInSeconds: expires,
	AutoMigrates:     []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}},
	AuthorizerConfig: testAuthorizerConfig,
}

func initializeTest(manager *database.Manager, expires int) *gin.Engine {
	tc := Config{
		SessionExpiresInSeconds: expires,
		DBManager:               manager,
		TemplateDir:             testConfig.TemplateDir,
		AuthorizerConfig:        testConfig.AuthorizerConfig,
		AutoMigrates:            testConfig.AutoMigrates,
		StaticDir:               testConfig.StaticDir,
		NodeNumber:              testConfig.NodeNumber,
	}
	return Initialize(&tc)
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
