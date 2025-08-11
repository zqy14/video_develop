package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// SMSProvider 接口定义
type SMSProvider interface {
	Send(phone, message string) error
	Name() string
}

// 第一个供应商实现
type Vendor1 struct{}

func (v *Vendor1) Send(phone, message string) error {
	// 模拟85%成功率
	if rand.Intn(100) < 85 {
		fmt.Printf("[%s] 发送成功到 %s\n", v.Name(), phone)
		return nil
	}
	return errors.New("供应商1发送失败")
}

func (v *Vendor1) Name() string {
	return "供应商1"
}

// 第二个供应商实现
type Vendor2 struct{}

func (v *Vendor2) Send(phone, message string) error {
	// 模拟75%成功率
	if rand.Intn(100) < 75 {
		fmt.Printf("[%s] 发送成功到 %s\n", v.Name(), phone)
		return nil
	}
	return errors.New("供应商2发送失败")
}

func (v *Vendor2) Name() string {
	return "供应商2"
}

// SMSService 短信服务
type SMSService struct {
	vendors []SMSProvider
}

func NewSMSService() *SMSService {
	return &SMSService{
		vendors: []SMSProvider{
			&Vendor1{},
			&Vendor2{},
		},
	}
}

// Send 发送短信，尝试所有供应商直到成功
func (s *SMSService) Send(phone, message string) error {
	for i, vendor := range s.vendors {
		err := vendor.Send(phone, message)
		if err == nil {
			return nil // 发送成功
		}

		fmt.Printf("尝试 %s 失败: %v\n", vendor.Name(), err)

		// 如果不是最后一个供应商，稍作延迟再试下一个
		if i < len(s.vendors)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return errors.New("所有供应商尝试失败")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	service := NewSMSService()

	phone := "13800138000"
	message := "您的验证码是: 123456"

	// 尝试发送3次演示
	for i := 1; i <= 3; i++ {
		fmt.Printf("\n=== 第%d次发送尝试 ===\n", i)
		err := service.Send(phone, message)
		if err != nil {
			fmt.Println("发送失败:", err)
		} else {
			fmt.Println("发送成功!")
			break
		}
	}
}
