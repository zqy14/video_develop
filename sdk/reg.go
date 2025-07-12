package sdk

import (
	"fmt"
	"regexp"
	"time"
)

// 验证手机号格式
func ValidatePhone(mobile string) bool {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return matched
}

// 邮箱校验
func ValidateEmail(email string) error {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    if !matched {
        return errors.New("邮箱格式错误")
    }
    return nil
}

// 格式化时间显示
func FormatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 5*time.Minute {
		return fmt.Sprintf("%d分钟之前", int(diff.Minutes()))
	} else if diff < 30*time.Minute {
		return "半小时之前"
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d小时之前", int(diff.Hours()))
	} else if diff < 3*24*time.Hour {
		return fmt.Sprintf("%d天之前", int(diff.Hours()/24))
	} else {
		return t.Format("2025-06-13 11:34")
	}
}
