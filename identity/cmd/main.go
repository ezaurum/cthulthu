package main

import (
	"github.com/ezaurum/cthulthu/identity"
	"github.com/ezaurum/cthulthu/runner"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu"
	"github.com/ezaurum/cthulthu/authorizer"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/render"
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/route"
)

func makeConfig() *config.Config {

	cf := config.Config{
		TemplateDir:            "./templates",
		StaticDir:                "./static",
		NodeNumber:              0,
		SessionExpiresInSeconds: session.DefaultSessionExpires,
		AuthorizerConfig: []interface{}{"./model.conf", "./policy.csv"},

		ConnectionString:"test.db",
		Dialect:"sqlite3",
		AutoMigrates:     []interface{}{&identity.Identity{}, &identity.CookieIDToken{}, &identity.FormIDToken{}, &identity.OAuthIDToken{}, } ,

		Routes: []func() route.Routes{ identity.Login, identity.Register},
	}

	cf.Initialize = func(engine *gin.Engine) {
		initialize(&cf, engine)
	}
	return &cf
}

func initialize(config *config.Config,
	r *gin.Engine) {

	if nil == config.DBManager {
		panic("DB manager cannot be nil")
	}

	manager := config.DBManager

	//TODO related
	manager.AutoMigrate(config.AutoMigrates...)

	identity.CreateIdentityByForm(identity.FormIDToken{
		AccountName:"likepc",
		AccountPassword:"like#pc$0218",
	}, manager)

	// 기본값으로 만들고
	// authenticator 를 초기화한다
	ca := authenticator.NewMem(config.NodeNumber, config.SessionExpiresInSeconds)
	ca.SetActions(identity.GetLoadCookieIDToken(manager),
		identity.GetLoadIdentity(manager),
		identity.GetPersistToken(manager))
	r.Use(cthulthu.GinMiddleware(ca).Handler())
	if len(config.AuthorizerConfig) > 0 {
		authorizer.Init(r, config.AuthorizerConfig...)
	}

	r.Use(manager.Handler())
	//TODO login redirect page 지정 필요
	// renderer
	if !helper.IsEmpty(config.TemplateDir) {
		r.HTMLRender = render.New(config.TemplateDir)
	}
}

func main() {
	runner.DefaultRun(makeConfig())
}
