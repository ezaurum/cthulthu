package admin

import (
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

	r := gin.Default()

	// 기본값으로 만들고
	// authenticator 를 초기화한다
	authorizer.InitWithAuthenticator(r, defaultConfig...)


	//TODO login redirect page 지정 필요

	//DB
	manager := database.Default()
	db := manager.Connect()
	defer db.Close()

	manager.AutoMigrate(&FormIDToken{}, &Identity{})

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

	// run
	r.Run(addr...)
}
