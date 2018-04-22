package database

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	contextKey = "cthulthu-DBM"
)

func Handler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetDatabase(c, db)
	}
}

func SetDatabase(c *gin.Context, db *gorm.DB) {
	c.Set(contextKey, db)
}

func GetDatabase(c *gin.Context) *gorm.DB {
	return c.MustGet(contextKey).(*gorm.DB)
}
