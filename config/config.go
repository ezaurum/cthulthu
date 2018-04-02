package config

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/gin-gonic/gin"
)

type Config struct {
	DBManager               *database.Manager
	ConnectionString 		string
	Dialect			 		string
	AutoMigrates            []interface{}

	TemplateDir             string
	StaticDir               string

	NodeNumber              int64
	SessionExpiresInSeconds int
	AuthorizerConfig        []interface{}

	Routes					[]func() route.Routes

	Initialize				func(engine *gin.Engine)
}
