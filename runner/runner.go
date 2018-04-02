package runner

import (
	"github.com/ezaurum/cthulthu/config"
	"github.com/ezaurum/cthulthu/route"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/database"
)

const (
	defaultAddress = ":9999"
)

func DefaultRun(config *config.Config) {
	//TODO 어드레스를 실행시나 콘피그에서 가져올 수 있도록
	Run(config, defaultAddress)
}

func Run(config *config.Config, addr...string) {

	//Init DB
	manager := database.New(config.ConnectionString, config.Dialect, config.NodeNumber)

	db := manager.Connect()
	defer db.Close()

	config.DBManager = manager

	r := gin.Default()

	config.Initialize(r)

	route.InitRoute(r, config.Routes...)

	r.Run(addr...)
}
