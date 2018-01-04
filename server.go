package cthulthu

import (
	"github.com/ezaurum/cthulthu/render"
	"github.com/gin-gonic/gin"
	"github.com/ezaurum/cthulthu/authorizer"
)

func Run(addr ...string) {

	r := gin.Default()

	// 기본값으로 만들고
	// authenticator 를 초기화한다
	authorizer.InitWithAuthenticator(r)

	r.HTMLRender = render.Default()

	//TODO disconnect := .SetupDB(r)
	//TODO defer disconnect()

	//TODO congkong.SetupRoute(r)

	r.Run(addr...)
}
