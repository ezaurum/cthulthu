package runner

import (
	"github.com/ezaurum/boongeoppang/gin"
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/route"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
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
		if gin.IsDebugging() {
			render.NewDebug(config.Dir.Template, config.FuncMap, e)
		} else {
			//TODO e.HTMLRender = render.New(config.Dir.Template, config.FuncMap)
		}
	}

	// 스태틱 파일 설정
	if !helper.IsEmpty(config.Dir.Static) {
		e.Static("/", config.Dir.Static)
	}

	e.Start(config.Address)
}
