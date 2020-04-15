package query

import (
	"github.com/ezaurum/cthulthu/conv"
	"github.com/labstack/echo/v4"
)

func (p *Param) FromContextInt64(c echo.Context, key string) bool {
	toInt64, hasKey := conv.ToInt64(c.QueryParam(key))
	if hasKey {
		p.SetValue(key, toInt64)
		return true
	}
	return false
}

func (p *Param) FromContextString(c echo.Context, key string) bool {
	keyValue := c.QueryParam(key)
	if len(keyValue) > 0 {
		p.SetValue(key, keyValue)
		return true
	}
	return false
}

func (p *Param) SetValue(key string, value interface{}) {
	if p.QueryValues != nil {
		p.QueryValues[key] = value
	} else {
		p.QueryValues = map[string]interface{}{
			key: value,
		}
	}
}

func New(c echo.Context) (*Param, error) {
	var p Param
	if err := c.Bind(&p); nil != err {
		return nil, err
	}
	return &p, nil
}
