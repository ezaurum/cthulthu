package helper

import "strings"

func Contains(target []string, substr string) bool {
	for _, v := range target {
		if strings.Contains(v, substr) {
			return true
		}
	}
	return false
}

func IsEmpty(target string) bool {
	return len(target) <1
}
