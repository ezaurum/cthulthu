package database

import (
	"github.com/gin-gonic/gin"
)

func (dbm *Manager) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		SetDatabase(c, dbm)
	}
}

func SetDatabase(c *gin.Context, db *Manager) {
	c.Set("DBM", db)
}

func GetDatabase(c *gin.Context) *Manager {
	return c.MustGet("DBM").(*Manager)
}
