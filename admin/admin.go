package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/authorizer"
	"github.com/ezaurum/cthulthu/render"
	"path/filepath"
	"os"
	"log"
)

const (
	staticDir = "admin/static"
	templateDir = "admin/templates"
)

var (
	defaultConfig = []interface{}{"admin/model.conf", "admin/policy.csv",}
)

func DefaultRun() {
	Run(":9999")
}

func Run(addr...string) {

	r := gin.Default()

	// 기본값으로 만들고
	// authenticator 를 초기화한다
	authorizer.InitWithAuthenticator(r, defaultConfig...)

	r.HTMLRender = render.New(templateDir)

	//TODO 디렉토리 여러 군데서 찾도록 하는 것도 필요
	//TODO skin 시스템을 상속가능하도록 하려면...

	// add static dirs
	r.Static("/js", staticDir + "/js")
	r.Static("/fonts", staticDir + "/fonts")
	r.Static("/css",staticDir +  "/css")

	filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if nil != err {
			log.Printf("err before %v, %v", path, err)
		}

		// 디렉토리는 패스
		if info.IsDir() {
			return nil
		}

		//TODO ignored files?

		base := filepath.Base(path)
		r.StaticFile(base,path)
		return nil
	})


	//TODO disconnect := .SetupDB(r)
	//TODO defer disconnect()

	//TODO congkong.SetupRoute(r)

	r.Run(addr...)
}
