package database

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func (dbm *Manager) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("WTF ma")
		fmt.Println("WTF ma")
		fmt.Println("WTF ma")
		SetDatabase(c, dbm)
	}
}

func SetDatabase(c *gin.Context, db *Manager) {
	c.Set("DBM", db)
}

func GetDatabase(c *gin.Context) *Manager {
	return c.MustGet("DBM").(*Manager)
}
