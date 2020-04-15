package app

import (
	"github.com/ezaurum/cthulthu/context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
)

import _ "github.com/go-sql-driver/mysql"

func Start() {
	initialize()
	// 웹 초기화
	e := echo.New()
	if err := context.Ctx().InitRoute(e); nil != err {
		log.Fatal("route error", err)
	}

	addr := os.Getenv("SS_HOST_ADDR")
	if len(addr) < 1 {
		addr = ":9998"
	}
	log.Fatal(e.Start(addr))
}
