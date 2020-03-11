package echo

import (
	"github.com/ezaurum/cthulthu/cmd/identity"
	cookie2 "github.com/ezaurum/cthulthu/cookie"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

const (
	IdentityTokenCookieName = "cthulthu-identity-token"
	IdentityContextKey      = "cthulthu-identity"
)

func MakeMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, e := c.Cookie(IdentityTokenCookieName)
			if _, b := GetIdentity(c); !b && nil == e {
				if er := SetIdentityToContext(db, c, cookie.Value);
					nil != er {
					return er
				}
			}
			return next(c)
		}
	}
}

func SetIdentityToContext(db *gorm.DB, c echo.Context, cookieValue string) error {
	var idToken identity.IdentifyToken
	r := db.Where("token = ?", cookieValue).Find(&idToken)
	if nil != r.Error {
		if gorm.IsRecordNotFoundError(r.Error) {
			cookie2.ClearCookie(c, IdentityTokenCookieName)
			return nil
		} else {
			return r.Error
		}
	}
	var ident identity.Identity
	r = db.Where("id = ?", idToken.IdentityID).Find(&ident)
	if nil != r.Error {
		if gorm.IsRecordNotFoundError(r.Error) {
			db.Delete(&idToken, idToken.ID)
			cookie2.ClearCookie(c, IdentityTokenCookieName)
		} else {
			return r.Error
		}
	}
	SetIdentity(c, &ident)
	return nil
}

func GetIdentity(c echo.Context) (*identity.Identity, bool) {
	get := c.Get(IdentityContextKey)
	if nil == get {
		return nil, false
	}
	return get.(*identity.Identity), true
}

func SetIdentity(c echo.Context, ident *identity.Identity) {
	c.Set(IdentityContextKey, &ident)
}
