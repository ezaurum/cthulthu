package main

import (
	"github.com/ezaurum/cthulthu/identity"
	"github.com/ezaurum/cthulthu/runner"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
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
		initialize(&cf)
	}
	return &cf
}

func initialize(config *config.Config) {

	if nil == config.DBManager {
		panic("DB manager cannot be nil")
	}

	manager := config.DBManager

	identity.CreateIdentityByForm(identity.FormIDToken{
		AccountName:"likepc",
		AccountPassword:"like#pc$0218",
	}, manager)
}

func main() {
	runner.DefaultRun(makeConfig())
}
