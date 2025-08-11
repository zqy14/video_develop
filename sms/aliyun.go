package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type AliyunProvider struct{}

func NewAliyunProvider() *AliyunProvider {
	return &AliyunProvider{}
}

func (a *AliyunProvider) Send(_ context.Context, phone, content string) error {
	// 实现阿里云短信发送逻辑
	fmt.Printf("通过阿里云发送短信到 %s: %s\n", phone, content)
	return nil // 模拟成功
}

func (a *AliyunProvider) Name() string {
	return "aliyun"
}

// 腾讯云短信实现
type TencentSMS struct{}

func (t *TencentSMS) Send(phone, content string) error {
	fmt.Printf("尝试通过腾讯云发送短信到 %s: %s\n", phone, content)
	// 模拟80%成功率
	if time.Now().UnixNano()%10 < 2 { // 20%失败率
		return errors.New("腾讯云短信发送失败")
	}
	return nil
}

func (t *TencentSMS) Name() string {
	return "tencent"
}
