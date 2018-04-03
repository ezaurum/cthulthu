package identity

import (
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/authorizer"
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/render"
	"github.com/ezaurum/cthulthu/session"
	"github.com/ezaurum/cthulthu/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"github.com/ezaurum/cthulthu/route"
)

func Initialize(config *config.Config,
	r *gin.Engine) {

	if nil == config.DBManager {
		panic("DB manager cannot be nil")
	}

	manager := config.DBManager

	//TODO related
	manager.AutoMigrate(config.AutoMigrates...)

	CreateIdentityByForm(FormIDToken{
		AccountName:     "likepc",
		AccountPassword: "like#pc$0218",
	}, manager)

	// 기본값으로 만들고
	// authenticator 를 초기화한다
	ca := authenticator.NewMem(config.NodeNumber, config.SessionExpiresInSeconds)
	ca.SetActions(GetLoadCookieIDToken(manager),
		GetLoadIdentityByCookie(manager),
		GetPersistToken(manager))
	var au authorizer.AuthorizeMiddleware
	if len(config.AuthorizerConfig) > 0 {
		au = authorizer.GetAuthorizer(config.AuthorizerConfig...)
	}

	r.Use(ca.Handler(), manager.Handler(), au.Handler())
	//TODO login redirect page 지정 필요
	// renderer
	if !helper.IsEmpty(config.TemplateDir) {
		r.HTMLRender = render.New(config.TemplateDir)
	}
}

func getRegisterFormPostData() url.Values {
	form := make(url.Values)
	form.Set("accountName", "test")
	form.Set("accountPassword", "test")
	return form
}

var testConfig = &config.Config{
	//DBManager:               manager,
	TemplateDir: "test/static",
	StaticDir:   "test/templates",
	NodeNumber:  0,
	//SessionExpiresInSeconds: expires,
	AutoMigrates:     []interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}},
	AuthorizerConfig: []interface{}{"test/model.conf", "test/policy.csv"},
}

func initializeTest(manager *database.Manager, expires int) (*gin.Engine, *config.Config) {
	tc := config.Config{
		SessionExpiresInSeconds: expires,
		DBManager:               manager,
		TemplateDir:             testConfig.TemplateDir,
		AuthorizerConfig:        testConfig.AuthorizerConfig,
		AutoMigrates:            testConfig.AutoMigrates,
		StaticDir:               testConfig.StaticDir,
		NodeNumber:              testConfig.NodeNumber,
		Routes: []func() route.Routes{ Login, Register},
	}
	engine := gin.Default()
	Initialize(&tc, engine)
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
