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
)

const (
	staticDir   = "admin/static"
	templateDir = "admin/templates"
)

var (
	defaultConfig = []interface{}{"admin/model.conf", "admin/policy.csv"}
)

func DefaultRun() {
	Run(":9999")
}

func Run(addr ...string) {

	//DB
	manager := database.Default()
	db := manager.Connect()
	defer db.Close()

	r := initialize(manager, templateDir, staticDir, defaultConfig...)

	// run
	r.Run(addr...)
}

func initialize(manager *database.Manager, templateDir string, staticDir string, params ...interface{}) *gin.Engine {

	//TODO related
	manager.AutoMigrate(&FormIDToken{}, &CookieIDToken{}, &Identity{})

	//gin
	r := gin.Default()
	// 기본값으로 만들고
	ca := authenticator.Default()
	ca.SetActions(GetLoadCookieIDToken(manager), GetLoadIdentity(manager))
	r.Use(ca.(cthulthu.GinMiddleware).Handler())
	// authenticator 를 초기화한다
	authorizer.Init(r, params...)
	//TODO login redirect page 지정 필요
	r.Use(manager.Handler())
	// renderer
	r.HTMLRender = render.New(templateDir)
	// static
	//TODO 디렉토리 여러 군데서 찾도록 하는 것도 필요
	//TODO skin 시스템을 상속가능하도록 하려면...
	route.AddAll(r, Login())
	route.AddAll(r, Register())
	r.Static("/js", staticDir+"/js")
	r.Static("/fonts", staticDir+"/fonts")
	r.Static("/css", staticDir+"/css")
	filepath.Walk(staticDir, route.SetStaticFile(r))
	return r
}
