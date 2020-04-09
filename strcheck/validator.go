package strcheck

import (
	"regexp"
	"strings"
)

func ValidatePhoneNumber(phone string, countryCode string) (bool, string) {

	// 숫자만
	phone = ExtractNumber(phone)
	if len(phone) < 1 {
		return false, phone
	}

	switch countryCode {
	case "82":
		// 앞에 82 있으면 버리기
		if strings.HasPrefix(phone, "82") {
			phone = phone[2:]
		}
		if strings.HasPrefix(phone, "+82") {
			phone = phone[3:]
		}
		// 한국인 경우 010으로 시작하도록
		if !strings.HasPrefix(phone, "0") {
			phone = "0" + phone
		}
		if !(strings.HasPrefix(phone, "010") ||
			strings.HasPrefix(phone, "011") ||
			strings.HasPrefix(phone, "018") ||
			strings.HasPrefix(phone, "017") ||
			strings.HasPrefix(phone, "019") ||
			strings.HasPrefix(phone, "016")) {
			return false, phone
		}

		// 전화번호는 10 or 11자리
		if len(phone) < 10 || len(phone) > 11 {
			return false, phone
		}
	}
	return true, phone
}

func IsEmpty(target string) bool {
	return len(target) < 1
}

func ExtractNumber(phone string) string {
	return strings.Join(OnlyNumberReg.FindAllString(phone, -1), "")
}

var (
	OnlyNumberReg = regexp.MustCompile("[0-9]+")
)
