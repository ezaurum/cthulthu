package admin

import (
	"github.com/ezaurum/cthulthu"
	"github.com/ezaurum/cthulthu/authenticator"
	"github.com/ezaurum/cthulthu/authorizer"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/render"
	"github.com/ezaurum/cthulthu/route"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"github.com/ezaurum/cthulthu/session"
)

func DefaultRun() {
	Run(":9999")
}

func Run(addr ...string) {

	//DB
	manager := database.Default()
	db := manager.Connect()
	defer db.Close()

	r := initialize(&Config{
		DBManager:manager,
		TemplateDir:templateDir,
		StaticDir:staticDir,
		NodeNumber:0,
		SessionExpiresInSeconds:session.DefaultSessionExpires,
		AutoMigrates:[]interface{}{&Identity{}, &CookieIDToken{}, &FormIDToken{}},
		AuthorizerConfig:defaultConfig,
	})

	// run
	r.Run(addr...)
}

func initialize(config *Config) *gin.Engine {

	manager := config.DBManager

	//TODO related
	manager.AutoMigrate(config.AutoMigrates...)

	//gin
	r := gin.Default()
	// 기본값으로 만들고
	ca := authenticator.NewMem(config.NodeNumber, config.SessionExpiresInSeconds)
	ca.SetActions(GetLoadCookieIDToken(manager), GetLoadIdentity(manager), GetPersistToken(manager))
	r.Use(cthulthu.GinMiddleware(ca).Handler())
	// authenticator 를 초기화한다
	authorizer.Init(r, config.AuthorizerConfig...)
	//TODO login redirect page 지정 필요
	r.Use(manager.Handler())
	// renderer
	r.HTMLRender = render.New(config.TemplateDir)
	// static
	//TODO 디렉토리 여러 군데서 찾도록 하는 것도 필요
	//TODO skin 시스템을 상속가능하도록 하려면...
	staticDir := config.StaticDir
	route.AddAll(r, Login())
	route.AddAll(r, Register())
	r.Static("/js", staticDir+"/js")
	r.Static("/fonts", staticDir+"/fonts")
	r.Static("/css", staticDir+"/css")
	filepath.Walk(staticDir, route.SetStaticFile(r))
	return r
}
