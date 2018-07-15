package echo

import (
	"strconv"
	"fmt"
	"time"
	"github.com/labstack/echo"
	"github.com/ezaurum/cthulthu/helper"
)

func BindID(c echo.Context, paramName string) (int64, bool) {
	accountIDString := c.Param(paramName)
	if "new" == accountIDString {
		return 0, false
	}
	accountID, err := strconv.ParseInt(accountIDString, 10, 64)
	if nil != err {
		fmt.Sprintf("[ERROR][%v] Parse error %v\n", time.Now(), err)
		return 0, false
	}
	return accountID, true
}

func BindQueryID(c echo.Context, queryIDName string) (int64, bool) {
	idString := c.QueryParam(queryIDName)
	if "new" == idString {
		return 0, false
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if nil != err {
		panic(fmt.Sprintf("[ERROR][%v] Parse error %v\n", time.Now(), err))
	}
	return id, true
}

func BindQueryIDString(c echo.Context, queryIDName string) (int64, bool, string) {
	idString := c.QueryParam(queryIDName)
	id, err := strconv.ParseInt(idString, 10, 64)
	if nil != err {
		return 0, false, idString
	}
	return id, true, idString
}

func BindIDString(c echo.Context, paramName string) (int64, bool, string) {
	accountIDString := c.Param(paramName)
	accountID, err := strconv.ParseInt(accountIDString, 10, 64)
	if nil != err {
		return 0, false, accountIDString
	}
	return accountID, true, accountIDString
}


func BindPhoneQuery(c echo.Context) string {
	phone := c.QueryParam("phone")
	if len(phone) > 0 {
		phone = helper.ExtractNumber(phone)
	}
	return phone
}
