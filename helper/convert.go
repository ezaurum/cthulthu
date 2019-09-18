package helper

import (
	"strconv"
	"strings"
)

func ToInt64(param string) (int64, bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(param), 10, 64)
	if nil != err {
		return 0, false
	}
	return id, true
}

func ToInt64FromHex(param string) (int64, bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(param), 16, 64)
	if nil != err {
		return 0, false
	}
	return id, true
}

func ToInt(param string) (int, bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(param), 10, 32)
	if nil != err {
		return 0, false
	}
	return int(id), true
}

func ToFloat64(param string) (float64, bool) {
	id, err := strconv.ParseFloat(strings.TrimSpace(param), 64)
	if nil != err {
		return 0, false
	}
	return id, true
}
