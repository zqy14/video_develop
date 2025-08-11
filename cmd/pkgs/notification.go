package pkgs

import "log"

// 通知接口
type Notifier interface {
	Send(message string) error
}

// 短信通知实现
type SMSNotifier struct{}

// 发送短信通知
func (s *SMSNotifier) Send(message string) error {
	// 这里实现实际的短信发送逻辑
	log.Printf("发送短信通知: %s", message)
	return nil
}
