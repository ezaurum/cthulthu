package config

import (
	"github.com/ezaurum/cthulthu/identity"
)

const (
	defaultAddress = ":9999"
)

func DefaultRun(config *Config) {
	//TODO 어드레스를 실행시나 콘피그에서 가져올 수 있도록
	Run(config, defaultAddress)
}

func Run(config *Config, addr ...string) {

	//DB
	manager := config.DBManager
	db := manager.Connect()
	defer db.Close()

	r := config.Initializer()

	config.InitRoute(r, config.DefaultRoutes...)
	identity.InitStateFiles(r, config)

	// run
	r.Run(addr...)
}
