package cthulthu

import (
	"github.com/gin-gonic/gin"
	"bitbucket.org/congkong-revivals/admin/congkong"
)

func Run(config string)  {

	r := gin.Default()
	r.HTMLRender = gin.Default()

	disconnect := congkong.SetupDB(r)
	defer disconnect()

	congkong.SetupRoute(r)

	r.Run(":8989")
}
