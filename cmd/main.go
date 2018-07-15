package main

import (
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/helper"
	"github.com/ezaurum/cthulthu/route"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	becho "github.com/ezaurum/boongeoppang/echo"
	"github.com/jinzhu/gorm"
	"github.com/ezaurum/cthulthu/database"
)

func main() {
	configFileName := "test-application.toml"
	cnf := &config.Config{}
	cnf.FromFile(configFileName)
	Run(cnf)
}

func Run(config *config.Config) {

	config.Generators = generators.New(func(_ string) generators.IDGenerator {
		return snowflake.New(config.NodeNumber)
	}, config.AutoMigrates...)

	db := initDB(config)

	defer db.Close()
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
		if e.Debug {
			e.Renderer = becho.NewDebug(config.Dir.Template, config.FuncMap)
		} else {
			e.Renderer = becho.New(config.Dir.Template, config.FuncMap)
		}
	}

	// 스태틱 파일 설정
	if !helper.IsEmpty(config.Dir.Static) {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:    config.Dir.Static,
			Skipper: middleware.DefaultSkipper,
		}))
	}

	e.Start(config.Address)
}

func initDB(config *config.Config) *gorm.DB {
	//Init DB
	db, err := gorm.Open(config.Db.Dialect, config.Db.Connection)
	if err != nil {
		panic(err)
	}
	database.RegisterAutoIDAssign(db, config.Generators)
	db.SingularTable(true)
	db.AutoMigrate(config.AutoMigrates...)
	return db
}
