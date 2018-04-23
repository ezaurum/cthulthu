package runner

import (
	"github.com/ezaurum/boongeoppang/gin"
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/identity"
	"github.com/ezaurum/cthulthu/route"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/generators"
)

func Run(config *config.Config) {

	idGenerators := generators.New(func() generators.IDGenerator {
		return snowflake.New(config.NodeNumber)
	}, config.AutoMigrates...)

	//Init DB
	db, err := database.Open(idGenerators, config.Db.Dialect, config.Db.Connection)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(config.AutoMigrates...)

	config.DB = db

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
		identity.DefaultMiddleware(config.NodeNumber, db, r,
			config.SessionExpiresInSeconds,
			config.AuthorizerConfig...)
	}

	// 라우터는 핸들러가 추가되고 나서
	route.InitRoute(r, config.Routes...)

	// 템틀릿 렌더러 설정
	if !helper.IsEmpty(config.Dir.Template) {
		if gin.IsDebugging() {
			render.NewDebug(config.Dir.Template, r)
		} else {
			r.HTMLRender = render.New(config.Dir.Template)
		}
	}

	// 스태틱 파일 설정
	if !helper.IsEmpty(config.Dir.Static) {
		route.InitStaticFiles(r, config.Dir.Static)
	}

	r.Run(config.Address)
}
