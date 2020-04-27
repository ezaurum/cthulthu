package app

import (
	"github.com/ezaurum/cthulthu/context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
)

import _ "github.com/go-sql-driver/mysql"

func Start(overrideAddr string) {
	Initialize()
	// 웹 초기화
	e := echo.New()
	if err := context.App().InitRoute(e); nil != err {
		log.Fatal("route error", err)
	}

	addr := os.Getenv("SS_HOST_ADDR")
	if len(addr) < 1 {
		addr = ":9998"
	}
	if len(overrideAddr) > 0 {
		addr = overrideAddr
	}
	log.Fatal(e.Start(addr))
}
