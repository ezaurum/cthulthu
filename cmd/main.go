package main

import (
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/strcheck"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	becho "github.com/ezaurum/boongeoppang/echo"
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

	db := config.InitDB()
	defer db.Close()

	// 웹 초기화

	e := echo.New()

	// 템틀릿 렌더러 설정
	if !strcheck.IsEmpty(config.Dir.Template) {
		if e.Debug {
			e.Renderer = becho.NewDebug(config.Dir.Template, config.FuncMap)
		} else {
			e.Renderer = becho.New(config.Dir.Template, config.FuncMap)
		}
	}

	// 스태틱 파일 설정
	if !strcheck.IsEmpty(config.Dir.Static) {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:    config.Dir.Static,
			Skipper: middleware.DefaultSkipper,
		}))
	}

	e.Start(config.Address)
}
