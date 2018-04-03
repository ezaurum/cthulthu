package runner

import (
	"github.com/ezaurum/cthulthu/config"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/render"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/authorizer"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/identity"
)

func Run(config *config.Config) {

	//Init DB
	manager := database.New(config.ConnectionString, config.Dialect, config.NodeNumber)

	db := manager.Connect()
	defer db.Close()

	manager.AutoMigrate(config.AutoMigrates...)

	config.DBManager = manager

	//웹 전에 초기화 해야 하는 것들
	if nil != config.OnInitializeDB {
		config.OnInitializeDB()
	}

	// 웹 초기화

	r := gin.Default()

	// 엔진 넘겨주고 시작
	if nil != config.Initialize {
		config.Initialize(r)
	}

	// 미들웨어 사용 초기화
	if nil != config.InitializeMiddleware {
		config.InitializeMiddleware(r)
	} else {
		// authenticator 를 초기화한다
		ca := authenticator.NewMem(config.NodeNumber, config.SessionExpiresInSeconds)
		ca.SetActions(identity.GetLoadCookieIDToken(manager),
			identity.GetLoadIdentity(manager),
			identity.GetPersistToken(manager))
		if len(config.AuthorizerConfig) > 0 {
			au := authorizer.Init(config.AuthorizerConfig...)
			r.Use(ca.Handler(), manager.Handler(), au.Handler())
		} else {
			r.Use(ca.Handler(), manager.Handler())
		}
	}

	// 라우터는 핸들러가 추가되고 나서
	route.InitRoute(r, config.Routes...)

	// 템틀릿 렌더러 설정
	if !helper.IsEmpty(config.TemplateDir) {
		r.HTMLRender = render.New(config.TemplateDir)
	}

	r.Run(config.Address)
}
