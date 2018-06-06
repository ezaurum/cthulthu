package runner

import (
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/route"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	becho "github.com/ezaurum/boongeoppang/echo"
)

func Run(config *config.Config) {

	config.Generators = generators.New(func() generators.IDGenerator {
		return snowflake.New(config.NodeNumber)
	}, config.AutoMigrates...)

	//Init DB
	db, err := database.Open(config.Generators,
		config.Db.Dialect, config.Db.Connection)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.SingularTable(true)
	db.AutoMigrate(config.AutoMigrates...)

	config.DB = db

	//웹 전에 초기화 해야 하는 것들
	if nil != config.OnInitializeDB {
		config.OnInitializeDB()
	}

	// 웹 초기화

	e := echo.New()

	// 엔진 넘겨주고 시작
	if nil != config.Initialize {
		config.Initialize(e)
	}

	// 미들웨어 사용 초기화
	if nil != config.InitializeMiddleware {
		config.InitializeMiddleware(e)
	}

	// 라우터는 핸들러가 추가되고 나서
	route.InitRoute(e, config.Routes...)

	// 템틀릿 렌더러 설정
	if !helper.IsEmpty(config.Dir.Template) {
		becho.NewDebug(config.Dir.Template, config.FuncMap)

		//TODO 디버그 아닌 상태가 필요
		/*if gin.IsDebugging() {
			render.NewDebug(config.Dir.Template, config.FuncMap, e)
		} else {
			//TODO e.HTMLRender = render.New(config.Dir.Template, config.FuncMap)
		}*/
	}

	// 스태틱 파일 설정
	if !helper.IsEmpty(config.Dir.Static) {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:   config.Dir.Static,
			Skipper:middleware.DefaultSkipper,
		}))
	}

	e.Start(config.Address)
}
