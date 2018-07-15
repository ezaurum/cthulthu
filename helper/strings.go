package helper

import (
	"strings"
	"regexp"
)

var (
	OnlyNumberReg = regexp.MustCompile("[0-9]+")
)

func Contains(target []string, substr string) bool {
	for _, v := range target {
		if strings.Contains(v, substr) {
			return true
		}
	}
	return false
}

func IsEmpty(target string) bool {
	return len(target) < 1
}

func ExtractNumber(phone string) string {
	return strings.Join(OnlyNumberReg.FindAllString(phone, -1), "")
}
