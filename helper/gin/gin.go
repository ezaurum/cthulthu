package gin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"time"
	"strings"
	"github.com/ezaurum/cthulthu/helper"
)

func BindID(c *gin.Context, paramName string) (int64, bool) {
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

func BindQueryID(c *gin.Context, queryIDName string) (int64, bool) {
	idString := c.Query(queryIDName)
	if "new" == idString {
		return 0, false
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if nil != err {
		panic(fmt.Sprintf("[ERROR][%v] Parse error %v\n", time.Now(), err))
	}
	return id, true
}

func BindQueryIDString(c *gin.Context, queryIDName string) (int64, bool, string) {
	idString := c.Query(queryIDName)
	id, err := strconv.ParseInt(idString, 10, 64)
	if nil != err {
		return 0, false, idString
	}
	return id, true, idString
}

func BindIDString(c *gin.Context, paramName string) (int64, bool, string) {
	accountIDString := c.Param(paramName)
	accountID, err := strconv.ParseInt(accountIDString, 10, 64)
	if nil != err {
		return 0, false, accountIDString
	}
	return accountID, true, accountIDString
}

func BindPhoneQuery(c *gin.Context) string {
	phone := c.Query("phone")
	if len(phone) > 0 {
		phone = helper.ExtractNumber(phone)
	}
	return phone
}