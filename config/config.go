package config

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/identity"
	"github.com/gin-gonic/gin"
)

const (
	staticDir   = "identity/static"
	templateDir = "identity/templates"
)

var (
	DefaultConfig       = []interface{}{"identity/model.conf", "identity/policy.csv"}
	DefaultAutoMigrates = []interface{}{&identity.Identity{}, &identity.CookieIDToken{},
	&identity.FormIDToken{}, &identity.OAuthIDToken{}, &identity.Score{}}
)

type Config struct {
	DBManager               *database.Manager
	TemplateDir             string
	StaticDir               string
	NodeNumber              int64
	SessionExpiresInSeconds int
	AuthorizerConfig        []interface{}
	AutoMigrates            []interface{}
	Initializer             func() *gin.Engine
	InitRoute               func(r *gin.Engine, routes ...func() route.Routes)
	DefaultRoutes           []func() route.Routes
}
