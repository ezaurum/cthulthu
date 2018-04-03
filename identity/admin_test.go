package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/render"
)

func getRegisterFormPostData() url.Values {
	form := make(url.Values)
	form.Set("accountName", "test")
	form.Set("accountPassword", "test")
	return form
}

var testConfig = &config.Config{
	//DBManager:               manager,
	Dir: config.DirConfig{
		Static:"test/static",
		Template:   "test/templates",
	},
	NodeNumber:  0,
	//SessionExpiresInSeconds: expires,
	AutoMigrates:     []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}},
	AuthorizerConfig: []interface{}{"test/model.conf", "test/policy.csv"},
}

func initializeTest(manager *database.Manager, expires int) (*gin.Engine, *config.Config) {
	tc := config.Config{
		SessionExpiresInSeconds: expires,
		DBManager:               manager,
		Dir: config.DirConfig{
			Static:"test/static",
			Template:   "test/templates",
		},
		AuthorizerConfig:        testConfig.AuthorizerConfig,
		AutoMigrates:            testConfig.AutoMigrates,
		NodeNumber:              testConfig.NodeNumber,
		Routes: []func() route.Routes{ Login, Register},
	}
	engine := gin.Default()

	manager.AutoMigrate(tc.AutoMigrates...)

	CreateIdentityByForm(FormIDToken{
		AccountName:     "test",
		AccountPassword: "test",
	}, manager)

	// 기본값으로 만들고

	DefaultMiddleware(tc.NodeNumber, manager, engine, tc.SessionExpiresInSeconds,
		tc.AuthorizerConfig...)

	// renderer
	if !helper.IsEmpty(tc.Dir.Template) {
		engine.HTMLRender = render.New(tc.Dir.Template)
	}

	return engine, &tc
}

func TestRoleAccess(t *testing.T) {

	testDB := database.Test()
	db := testDB.Connect()
	defer db.Close()
	r, _ := initializeTest(testDB, session.DefaultSessionExpires)

	client := test.HttpClient{}
	w0 := client.GetRequest(r, "/")
	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/login", test.GetRedirect(w0))
}
