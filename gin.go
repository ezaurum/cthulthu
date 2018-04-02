package cthulthu

import (
	"github.com/gin-gonic/gin"
)

// gin milddleware
type GinMiddleware interface {
	Handler() gin.HandlerFunc
}
