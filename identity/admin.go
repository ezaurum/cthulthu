package identity

import (
	"github.com/ezaurum/cthulthu"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/authorizer"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/render"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func DefaultRun() {
	Run(":9999")
}

func Run(addr ...string) {

	//DB
	manager := database.Default()
	db := manager.Connect()
	defer db.Close()

	config := Config{
		DBManager:               manager,
		TemplateDir:             templateDir,
		StaticDir:               staticDir,
		NodeNumber:              0,
		SessionExpiresInSeconds: session.DefaultSessionExpires,
		AutoMigrates:            DefaultAutoMigrates,
		AuthorizerConfig:        DefaultConfig,
	}

	r := Initialize(&config)

	initRoute(r, &config)
	InitStateFiles(r, &config)

	// run
	r.Run(addr...)
}

func Initialize(config *Config) *gin.Engine {

	if nil == config.DBManager {
		panic("DB manager cannot be nil")
	}

	manager := config.DBManager

	//TODO related
	manager.AutoMigrate(config.AutoMigrates...)

	CreateIdentityByForm(FormIDToken{
		AccountName:"likepc",
		AccountPassword:"like#pc$0218",
	}, manager)
	//gin
	r := gin.Default()
	// 기본값으로 만들고
	// authenticator 를 초기화한다
	ca := authenticator.NewMem(config.NodeNumber, config.SessionExpiresInSeconds)
	ca.SetActions(GetLoadCookieIDToken(manager), GetLoadIdentity(manager), GetPersistToken(manager))
	r.Use(cthulthu.GinMiddleware(ca).Handler())
	if len(config.AuthorizerConfig) > 0 {
		authorizer.Init(r, config.AuthorizerConfig...)
	}

	r.Use(manager.Handler())
	//TODO login redirect page 지정 필요
	// renderer
	if helper.IsEmpty(config.TemplateDir) {
		r.HTMLRender = render.New(config.TemplateDir)
	}
	return r
}

func initRoute(r *gin.Engine, config *Config) {
	route.AddAll(r, Login())
	//route.AddAll(r, Register())

	//TODO 분리 필요
	route.AddAll(r, ScoreRoute())
}

func InitStateFiles(r *gin.Engine, config *Config) {
	// static
	//TODO 디렉토리 여러 군데서 찾도록 하는 것도 필요
	//TODO skin 시스템을 상속가능하도록 하려면...
	staticDir := config.StaticDir
	r.Static("/images", staticDir+"/images")
	r.Static("/js", staticDir+"/js")
	r.Static("/fonts", staticDir+"/fonts")
	r.Static("/css", staticDir+"/css")
	filepath.Walk(staticDir, route.SetStaticFile(r))
}
