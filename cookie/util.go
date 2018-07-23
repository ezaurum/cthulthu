package cookie

import (
	"github.com/labstack/echo"
	"net/http"
)

func ClearCookie(c echo.Context, cookieName string) {
	c.SetCookie(&http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		Domain:   "",
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})
}
