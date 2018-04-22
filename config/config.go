package config

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
)

type Config struct {
	DBManager    *database.Manager
	AutoMigrates []interface{}

	NodeNumber              int64
	SessionExpiresInSeconds int
	AuthorizerConfig        []interface{}

	Routes []func() route.Routes

	OnInitializeDB       func()
	Initialize           func(engine *gin.Engine)
	InitializeMiddleware func(engine *gin.Engine)

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
