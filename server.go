package cthulthu

import (
	"github.com/ezaurum/cthulthu/render"
	"github.com/gin-gonic/gin"
)

func Run(addr ...string) {

	r := gin.Default()
	r.HTMLRender = render.Default()

	//TODO disconnect := .SetupDB(r)
	//TODO defer disconnect()

	//TODO congkong.SetupRoute(r)

	r.Run(addr...)
}
