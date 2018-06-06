package config

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/route"
	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"html/template"
	"github.com/labstack/echo"
)

type Config struct {
	DB           *gorm.DB
	AutoMigrates []interface{}

	NodeNumber              int64
	Generators              generators.IDGenerators
	SessionExpiresInSeconds int
	AuthorizerConfig        []interface{}

	Routes []func() route.Routes

	OnInitializeDB       func()
	Initialize           func(engine *echo.Echo)
	InitializeMiddleware func(engine *echo.Echo)

	FuncMap template.FuncMap

	Address string

	Db  DBConfig
	Dir DirConfig
}

type DBConfig struct {
	Connection string
	Dialect    string
}

type DirConfig struct {
	Static   string
	Template string
}

func (cnf *Config) FromFile(configFile string) {
	toml, err := toml.LoadFile(configFile)
	if nil != err {
		panic(err)
	}
	toml.Unmarshal(cnf)
}
