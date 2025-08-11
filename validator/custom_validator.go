package validator

import (
	"regexp"
	"unicode"
)

// 检查字符串是否只包含中文字符
func IsChinese(str string) bool {
	for _, r := range str {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return len(str) > 0
}

// 检查手机号格式 (中国手机号)
func IsPhone(phone string) bool {
	// 匹配1开头，第二位3-9，后面9位数字
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

// 检查字符串是否为纯数字
func IsNumber(str string) bool {
	matched, _ := regexp.MatchString(`^\d+$`, str)
	return matched
}

// 生成幂等键
func GenerateIdempotentKey(phone, name string) string {
	return "rider_add_" + phone + "_" + name
}
