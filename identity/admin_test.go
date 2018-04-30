package identity

import (
	"github.com/ezaurum/boongeoppang/gin"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/helper"
	itest "github.com/ezaurum/cthulthu/identity/test"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/url"
	"reflect"
	"testing"
	"github.com/ezaurum/cthulthu/generators/snowflake"
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
		Static:   "test/static",
		Template: "test/templates",
	},
	NodeNumber: 0,
	//SessionExpiresInSeconds: expires,
	AutoMigrates:     []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}},
	AuthorizerConfig: []interface{}{"test/model.conf", "test/policy.csv"},
}

func initializeTest(gens generators.IDGenerators, manager *gorm.DB, expires int) (*gin.Engine, *config.Config) {

	name := reflect.TypeOf(&Identity{}).String()
	idGenerator := gens[name]

	tc := config.Config{
		SessionExpiresInSeconds: expires,
		DB: manager,
		Dir: config.DirConfig{
			Static:   "test/static",
			Template: "test/templates",
		},
		AuthorizerConfig: testConfig.AuthorizerConfig,
		AutoMigrates:     testConfig.AutoMigrates,
		NodeNumber:       testConfig.NodeNumber,
		Routes:           []func() route.Routes{Login, MakeRegister(idGenerator)},
	}
	engine := gin.Default()

	manager.AutoMigrate(tc.AutoMigrates...)

	CreateIdentityByForm(FormIDToken{
		AccountName:     "test",
		AccountPassword: "test",
	}, manager, idGenerator)

	// 기본값으로 만들고

	DefaultMiddleware(tc.NodeNumber, manager, engine, tc.SessionExpiresInSeconds,
		tc.AuthorizerConfig...)

	// renderer
	if !helper.IsEmpty(tc.Dir.Template) {
		engine.HTMLRender = render.New(tc.Dir.Template, nil)
	}

	return engine, &tc
}

func TestRoleAccess(t *testing.T) {

	targets := []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}}
	gens := snowflake.GetGenerators(0, targets...)

	testDB := itest.DB(gens)
	defer testDB.Close()
	r, _ := initializeTest(gens, testDB, authenticator.DefaultSessionExpires)

	client := test.HttpClient{}
	w0 := client.GetRequest(r, "/")
	assert.True(t, test.IsRedirect(w0))
	assert.Equal(t, "/login", test.GetRedirect(w0))
}
