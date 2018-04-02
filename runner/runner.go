package runner

import (
	"github.com/ezaurum/cthulthu/config"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/render"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/identity"
	"github.com/ezaurum/cthulthu/authorizer"
)

const (
	defaultAddress = ":9999"
)

func DefaultRun(config *config.Config) {
	//TODO 어드레스를 실행시나 콘피그에서 가져올 수 있도록
	Run(config, defaultAddress)
}

func Run(config *config.Config, addr...string) {

	//Init DB
	manager := database.New(config.ConnectionString, config.Dialect, config.NodeNumber)

	db := manager.Connect()
	defer db.Close()

	// DB 초기화
	manager.AutoMigrate(config.AutoMigrates...)

	config.DBManager = manager

	r := gin.Default()

	config.Initialize(r)

	// authenticator 를 초기화한다
	ca := authenticator.NewMem(config.NodeNumber, config.SessionExpiresInSeconds)
	ca.SetActions(identity.GetLoadCookieIDToken(manager),
		identity.GetLoadIdentity(manager),
		identity.GetPersistToken(manager))
	if len(config.AuthorizerConfig) > 0 {
		authorizer.Init(r, config.AuthorizerConfig...)
	}

	route.InitRoute(r, config.Routes...)

	// DB 핸들러 설정
	r.Use(ca.Handler(), manager.Handler())

	// 템틀릿 렌더러 설정
	if !helper.IsEmpty(config.TemplateDir) {
		r.HTMLRender = render.New(config.TemplateDir)
	}

	r.Run(addr...)
}
