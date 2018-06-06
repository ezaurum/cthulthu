package echo

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

const (
	contextKey = "cthulthu-DBM"
)

func Handler(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			SetDatabase(c, db)
			return next(c)
		}
	}
}

func SetDatabase(c echo.Context, db *gorm.DB) {
	c.Set(contextKey, db)
}

func GetDatabase(c echo.Context) *gorm.DB {
	return c.Get(contextKey).(*gorm.DB)
}
